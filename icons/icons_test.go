package icons

import (
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/jobrunner/fieldworksdiary-designsystem/internal/farbe"
)

func TestHuelleFolgtDerBauregel(t *testing.T) {
	got := string(huelle(`<path d="M6 9l6 6 6-6"/>`, false))

	for _, pflicht := range []string{
		`viewBox="0 0 24 24"`,
		`fill="none"`,
		`stroke="currentColor"`,
		`stroke-width="2"`,
		`stroke-linecap="round"`,
		`stroke-linejoin="round"`,
		`aria-hidden="true"`,
		`focusable="false"`,
		`<path d="M6 9l6 6 6-6"/>`,
	} {
		if !contains(got, pflicht) {
			t.Errorf("die Hülle enthält %q nicht:\n%s", pflicht, got)
		}
	}
	// Feste Maße würden die Größe am Einsatzort festnageln; sie soll über
	// CSS bestimmt werden. Regex erfasst Attribute nach Leerraum, Tabulator,
	// Zeilenumbruch — aber nicht stroke-width/stroke-height (Bindestrich ist
	// kein Leerraum).
	for name, re := range map[string]*regexp.Regexp{
		"width":  regexp.MustCompile(`(^|[\s])width\s*=`),
		"height": regexp.MustCompile(`(^|[\s])height\s*=`),
	} {
		if re.MatchString(got) {
			t.Errorf("die Hülle enthält Attribut %q — die Größe gehört ins CSS:\n%s", name, got)
		}
	}
	// Im Strich-Zweig muss fill="none" vorhanden sein und
	// fill="currentColor" oder stroke="none" darf nicht vorhanden sein.
	if !contains(got, `fill="none"`) {
		t.Errorf("Strich-Hülle enthält fill=none nicht:\n%s", got)
	}
	if contains(got, `fill="currentColor"`) {
		t.Errorf("Strich-Hülle enthält versehentlich fill=currentColor:\n%s", got)
	}
	if contains(got, `stroke="none"`) {
		t.Errorf("Strich-Hülle enthält versehentlich stroke=none:\n%s", got)
	}
}

func TestGefuellteHuelleFuerFlaechensymbole(t *testing.T) {
	// Die Mondphasen sind die einzige Ausnahme von der Strichzeichnung:
	// bei ihnen trägt die Flächenaufteilung die Aussage.
	got := string(huelle(`<circle cx="12" cy="12" r="9"/>`, true))
	if !contains(got, `fill="currentColor"`) {
		t.Errorf("gefüllte Hülle nutzt kein fill=currentColor:\n%s", got)
	}
	if contains(got, `fill="none"`) {
		t.Errorf("gefüllte Hülle enthält noch fill=none:\n%s", got)
	}
	if contains(got, `stroke="currentColor"`) {
		t.Errorf("gefüllte Hülle enthält versehentlich stroke=currentColor:\n%s", got)
	}
}

func TestIconString(t *testing.T) {
	// Icon.String() wandelt den Typ in seinen String-Wert um.
	icon := huelle(`<path d="M12 12"/>`, false)
	result := icon.String()
	if !contains(result, `viewBox="0 0 24 24"`) {
		t.Errorf("Icon.String() gab kein gültiges SVG zurück: %s", result)
	}
}

func TestJedesSymbolFolgtDerBauregel(t *testing.T) {
	if len(Alle()) == 0 {
		t.Fatal("Alle() ist leer — dann prüft dieser Test nichts")
	}
	for name, icon := range Alle() {
		s := string(icon)
		t.Run(name, func(t *testing.T) {
			for _, pflicht := range []string{
				`viewBox="0 0 24 24"`,
				`aria-hidden="true"`,
				`focusable="false"`,
				`currentColor`,
			} {
				if !contains(s, pflicht) {
					t.Errorf("%s enthält %q nicht", name, pflicht)
				}
			}
			// Prüfe freistehende width/height nur auf dem SVG-Element selbst
			// (zwischen <svg und dem ersten >), nicht auf inneren Elementen.
			svgOpenRe := regexp.MustCompile(`<svg\s+[^>]*?>`)
			svgOpen := svgOpenRe.FindString(s)
			if svgOpen != "" {
				for _, attr := range []string{"width", "height"} {
					if regexp.MustCompile(`(^|[\s])` + attr + `\s*=`).MatchString(svgOpen) {
						t.Errorf("%s hat %q auf dem SVG-Element — Größe gehört ins CSS", name, attr)
					}
				}
			}
			// Farben in Klartext sind überall verboten — mit derselben
			// Strenge wie in base.css (css_test.go): nicht nur "#",
			// "rgb(", "hsl(", sondern auch benannte Farben wie "black"
			// und moderne Funktionen wie oklch(). Belegt: stroke="black"
			// kam an der alten, dreiteiligen Prüfung unbemerkt vorbei.
			for _, wert := range farbe.Funde(s) {
				t.Errorf("%s enthält Farbe %q — Farbe kommt aus currentColor", name, wert)
			}
			if !contains(s, `<svg`) || !contains(s, `</svg>`) {
				t.Errorf("%s ist kein vollständiges SVG", name)
			}
		})
	}
}

// TestHTMLLiefertTemplateHTML belegt die Zusage aus icons.go: HTML()
// liefert template.HTML statt einem gewöhnlichen string, damit
// html/template das Symbol als Markup einsetzt statt es zu maskieren.
func TestHTMLLiefertTemplateHTML(t *testing.T) {
	icon := Standort()
	got := icon.HTML()
	if string(got) != string(icon) {
		t.Errorf("icon.HTML() = %q, erwartet %q", got, icon)
	}
}

// TestOhneHTMLWuerdeHTMLTemplateDasSymbolMaskieren belegt das eigentliche
// Problem, das HTML() löst: schreibt eine html/template-Vorlage
// {{.Symbol}} statt {{.Symbol.HTML}}, bekommt sie sichtbaren
// SVG-Quelltext statt eines Symbols, weil Icon für html/template
// gewöhnlicher Text ist.
func TestOhneHTMLWuerdeHTMLTemplateDasSymbolMaskieren(t *testing.T) {
	icon := Standort()
	tplOhne := template.Must(template.New("ohne").Parse(`<button>{{.}}</button>`))
	var b strings.Builder
	if err := tplOhne.Execute(&b, icon); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(b.String(), "<svg") {
		t.Fatalf("Testaufbau fehlerhaft: {{.}} auf Icon sollte SVG maskieren, tat es hier nicht: %s", b.String())
	}
	if !strings.Contains(b.String(), "&lt;svg") {
		t.Errorf("{{.}} auf Icon sollte den maskierten SVG-Quelltext zeigen (Beleg für das Problem), zeigt aber: %s", b.String())
	}

	tplMit := template.Must(template.New("mit").Parse(`<button>{{.HTML}}</button>`))
	var b2 strings.Builder
	if err := tplMit.Execute(&b2, icon); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b2.String(), "<svg") {
		t.Errorf("{{.HTML}} sollte das Symbol als Markup einsetzen, tat es hier nicht: %s", b2.String())
	}
}

// TestMitKlasseSetztKlasseAufsSvgElement belegt icons.MitKlasse(): sie
// erspart Dienste die Zeichenketten-Chirurgie, die zuvor in demo.go stand
// (strings.Replace(svg, "<svg ", …)).
func TestMitKlasseSetztKlasseAufsSvgElement(t *testing.T) {
	icon := Standort()
	got := string(MitKlasse(icon, "icon"))
	if !contains(got, `<svg class="icon" `) {
		t.Errorf("MitKlasse setzt die Klasse nicht auf das <svg>-Element: %s", got)
	}
	// Der restliche Inhalt bleibt unverändert.
	if !contains(got, `viewBox="0 0 24 24"`) {
		t.Errorf("MitKlasse hat den Rest des Symbols verändert: %s", got)
	}
}

// TestMitBeschriftungMachtSymbolAnsagbar belegt icons.MitBeschriftung():
// ein allein stehendes Symbol (etwa in einem Knopf ohne Text) braucht
// role="img" und aria-label statt aria-hidden="true".
func TestMitBeschriftungMachtSymbolAnsagbar(t *testing.T) {
	icon := Standort()
	got := string(MitBeschriftung(icon, "Standort bestimmen"))
	if contains(got, `aria-hidden="true"`) {
		t.Errorf("MitBeschriftung hat aria-hidden=\"true\" nicht entfernt: %s", got)
	}
	if !contains(got, `role="img"`) {
		t.Errorf("MitBeschriftung setzt role=\"img\" nicht: %s", got)
	}
	if !contains(got, `aria-label="Standort bestimmen"`) {
		t.Errorf("MitBeschriftung setzt aria-label nicht mit dem übergebenen Text: %s", got)
	}
}

// TestMitBeschriftungMaskiertDenText belegt, dass ein Anführungszeichen im
// Beschriftungstext nicht aus dem Attribut ausbricht.
func TestMitBeschriftungMaskiertDenText(t *testing.T) {
	icon := Standort()
	got := string(MitBeschriftung(icon, `x" onclick="alert(1)`))
	if contains(got, `onclick="alert(1)"`) {
		t.Errorf("MitBeschriftung maskiert Anführungszeichen im Beschriftungstext nicht — Einschleusung ins Attribut möglich: %s", got)
	}
}

func TestAlleEnthaeltDieBediensymbole(t *testing.T) {
	// Diese Namen sind eine Zusage an die Dienste: sie rufen sie auf.
	for _, name := range []string{
		"standort", "chevron-unten", "kalender", "schliessen", "suche",
		"herunterladen", "kopieren", "haken", "warnung", "information",
		"fehler", "menue",
	} {
		if _, da := Alle()[name]; !da {
			t.Errorf("Alle() kennt %q nicht", name)
		}
	}
}

func TestAlleSammelkeineDuplikate(t *testing.T) {
	// Wenn zwei Teil-Sammlungen denselben Schlüssel haben, geht einer
	// stillschweigend verloren. Diese Prüfung erkennt das, bevor eine
	// Aufgabe folgende Symbole vergisst.

	// Die Teil-Sammlungen — Reihenfolge folgt Alle()
	teile := []map[string]Icon{
		bediensymbole(),
		wettersymbole(),
		himmelssymbole(),
		messwertsymbole(),
	}

	// Summe der Längen aller Teile
	summe := 0
	for _, teil := range teile {
		summe += len(teil)
	}

	// Größe der Gesamtmenge
	gesamt := len(Alle())

	if summe != gesamt {
		// Sammeln, welche Namen vorkommen, um die Meldung zu schärfen
		seheneigene := make(map[string]bool)
		for _, teil := range teile {
			for k := range teil {
				if seheneigene[k] {
					t.Errorf("Name %q kommt in mehreren Teil-Sammlungen vor — "+
						"Alle() überschreibt stillschweigend", k)
				}
				seheneigene[k] = true
			}
		}
		t.Errorf("Duplikate in Teil-Sammlungen: Summe %d, Gesamtmenge %d",
			summe, gesamt)
	}
}

// TestFarbeInKlartextAlsWortWirdErkannt belegt, dass die geschärfte Prüfung
// aus TestJedesSymbolFolgtDerBauregel tatsächlich anschlägt, wenn ein
// Symbol eine Farbe im Klartext trägt — nicht nur bei "#", "rgb(", "hsl(",
// sondern auch bei einem Farbwort. Ein Test, der nie rot war, beweist
// nichts: nachgewiesen ist der reale Fall aus der Schlussprüfung,
// stroke="black" statt stroke="currentColor", der die alte, dreiteilige
// Prüfung unbemerkt passierte.
func TestFarbeInKlartextAlsWortWirdErkannt(t *testing.T) {
	mutiert := strings.Replace(string(Standort()), `stroke="currentColor"`, `stroke="black"`, 1)
	if mutiert == string(Standort()) {
		t.Fatal("Mutation hat stroke=currentColor nicht ersetzt — Testaufbau prüft nicht das Vorgesehene")
	}
	funde := farbe.Funde(mutiert)
	if len(funde) == 0 {
		t.Error(`farbe.Funde erkennt stroke="black" nicht — sollte als Farbliteral gelten`)
	}
}

// TestAlleIstPaarweiseVerschieden prüft, dass keine zwei Symbole
// identisches Markup tragen. Die bisherigen Tests vergleichen nur Anzahlen
// und Schlüsselnamen, nie Zeichnungen — ein falsch verdrahteter
// Karteneintrag (etwa "suche": Standort()) fiel dadurch nicht auf: er
// ändert weder die Anzahl noch die Namen, nur den Inhalt hinter einem der
// Namen.
func TestAlleIstPaarweiseVerschieden(t *testing.T) {
	gesehen := map[Icon]string{}
	for name, icon := range Alle() {
		if vorher, da := gesehen[icon]; da {
			t.Errorf("%q ist identisch gezeichnet wie %q — Copy-Paste-Fehler oder ins falsche Ziel verdrahtet", name, vorher)
			continue
		}
		gesehen[icon] = name
	}
}

// TestFalschVerdrahteterEintragWirdErkannt belegt, dass
// TestAlleIstPaarweiseVerschieden tatsächlich anschlägt: derselbe Fall wie
// im Befund, "suche" trüge das Markup von Standort().
func TestFalschVerdrahteterEintragWirdErkannt(t *testing.T) {
	alle := Alle()
	alle["suche"] = alle["standort"] // Nachbau der Fehlverdrahtung aus dem Befund

	gesehen := map[Icon]string{}
	kollision := false
	for name, icon := range alle {
		if vorher, da := gesehen[icon]; da {
			kollision = true
			_ = vorher
		}
		gesehen[icon] = name
	}
	if !kollision {
		t.Error("die Duplikatsprüfung erkennt eine Fehlverdrahtung (suche == standort) nicht")
	}
}

func TestAlleEnthaeltAlleSymbole(t *testing.T) {
	// Diese Prüfung zählt alle exportierten Icon-Funktionen im Paket
	// und vergleicht mit der Größe von Alle(). Wenn die Zahlen nicht
	// übereinstimmen, wurde eine Symbolfunktion nicht eingetragen.

	// Lese alle .go-Dateien im Package und zähle die Symbolfunktionen
	dir := "."
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Konnte Paketverzeichnis nicht lesen: %v", err)
	}

	var functionCount int
	var functionNames []string
	functionRe := regexp.MustCompile(`(?m)^func ([A-Z]\w*)\(\) Icon \{`)

	for _, entry := range entries {
		if entry.IsDir() || !regexp.MustCompile(`\.go$`).MatchString(entry.Name()) {
			continue
		}
		if entry.Name() == "icons_test.go" {
			continue // Testdatei ignorieren
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Konnte Datei %s nicht lesen: %v", path, err)
		}

		matches := functionRe.FindAllStringSubmatch(string(data), -1)
		for _, match := range matches {
			functionCount++
			if len(match) > 1 {
				functionNames = append(functionNames, match[1])
			}
		}
	}

	actualCount := len(Alle())

	if functionCount != actualCount {
		t.Errorf("Anzahl der Symbolfunktionen (%d) passt nicht zu Alle() (%d). "+
			"Wahrscheinlich wurde eine Symbolfunktion nicht in die Sammlung eingetragen. "+
			"Gefundene Funktionen: %v",
			functionCount, actualCount, functionNames)
	}
}

func contains(h, n string) bool {
	return len(n) <= len(h) && indexOf(h, n) >= 0
}

func indexOf(h, n string) int {
	for i := 0; i+len(n) <= len(h); i++ {
		if h[i:i+len(n)] == n {
			return i
		}
	}
	return -1
}
