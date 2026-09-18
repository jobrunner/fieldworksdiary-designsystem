package designsystem

import (
	"strings"
	"testing"
)

// TestBaseCSSEnthaeltAkkordeonRegeln prüft die Zusagen aus Teil A: das
// Akkordeon-Muster (<details>/<summary>) bekommt seine Gestaltung aus
// base.css, ohne dass ein Dienst sie selbst nachbauen muss. Geprüft wird
// insbesondere, dass das Vorgabedreieck des Browsers abgeschaltet ist —
// laut Aufgabenstellung leicht zu vergessen und sähe sonst neben dem
// Chevron aus dem Symbolpaket doppelt aus.
// akkordeonMuster sind die Textmuster, an denen sich die Zusagen aus Teil A
// erkennen lassen — als map von Muster auf Begründung, damit ein
// Fehlschlag sofort sagt, welche Zusage fehlt.
var akkordeonMuster = map[string]string{
	".akkordeon":                       "Rahmen und Radius des <details>",
	".akkordeon > summary":             "Kopfzeile, Mindesthöhe, sichtbarer Fokus",
	".akkordeon-inhalt":                "Innenabstand und Trennlinie",
	".akkordeon[open] > summary .icon": "der Chevron dreht sich beim Öffnen",
	"list-style: none":                 "das Vorgabedreieck in Firefox ist abgeschaltet",
	"::-webkit-details-marker":         "das Vorgabedreieck in Chrome/Safari ist abgeschaltet",
	"rotate(180deg)":                   "der Chevron dreht tatsächlich, nicht nur der Selektor greift",
}

func fehlendeAkkordeonMuster(css string) []string {
	var fehlt []string
	for muster := range akkordeonMuster {
		if !strings.Contains(css, muster) {
			fehlt = append(fehlt, muster)
		}
	}
	return fehlt
}

func TestBaseCSSEnthaeltAkkordeonRegeln(t *testing.T) {
	for _, muster := range fehlendeAkkordeonMuster(string(BaseCSS())) {
		t.Errorf("base.css enthält %q nicht (%s)", muster, akkordeonMuster[muster])
	}
}

// TestFehlendeAkkordeonRegelWirdErkannt belegt, dass die vorige Prüfung
// tatsächlich anschlägt, wenn eine Regel fehlt — ein Test, der nie rot war,
// beweist nichts. Geprüft an CSS-Text ohne die Marker-Abschaltung für
// Chrome/Safari, dem Fall, den die Aufgabenstellung als "leicht zu
// vergessen" nennt.
func TestFehlendeAkkordeonRegelWirdErkannt(t *testing.T) {
	verstuemmelt := strings.ReplaceAll(string(BaseCSS()), "::-webkit-details-marker", "")

	fehlt := fehlendeAkkordeonMuster(verstuemmelt)
	gefunden := false
	for _, m := range fehlt {
		if m == "::-webkit-details-marker" {
			gefunden = true
		}
	}
	if !gefunden {
		t.Error("fehlendeAkkordeonMuster meldet eine entfernte Marker-Abschaltung nicht als fehlend")
	}
}

// TestBaseCSSEnthaeltComboboxRegeln prüft die Gestaltung aus Teil B: Hülle,
// Vorschlagsliste und Einträge, sowie die Hervorhebung über aria-selected
// statt über eine eigene Klasse — Darstellung und Ansage dürfen nicht
// auseinanderlaufen können.
func TestBaseCSSEnthaeltComboboxRegeln(t *testing.T) {
	css := string(BaseCSS())

	for _, muster := range []string{
		".combobox",
		".combobox-liste",
		".combobox-liste .option",
		`[aria-selected="true"]`,
		"position: relative",
		"overflow-y: auto",
	} {
		if !strings.Contains(css, muster) {
			t.Errorf("base.css enthält %q nicht", muster)
		}
	}
}
