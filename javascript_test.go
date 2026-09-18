package designsystem

import (
	"strings"
	"testing"
)

// TestJSLiefertNichtLeeresModul prüft die Grundvoraussetzung: JS() liefert
// überhaupt etwas, und das Gelieferte ist das erwartete Modul, nicht ein
// leerer oder vertauschter Puffer.
func TestJSLiefertNichtLeeresModul(t *testing.T) {
	js := string(JS())
	if len(js) == 0 {
		t.Fatal("JS() liefert einen leeren Puffer")
	}
	if !strings.Contains(js, "export function mountCombobox") {
		t.Error(`JS() enthält "export function mountCombobox" nicht`)
	}
}

// TestJSLiefertKopie belegt wie TokensCSS/BaseCSS, dass JS() eine Kopie
// liefert: ein Aufrufer, der das Ergebnis in-place verändert, darf den
// eingebetteten Puffer nicht für nachfolgende Aufrufer beschädigen.
func TestJSLiefertKopie(t *testing.T) {
	erste := JS()
	if len(erste) == 0 {
		t.Fatal("JS() liefert einen leeren Puffer, kann Veränderung nicht prüfen")
	}
	erste[0] = '!'
	zweite := JS()
	if zweite[0] == '!' {
		t.Error("JS() liefert eine Referenz auf den eingebetteten Puffer statt einer Kopie — Veränderung durch einen Aufrufer wirkt auf alle weiteren Aufrufe durch")
	}
}

// TestJSEnthaeltKeinenExternenImport belegt die Zusage, dass das Modul ohne
// Bündler und ohne Abhängigkeiten läuft: ein "import" von außerhalb (aus
// einer anderen Datei oder einem Paket) würde das im Browser ohne Bündler
// scheitern lassen oder stillschweigend eine fremde Abhängigkeit
// einschleusen.
func TestJSEnthaeltKeinenExternenImport(t *testing.T) {
	js := string(JS())
	for _, zeile := range strings.Split(js, "\n") {
		trimmed := strings.TrimSpace(zeile)
		if strings.HasPrefix(trimmed, "import ") || strings.HasPrefix(trimmed, "import(") {
			t.Errorf("JS() enthält einen import: %q — das Modul soll ohne Bündler und ohne Abhängigkeiten laufen", trimmed)
		}
	}
}

// TestJSEnthaeltDieSiebenVorkehrungen prüft von Go aus, was ohne
// JavaScript-Testlauf prüfbar ist: dass die sieben Vorkehrungen aus
// Expertus' Vorlage (Entprellung, AbortController, mousedown statt click,
// aria-expanded/aria-activedescendant/aria-selected, Tastaturbedienung,
// Freitext mit id: null, destroy()) im Quelltext anwesend sind. Das ist eine
// schwächere Zusage als ein echter Testlauf: geprüft wird die ANWESENHEIT
// der Vorkehrungen, nicht ihre WIRKUNG — ob die Entprellung tatsächlich
// verzögert, ob destroy() tatsächlich jeden Zeitgeber trifft, bliebe auch
// nach diesem Test offen. Dieses Modul hat keinen JavaScript-Testlauf und
// soll auch keinen einführen (siehe Aufgabenstellung).
// vorkehrungen listet die Textmuster, an denen sich die sieben aus Expertus
// übernommenen Vorkehrungen erkennen lassen. Eine map von Muster auf
// Begründung, damit ein Fehlschlag sofort sagt, WELCHE Vorkehrung fehlt und
// WARUM sie gefordert ist — nicht nur, dass irgendein String nicht vorkam.
var vorkehrungen = map[string]string{
	"setTimeout":            "Entprellung vor dem Netzaufruf",
	"AbortController":       "Abbruch überholter Anfragen",
	"mousedown":             "Auswahl per mousedown statt click",
	"aria-expanded":         "Zustand \"offen\" am Eingabefeld",
	"aria-activedescendant": "hervorgehobener Eintrag am Eingabefeld",
	"aria-selected":         "Ansage der Hervorhebung an den Einträgen",
	"ArrowDown":             "Tastaturbedienung: Hervorhebung bewegen",
	"Escape":                "Tastaturbedienung: schließen und abbrechen",
	"id: null":              "Freitext bleibt zulässig",
	"destroy":               "Aufräumen von Zeitgeber und laufender Anfrage",
}

// fehlendeVorkehrungen prüft von Go aus, was ohne JavaScript-Testlauf
// prüfbar ist: dass jede der sieben Vorkehrungen aus Expertus' Vorlage
// (Entprellung, AbortController, mousedown statt click, die drei
// aria-Attribute, Tastaturbedienung, Freitext mit id: null, destroy()) im
// Quelltext anwesend ist. Das ist eine schwächere Zusage als ein echter
// Testlauf: geprüft wird die ANWESENHEIT der Vorkehrungen, nicht ihre
// WIRKUNG — ob die Entprellung tatsächlich verzögert oder destroy()
// tatsächlich jeden Zeitgeber trifft, bliebe auch danach offen. Dieses
// Modul hat keinen JavaScript-Testlauf und soll auch keinen einführen
// (siehe Aufgabenstellung).
func fehlendeVorkehrungen(js string) []string {
	var fehlt []string
	for muster := range vorkehrungen {
		if !strings.Contains(js, muster) {
			fehlt = append(fehlt, muster)
		}
	}
	return fehlt
}

func TestJSEnthaeltDieSiebenVorkehrungen(t *testing.T) {
	for _, muster := range fehlendeVorkehrungen(string(JS())) {
		t.Errorf("JS() enthält %q nicht (%s)", muster, vorkehrungen[muster])
	}
}

// TestFehlendeVorkehrungErkanntBelegt macht wahr, was die vorige Prüfung nur
// behauptet: dass sie tatsächlich anschlägt, wenn eine Vorkehrung fehlt.
// Ein Test, der nie rot war, beweist nichts — deshalb hier derselbe
// Erkennungsweg (fehlendeVorkehrungen), einmal absichtlich gegen einen Text
// ohne AbortController und ohne destroy geführt.
func TestFehlendeVorkehrungErkanntBelegt(t *testing.T) {
	verstuemmelt := strings.ReplaceAll(string(JS()), "AbortController", "XxxEntfernt")
	verstuemmelt = strings.ReplaceAll(verstuemmelt, "destroy", "XxxEntfernt2")

	fehlt := fehlendeVorkehrungen(verstuemmelt)
	gefunden := map[string]bool{}
	for _, m := range fehlt {
		gefunden[m] = true
	}
	if !gefunden["AbortController"] {
		t.Error("fehlendeVorkehrungen meldet ein entferntes AbortController nicht als fehlend")
	}
	if !gefunden["destroy"] {
		t.Error("fehlendeVorkehrungen meldet ein entferntes destroy nicht als fehlend")
	}
}
