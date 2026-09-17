package demo

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	designsystem "github.com/jobrunner/fieldworksdiary-designsystem"
)

// klassenRe findet die Klassennamen, die base.css definiert.
var klassenRe = regexp.MustCompile(`^\.([a-z][a-z0-9-]*)`)

// classAttrRe findet alle class="..."-Attribute in HTML.
var classAttrRe = regexp.MustCompile(`class="([^"]*)"`)

// benutzteKlassen extrahiert die tatsächlich vergebenen Klassennamen aus HTML,
// indem es nur class="..."-Attribute berücksichtigt. So werden id-Attribute,
// Klassenerwähnungen im Fließtext oder per display:none versteckte Elemente
// nicht fälschlicherweise als „gezeigt" gewertet.
func benutzteKlassen(html string) map[string]bool {
	klassen := make(map[string]bool)
	for _, m := range classAttrRe.FindAllStringSubmatch(html, -1) {
		if len(m) > 1 {
			for _, k := range strings.Fields(m[1]) {
				klassen[k] = true
			}
		}
	}
	return klassen
}

func TestBenutzteKlassenExtrahiertNurAusClassAttributen(t *testing.T) {
	tests := []struct {
		html         string
		erwartet     map[string]bool
		beschreibung string
	}{
		{
			`<span id="btn">Text</span>`,
			map[string]bool{},
			"id-Attribut ist keine Klasse",
		},
		{
			`<div class="card btn-row">Text</div>`,
			map[string]bool{"card": true, "btn-row": true},
			"mehrere Klassen im Attribut",
		},
		{
			`<p>Die Klasse btn ist blau.</p>`,
			map[string]bool{},
			"Klassenname im Text ergibt keine Klasse",
		},
		{
			`<div class="card">A</div><div class="btn">B</div>`,
			map[string]bool{"card": true, "btn": true},
			"mehrere Attribute zusammen",
		},
	}

	for _, test := range tests {
		got := benutzteKlassen(test.html)
		if len(got) != len(test.erwartet) {
			t.Errorf("%s: Länge %d, erwartet %d", test.beschreibung, len(got), len(test.erwartet))
		}
		for k := range test.erwartet {
			if !got[k] {
				t.Errorf("%s: .%s fehlt", test.beschreibung, k)
			}
		}
		for k := range got {
			if !test.erwartet[k] {
				t.Errorf("%s: unerwartete .%s", test.beschreibung, k)
			}
		}
	}
}

// kommentarRe entfernt CSS-Kommentare vor der Auswertung, damit ein
// Erklärtext wie "a { color }" in einer Begründung nicht als Regel
// mitgezählt wird.
var kommentarRe = regexp.MustCompile(`(?s)/\*.*?\*/`)

// regelKopfRe findet die Selektorliste jeder Regel: alles bis zum
// öffnenden "{". Das trifft auch auf @media- und @keyframes-Blöcke zu
// (deren Kopf beginnt mit "@" und wird unten verworfen) sowie auf die
// Regeln, die darin verschachtelt sind — die kommen als eigener Treffer.
var regelKopfRe = regexp.MustCompile(`([^{}]+)\{`)

// elementNameRe erkennt einen bloßen HTML-Elementnamen wie "a" oder
// "input", nicht aber "*" oder ":focus-visible".
var elementNameRe = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

// tagRe findet öffnende HTML-Tags samt ihrer Attribute in einem HTML-Text,
// um danach prüfen zu können, ob das Attribut class="skip-link" dabei ist.
var tagRe = regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)([^>]*)>`)

// istSprunglink erkennt, ob die Attribute eines Tags die Klasse
// "skip-link" tragen. Der Sprunglink ist per Gestaltung unsichtbar
// (position: absolute; left: -999px, erst bei Fokus sichtbar über
// .skip-link:focus) und dient nur der Tastaturnavigation, nicht der
// Abnahme der Komponente — ein <a class="skip-link"> darf deshalb nicht
// als Beleg dafür zählen, dass die Link-Komponente auf der Seite zu sehen
// ist.
func istSprunglink(attribute string) bool {
	m := classAttrRe.FindStringSubmatch(attribute)
	if m == nil {
		return false
	}
	for _, klasse := range strings.Fields(m[1]) {
		if klasse == "skip-link" {
			return true
		}
	}
	return false
}

// keineKomponente führt Elemente, die base.css zwar per Element-Selektor
// gestalten könnte, die aber keine für sich sichtbare Komponente sind —
// sie umschließen die ganze Seite oder alles auf einmal und lassen sich
// nicht wie ein Knopf oder ein Link isoliert beurteilen. Die Liste ist
// bewusst benannt statt implizit, damit ihre Kürze als Absicht erkennbar
// ist und nicht als vergessener Fall.
var keineKomponente = map[string]bool{
	"html": true,
	"body": true,
	"*":    true,
}

// entferneKeyframes schneidet @keyframes-Blöcke vollständig heraus, bevor
// elementSelektoren läuft. Ohne das läse die Regelkopf-Suche das "to" in
// "@keyframes spin { to { transform: … } }" als Element-Selektor "to" —
// das gibt es in HTML nicht, die Prüfung dürfte es also nie verlangen.
// Die Klammertiefe wird gezählt, weil @media auf dieselbe Art verschachtelt
// ist und ein einfaches "bis zur nächsten schließenden Klammer" bei einem
// verschachtelten Block zu früh abschneiden würde.
func entferneKeyframes(css string) string {
	for {
		i := strings.Index(css, "@keyframes")
		if i < 0 {
			return css
		}
		open := strings.Index(css[i:], "{")
		if open < 0 {
			return css
		}
		start := i + open
		tiefe := 0
		ende := len(css)
		for j := start; j < len(css); j++ {
			switch css[j] {
			case '{':
				tiefe++
			case '}':
				tiefe--
				if tiefe == 0 {
					ende = j + 1
					goto fertig
				}
			}
		}
	fertig:
		css = css[:i] + css[ende:]
	}
}

// elementSelektoren sammelt die Element-Selektoren, die ein CSS-Text
// gestaltet — "a" aus "a { … }" ebenso wie aus "a:visited { … }" oder aus
// "thead th { … }" (letztes Glied der Nachfolger-Kette). Klassen- (".x"),
// ID- (#x) und reine Pseudo-Klassen-Selektoren (":focus-visible") liefern
// keinen Treffer, weil sie kein HTML-Element benennen, sondern nur eine
// Markierung oder einen Zustand.
func elementSelektoren(css string) map[string]bool {
	css = kommentarRe.ReplaceAllString(css, "")
	css = entferneKeyframes(css)

	gefunden := make(map[string]bool)
	for _, m := range regelKopfRe.FindAllStringSubmatch(css, -1) {
		kopf := strings.TrimSpace(m[1])
		if kopf == "" || strings.HasPrefix(kopf, "@") {
			continue
		}
		for _, selektor := range strings.Split(kopf, ",") {
			teile := strings.Fields(strings.TrimSpace(selektor))
			if len(teile) == 0 {
				continue
			}
			letztes := teile[len(teile)-1]
			// Am ersten Sonderzeichen abschneiden: "input:disabled" -> "input",
			// "input[aria-disabled=…]" -> "input", ".checkbox-row" -> "".
			ende := len(letztes)
			for i, r := range letztes {
				if r == ':' || r == '[' || r == '.' || r == '#' || r == '*' || r == '&' {
					ende = i
					break
				}
			}
			name := letztes[:ende]
			if elementNameRe.MatchString(name) {
				gefunden[name] = true
			}
		}
	}
	return gefunden
}

// benutzteElemente extrahiert die Tags, die in einem HTML-Text tatsächlich
// vorkommen — Gegenstück zu benutzteKlassen, nur für Elemente statt
// Klassen. Ein Vorkommen mit class="skip-link" zählt nicht: es zeigt die
// Komponente nicht, sondern ist reine Gerüststruktur für die
// Tastaturnavigation (siehe istSprunglink). Andere Tags derselben Art an
// anderer Stelle im Dokument zählen weiterhin normal.
func benutzteElemente(html string) map[string]bool {
	elemente := make(map[string]bool)
	for _, m := range tagRe.FindAllStringSubmatch(html, -1) {
		if istSprunglink(m[2]) {
			continue
		}
		elemente[strings.ToLower(m[1])] = true
	}
	return elemente
}

func TestElementSelektorenErkenntFehlendesElement(t *testing.T) {
	css := `
/* Ohne eigene Regel greift die Browser-Vorgabe. */
a {
  color: red;
}

a:visited {
  color: red;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.card {
  padding: 1rem;
}
`
	elemente := elementSelektoren(css)
	if !elemente["a"] {
		t.Error(`elementSelektoren hat "a" nicht gefunden`)
	}
	if elemente["to"] {
		t.Error(`elementSelektoren hat den @keyframes-Schritt "to" fälschlich als Element gewertet`)
	}
	if elemente["card"] {
		t.Error(`elementSelektoren hat die Klasse ".card" fälschlich als Element gewertet`)
	}

	ohneLink := `<div class="card"><p>Kein Verweis hier.</p></div>`
	mitLink := `<div class="card"><p><a href="#x">Verweis</a></p></div>`

	if benutzteElemente(ohneLink)["a"] {
		t.Error(`benutzteElemente hat in HTML ohne <a> fälschlich "a" gemeldet`)
	}
	if !benutzteElemente(mitLink)["a"] {
		t.Error(`benutzteElemente hat <a> in HTML mit einem Verweis nicht gefunden`)
	}
}

// TestBenutzteElementeZaehltSprunglinkNicht belegt genau den Fall, der die
// erste Fassung der Prüfung unbemerkt bestehen ließ: eine Demo-Seite, auf
// der der einzige <a> der unsichtbare Sprunglink ist. benutzteElemente muss
// "a" hier als nicht gezeigt melden — sonst besteht der Zustand, den der
// ursprüngliche Befund (kein sichtbarer Link auf der Referenzseite) genau
// beschreibt, die Prüfung unentdeckt.
func TestBenutzteElementeZaehltSprunglinkNicht(t *testing.T) {
	nurSprunglink := `<body>
  <a class="skip-link" href="#inhalt">Sprunglink (nur bei Fokus sichtbar)</a>
  <main id="inhalt"><p>Kein weiterer Verweis auf der Seite.</p></main>
</body>`

	if benutzteElemente(nurSprunglink)["a"] {
		t.Error(`benutzteElemente hat ein <a class="skip-link"> fälschlich als gezeigten Link gewertet — der Sprunglink ist per Gestaltung unsichtbar und taugt nicht als Beleg`)
	}

	mitEchtemLinkUndSprunglink := `<body>
  <a class="skip-link" href="#inhalt">Sprunglink</a>
  <main id="inhalt"><p><a href="/designsystem.css">Stylesheet</a></p></main>
</body>`

	if !benutzteElemente(mitEchtemLinkUndSprunglink)["a"] {
		t.Error(`benutzteElemente hat einen echten Link neben dem Sprunglink nicht gefunden`)
	}
}

// TestDemoZeigtJedesElement schließt die Lücke, die TestDemoZeigtJedeKomponente
// lässt: jener Test leitet seine erwartete Liste aus Klassennamen ab und
// übersieht deshalb Komponenten, die base.css über einen bloßen
// Element-Selektor gestaltet — a, a:visited, table, th, td, input, select,
// textarea, label. Bei den Formularelementen und der Tabelle blieb das
// folgenlos, weil die Demo sie ohnehin zeigt; beim Link fiel es der
// Schlussprüfung auf, weil er als einzige Komponente fehlte, die genau die
// im Befund behobene Regel überhaupt sichtbar macht.
func TestDemoZeigtJedesElement(t *testing.T) {
	angeboten := elementSelektoren(string(designsystem.BaseCSS()))
	for name := range keineKomponente {
		delete(angeboten, name)
	}

	benutzt := benutzteElemente(string(indexHTML))
	for element := range angeboten {
		if !benutzt[element] {
			t.Errorf("Die Demo-Seite zeigt kein <%s> — das Element wäre nirgends vor dem Einsatz zu sehen", element)
		}
	}
}

func TestDemoZeigtJedeKomponente(t *testing.T) {
	// Welche Klassen base.css anbietet, steht in base.css — nicht in einer
	// Liste hier, die beim nächsten Zuwachs vergessen würde.
	angeboten := map[string]bool{}
	for _, zeile := range strings.Split(string(designsystem.BaseCSS()), "\n") {
		if m := klassenRe.FindStringSubmatch(strings.TrimSpace(zeile)); m != nil {
			angeboten[m[1]] = true
		}
	}
	// Zustandsklassen erscheinen nur zusammen mit ihrer Grundklasse.
	for _, nur := range []string{"active"} {
		delete(angeboten, nur)
	}

	demo := string(indexHTML)
	benutzt := benutzteKlassen(demo)
	for klasse := range angeboten {
		if !benutzt[klasse] {
			t.Errorf("Die Demo-Seite zeigt .%s nicht — die Klasse wäre nirgends vor dem Einsatz zu sehen", klasse)
		}
	}
}

// TestIndexBindetStylesheetEinUndSchummeltNicht hält die Grundvoraussetzung
// der Referenzseite fest: sie muss das echte Design-System-Stylesheet
// einbinden und darf keinen eigenen <style>-Block mitbringen. Ohne diese
// Prüfung ließe sich jeder Test hier durch einen lokal „reparierten"
// <style>-Block oder das stille Entfernen des <link> unterlaufen — die
// Demo-Seite ist laut Entwurf für die Mehrzahl der Komponenten die einzige
// Prüfstelle, und eine Seite, die ihr eigenes Stylesheet flickt, gäbe eine
// gute Abnahme für ein tatsächlich kaputtes designsystem.css.
func TestIndexBindetStylesheetEinUndSchummeltNicht(t *testing.T) {
	html := string(indexHTML)
	if !strings.Contains(html, `<link rel="stylesheet" href="/designsystem.css" />`) &&
		!strings.Contains(html, `<link rel="stylesheet" href="/designsystem.css">`) {
		t.Error("demo/index.html enthält nicht <link rel=\"stylesheet\" href=\"/designsystem.css\">")
	}
	if strings.Contains(html, "<style") {
		t.Error("demo/index.html enthält einen eigenen <style>-Block — das würde ein kaputtes designsystem.css lokal überdecken")
	}
}

// TestSprunglinkIstErstesFokussierbaresElement prüft, dass der Sprunglink
// unmittelbar hinter dem öffnenden <body> steht — vor <main>, nicht danach
// und nicht innerhalb davon. Ein Sprunglink, der irgendwo später im
// Dokument steht (etwa am Ende von <main>, dem Bereich, den er überspringen
// soll), ist funktionslos: er muss das erste fokussierbare Element im
// Dokument sein, sonst hat Tab ihn längst überholt, bevor er greift.
func TestSprunglinkIstErstesFokussierbaresElement(t *testing.T) {
	html := string(indexHTML)
	iBody := strings.Index(html, "<body>")
	iSkip := strings.Index(html, `class="skip-link"`)
	iMain := strings.Index(html, "<main")
	if iBody < 0 || iSkip < 0 || iMain < 0 {
		t.Fatal("demo/index.html enthält nicht <body>, .skip-link und <main wie erwartet")
	}
	if !(iBody < iSkip && iSkip < iMain) {
		t.Errorf("Reihenfolge ist %d (<body>) < %d (.skip-link) < %d (<main) nicht erfüllt — der Sprunglink muss zwischen <body> und <main stehen", iBody, iSkip, iMain)
	}
}

func TestHandlerLiefertSeiteUndCSS(t *testing.T) {
	srv := httptest.NewServer(Handler())
	defer srv.Close()

	for _, f := range []struct {
		pfad string
		typ  string
		teil string
	}{
		{"/", "text/html; charset=utf-8", "<!doctype html>"},
		{"/designsystem.css", "text/css; charset=utf-8", "--control-line:"},
	} {
		res, err := http.Get(srv.URL + f.pfad)
		if err != nil {
			t.Fatalf("GET %s: %v", f.pfad, err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Errorf("GET %s: Status %d, erwartet 200", f.pfad, res.StatusCode)
		}
		if got := res.Header.Get("Content-Type"); got != f.typ {
			t.Errorf("GET %s: Content-Type %q, erwartet %q", f.pfad, got, f.typ)
		}
		buf := make([]byte, 4096)
		n, _ := res.Body.Read(buf)
		if !strings.Contains(strings.ToLower(string(buf[:n])), strings.ToLower(f.teil)) {
			t.Errorf("GET %s: %q kommt im Anfang der Antwort nicht vor", f.pfad, f.teil)
		}
	}
}
