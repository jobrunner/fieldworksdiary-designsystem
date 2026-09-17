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
