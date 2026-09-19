package designsystem

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestJSEnthaeltMountKoordinaten prüft die Grundvoraussetzung: JS() liefert
// auch die Koordinateneingabe, nicht nur die Combobox.
func TestJSEnthaeltMountKoordinaten(t *testing.T) {
	js := string(JS())
	if !strings.Contains(js, "export function mountKoordinaten") {
		t.Error(`JS() enthält "export function mountKoordinaten" nicht`)
	}
	if !strings.Contains(js, "export function parseKoordinatenPaar") {
		t.Error(`JS() enthält "export function parseKoordinatenPaar" nicht`)
	}
}

// koordinatenVorkehrungen listet die Textmuster, an denen sich die
// Zusagen aus dem Auftrag erkennen lassen — als map von Muster auf
// Begründung, damit ein Fehlschlag sofort sagt, welche Zusage fehlt.
// Wie bei den Combobox-Vorkehrungen in javascript_test.go wird auf
// kommentarfreiem Text geprüft (siehe fehlendeKoordinatenVorkehrungen),
// sonst wäre die Prüfung durch einen erklärenden Kommentar allein
// erfüllbar, ohne dass der zugehörige Code vorhanden ist.
var koordinatenVorkehrungen = map[string]string{
	"yZuerst":      "Auswahl, welches Feld sichtbar zuerst steht",
	"einzelfeld":   "Umschalten auf ein einzelnes Feld (MGRS)",
	"insertBefore": "die Reihenfolge tauscht im Baum, nicht die id/label-Bindung",
	"idPrefix":     "konfigurierbares id-Präfix gegen Kollisionen zwischen zwei Eingaben",
}

func fehlendeKoordinatenVorkehrungen(js string) []string {
	kommentarlos := removeComments(js)
	var fehlt []string
	for muster := range koordinatenVorkehrungen {
		if !strings.Contains(kommentarlos, muster) {
			fehlt = append(fehlt, muster)
		}
	}
	return fehlt
}

func TestKoordinatenEnthaeltDieVorkehrungen(t *testing.T) {
	for _, muster := range fehlendeKoordinatenVorkehrungen(string(JS())) {
		t.Errorf("JS() enthält %q nicht (%s)", muster, koordinatenVorkehrungen[muster])
	}
}

// TestKoordinatenFehlendeVorkehrungErkanntBelegt macht wahr, was die vorige
// Prüfung nur behauptet: dass sie tatsächlich anschlägt, wenn eine
// Vorkehrung fehlt. Ein Test, der nie rot war, beweist nichts. Getroffen
// wird ausschließlich der CODE, das erklärende "yZuerst tauscht NUR DIE
// POSITION..."-Kommentarargument bleibt unverändert stehen.
func TestKoordinatenFehlendeVorkehrungErkanntBelegt(t *testing.T) {
	js := string(JS())

	verstuemmelt := strings.Replace(js,
		"if (sys.yZuerst) gitter.insertBefore(gruppeY, gruppeX)",
		"if (false) { /* Reihenfolge tauscht hier nicht mehr */ }", 1)
	verstuemmelt = strings.Replace(verstuemmelt,
		"const yZuerst = !!aktuelles.yZuerst\n    ;(yZuerst ? feldY : feldX).value = paar[0]\n    ;(yZuerst ? feldX : feldY).value = paar[1]",
		"feldX.value = paar[0]\n    feldY.value = paar[1]", 1)

	if strings.Contains(removeComments(verstuemmelt), "yZuerst") {
		t.Fatal("Mutation hat yZuerst im Code nicht vollständig entfernt — Testaufbau prüft nicht das Vorgesehene")
	}
	// Der erklärende Kommentar über anwenden() nennt "yZuerst" weiterhin —
	// genau das ist der Fall, den removeComments abfangen muss.
	if !strings.Contains(verstuemmelt, "yZuerst") {
		t.Fatal("Testaufbau fehlerhaft: der erklärende Kommentar, der \"yZuerst\" nennt, wurde versehentlich mitentfernt")
	}

	fehlt := fehlendeKoordinatenVorkehrungen(verstuemmelt)
	gefunden := map[string]bool{}
	for _, m := range fehlt {
		gefunden[m] = true
	}
	if !gefunden["yZuerst"] {
		t.Error("fehlendeKoordinatenVorkehrungen meldet ein aus dem Code entferntes yZuerst nicht als fehlend, obwohl es nur noch im Kommentar steht")
	}
}

// --- Verhaltensprüfung der Zerlegung über einen echten node-Lauf --------
//
// Die fünf Fälle aus dem Auftrag werden hier nicht nur an Textmustern
// erkannt, sondern TATSÄCHLICH ausgeführt: js/koordinaten.js ist ein reines
// ES-Modul ohne DOM-Abhängigkeit für parseKoordinatenPaar, lässt sich also
// mit einem kurzen node-Aufruf importieren und aufrufen — ohne einen neuen
// JavaScript-Testlauf für das ganze Modul einzuführen (siehe
// javascript_syntax_test.go: das ist dieselbe Art Aufruf wie "node
// --check", nur mit "node datei.mjs" statt "--check"). Überspringt sich
// selbst, wenn "node" fehlt.
func TestParseKoordinatenPaarFuenfFaelle(t *testing.T) {
	nodePath, err := exec.LookPath("node")
	if err != nil {
		t.Skip(`"node" ist nicht installiert`)
	}

	tests := []struct {
		name      string
		eingabe   string
		erwartetA string
		erwartetB string
		istPaar   bool
	}{
		{"Semikolon, Komma ist Dezimaltrenner", "35,01;32,67", "35.01", "32.67", true},
		{"Komma trennt, Punkt ist Dezimaltrenner", "35.01, 32.67", "35.01", "32.67", true},
		{"Leerzeichen trennt, Komma ist Dezimaltrenner", "35,01 32,67", "35.01", "32.67", true},
		{"Leerzeichen trennt, Punkt ist Dezimaltrenner", "35.01 32.67", "35.01", "32.67", true},
		{"einzelner Wert mit deutschem Dezimalkomma ist KEIN Paar", "35,016132", "", "", false},
	}

	dir := t.TempDir()
	modulPfad := filepath.Join(dir, "koordinaten.mjs")
	if err := os.WriteFile(modulPfad, JS(), 0o644); err != nil {
		t.Fatalf("konnte Modul nicht schreiben: %v", err)
	}

	eingaben := make([]string, len(tests))
	for i, tc := range tests {
		eingaben[i] = tc.eingabe
	}
	eingabenJSON, err := json.Marshal(eingaben)
	if err != nil {
		t.Fatalf("konnte Eingaben nicht kodieren: %v", err)
	}

	runnerPfad := filepath.Join(dir, "runner.mjs")
	runner := "import { parseKoordinatenPaar } from " +
		strconvQuote(modulPfad) + ";\n" +
		"const eingaben = " + string(eingabenJSON) + ";\n" +
		"console.log(JSON.stringify(eingaben.map((e) => parseKoordinatenPaar(e))));\n"
	if err := os.WriteFile(runnerPfad, []byte(runner), 0o644); err != nil {
		t.Fatalf("konnte Runner nicht schreiben: %v", err)
	}

	out, err := exec.Command(nodePath, runnerPfad).CombinedOutput()
	if err != nil {
		t.Fatalf("node-Lauf ist fehlgeschlagen: %v\n%s", err, out)
	}

	var ergebnisse []*[2]string
	if err := json.Unmarshal(out, &ergebnisse); err != nil {
		t.Fatalf("konnte Ausgabe nicht als JSON lesen: %v\nAusgabe: %s", err, out)
	}
	if len(ergebnisse) != len(tests) {
		t.Fatalf("%d Ergebnisse, erwartet %d", len(ergebnisse), len(tests))
	}

	for i, tc := range tests {
		got := ergebnisse[i]
		if !tc.istPaar {
			if got != nil {
				t.Errorf("%s: parseKoordinatenPaar(%q) = %v, erwartet null (kein Paar)", tc.name, tc.eingabe, *got)
			}
			continue
		}
		if got == nil {
			t.Errorf("%s: parseKoordinatenPaar(%q) = null, erwartet [%q, %q]", tc.name, tc.eingabe, tc.erwartetA, tc.erwartetB)
			continue
		}
		if got[0] != tc.erwartetA || got[1] != tc.erwartetB {
			t.Errorf("%s: parseKoordinatenPaar(%q) = %v, erwartet [%q, %q]", tc.name, tc.eingabe, *got, tc.erwartetA, tc.erwartetB)
		}
	}
}

// strconvQuote liefert eine für ein JavaScript-Import-Statement gültige
// Zeichenkette — Go-Pfade unter macOS/Linux enthalten kein Anführungszeichen,
// aber json.Marshal maskiert zuverlässig, falls doch (etwa Leerzeichen im
// Verzeichnisnamen des Systemtemp).
func strconvQuote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
