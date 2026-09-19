package designsystem

import (
	"regexp"
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
//
// Kommentare werden vor der Suche entfernt (removeComments, siehe
// css_test.go — dort inzwischen auch für //-Zeilenkommentare erweitert, die
// CSS nicht kennt, JavaScript aber schon). Ohne das ist die Prüfung durch
// Kommentare erfüllbar: sechs der zehn Muster (AbortController, mousedown,
// aria-activedescendant, Escape, id: null, destroy) stehen in
// js/combobox.js sowohl im Code als auch im erklärenden Kommentar darüber
// — wird der ganze Text durchsucht, meldet die Funktion eine Vorkehrung
// als vorhanden, obwohl der zugehörige Code vollständig entfernt wurde.
func fehlendeVorkehrungen(js string) []string {
	kommentarlos := removeComments(js)
	var fehlt []string
	for muster := range vorkehrungen {
		if muster == "destroy" {
			if !destroyRaeumtAuf(kommentarlos) {
				fehlt = append(fehlt, muster)
			}
			continue
		}
		if !strings.Contains(kommentarlos, muster) {
			fehlt = append(fehlt, muster)
		}
	}
	return fehlt
}

// destroyRaeumtAuf prüft die Vorkehrung "destroy" stärker als reine
// Wortanwesenheit: die Methode heißt zwangsläufig "destroy" — dieses Wort
// bliebe selbst dann im Code stehen, wenn ihr Rumpf komplett geleert würde
// (belegt: "destroy() { … }" auf "destroy() {}" geleert → der alte,
// wortbasierte Test blieb grün). Geprüft wird deshalb der RUMPF der
// destroy()-Methode: er muss tatsächlich aufräumen (mindestens
// clearTimeout und einen Abbruch über .abort() enthalten), nicht nur
// existieren.
var destroyMethodeRe = regexp.MustCompile(`destroy\(\)\s*\{([^}]*)\}`)

func destroyRaeumtAuf(js string) bool {
	m := destroyMethodeRe.FindStringSubmatch(js)
	if m == nil {
		return false
	}
	rumpf := m[1]
	return strings.Contains(rumpf, "clearTimeout") && strings.Contains(rumpf, ".abort()")
}

func TestJSEnthaeltDieSiebenVorkehrungen(t *testing.T) {
	for _, muster := range fehlendeVorkehrungen(string(JS())) {
		t.Errorf("JS() enthält %q nicht (%s)", muster, vorkehrungen[muster])
	}
}

// TestFehlendeVorkehrungErkanntBelegt macht wahr, was die vorige Prüfung nur
// behauptet: dass sie tatsächlich anschlägt, wenn eine Vorkehrung fehlt. Ein
// Test, der nie rot war, beweist nichts.
//
// Die Mutation trifft ausschließlich den CODE, nie den erklärenden
// Kommentar darüber — sonst bliebe das Muster über den Kommentar hinweg
// weiter auffindbar, und der Test bewiese nur, dass ein Kommentar erkannt
// wird, nicht dass fehlender CODE erkannt wird. Betroffen sind:
//   - AbortController: die einzige Codestelle, die den Namen nennt
//     ("laufend = new AbortController()"), wird entfernt; der erklärende
//     Kommentar bei "Die Entprellung ist nicht verhandelbar…" bleibt
//     unverändert stehen.
//   - destroy: der Methodenkörper wird geleert (destroy() {}); Name und
//     der erklärende Kommentar darüber bleiben stehen.
func TestFehlendeVorkehrungErkanntBelegt(t *testing.T) {
	js := string(JS())

	verstuemmelt := strings.Replace(js,
		"laufend = new AbortController()", "laufend = null", 1)
	verstuemmelt = strings.Replace(verstuemmelt,
		"destroy() {\n      clearTimeout(timer)\n      laufend?.abort()\n    },",
		"destroy() {},", 1)

	// Voraussetzung der Mutationsprobe: beide Ersetzungen müssen tatsächlich
	// gegriffen haben, sonst prüft der Rest dieses Tests nichts.
	if strings.Contains(verstuemmelt, "new AbortController()") {
		t.Fatal("Mutation hat AbortController im Code nicht entfernt — Testaufbau prüft nicht das Vorgesehene")
	}
	if strings.Contains(verstuemmelt, "clearTimeout(timer)\n      laufend?.abort()\n    },") {
		t.Fatal("Mutation hat den destroy()-Körper nicht geleert — Testaufbau prüft nicht das Vorgesehene")
	}
	// Die erklärenden Kommentare bleiben absichtlich stehen — genau das ist
	// der Fall, den der alte, kommentarblinde Test überging.
	if !strings.Contains(verstuemmelt, "AbortController") {
		t.Fatal("Testaufbau fehlerhaft: der Kommentar, der \"AbortController\" nennt, wurde versehentlich mitentfernt")
	}
	if !strings.Contains(verstuemmelt, "destroy()") {
		t.Fatal("Testaufbau fehlerhaft: die Erwähnung von destroy() im Kommentar wurde versehentlich mitentfernt")
	}

	fehlt := fehlendeVorkehrungen(verstuemmelt)
	gefunden := map[string]bool{}
	for _, m := range fehlt {
		gefunden[m] = true
	}
	if !gefunden["AbortController"] {
		t.Error("fehlendeVorkehrungen meldet ein aus dem Code entferntes AbortController nicht als fehlend, obwohl es nur noch im Kommentar steht")
	}
	if !gefunden["destroy"] {
		t.Error("fehlendeVorkehrungen meldet ein geleertes destroy() nicht als fehlend, obwohl der Name nur noch im Kommentar steht")
	}
}

// TestKommentierteVorkehrungWirdNichtAlsAnwesendGezaehlt belegt das
// eigentliche Kernstück von Teil A direkt: ein Text, der ein Muster NUR im
// Kommentar enthält (kein Code), muss als fehlend gelten. Das ist die
// Mutationsprobe aus dem Schlussprüfungsbefund nachgebaut — vorher lieferte
// removeComments die Kommentare unverändert mit, dieser Test hätte fälsch-
// lich "vorhanden" gemeldet.
func TestKommentierteVorkehrungWirdNichtAlsAnwesendGezaehlt(t *testing.T) {
	text := "// AbortController wird hier nur erwähnt, nicht verwendet\nconst x = 1\n"
	fehlt := fehlendeVorkehrungen(text)
	gefunden := map[string]bool{}
	for _, m := range fehlt {
		gefunden[m] = true
	}
	if !gefunden["AbortController"] {
		t.Error("ein Muster, das nur im Kommentar steht, wird fälschlich als im Code vorhanden gewertet")
	}
}
