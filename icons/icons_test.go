package icons

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
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
			// Farben in Klartext sind überall verboten.
			for _, farbeVerboten := range []string{`#`, `rgb(`, `hsl(`} {
				if contains(s, farbeVerboten) {
					t.Errorf("%s enthält Farbe %q — Farbe kommt aus currentColor", name, farbeVerboten)
				}
			}
			if !contains(s, `<svg`) || !contains(s, `</svg>`) {
				t.Errorf("%s ist kein vollständiges SVG", name)
			}
		})
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
