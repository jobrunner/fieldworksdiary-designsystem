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

// --- .koord-gitter: zwei Spalten als Regelfall --------------------------
//
// TestKoordGitterHatZweiSpaltenAlsRegelfall belegt die Nachbesserung: eine
// frühere Fassung hatte grid-template-columns: 1fr als Grundregel und
// schaltete erst ab 640px auf zwei Spalten um — die Felder standen damit
// im häufigen, schmalen Fall (Expertus im Gelände auf dem Telefon)
// untereinander, obwohl "nebeneinander" ausdrücklich verlangt war. Dieser
// Test prüft NUR, dass die CSS-Regel die richtige Form hat
// (repeat(auto-fit, minmax(…)) als Grundregel, keine Fensterbreiten-Abfrage
// mehr für .koord-gitter) — er belegt die ANWESENHEIT der Regel, nicht die
// tatsächliche Darstellung im Browser; das bleibt ohne einen echten
// Browser-Testlauf offen.
func TestKoordGitterHatZweiSpaltenAlsRegelfall(t *testing.T) {
	css := removeComments(string(BaseCSS()))

	if !strings.Contains(css, "grid-template-columns: repeat(auto-fit, minmax(") {
		t.Error(`base.css enthält "grid-template-columns: repeat(auto-fit, minmax(...))" für .koord-gitter nicht — zwei Spalten sollen der Regelfall sein, unabhängig von der Fensterbreite`)
	}

	// Die Grundregel darf nicht mehr "1fr" ohne repeat(auto-fit, ...) sein
	// (der alte, fehlerhafte Zustand: eine Spalte als Vorgabe).
	i := strings.Index(css, ".koord-gitter")
	if i < 0 {
		t.Fatal(`base.css enthält ".koord-gitter" nicht`)
	}
	open := strings.Index(css[i:], "{")
	close := strings.Index(css[i:], "}")
	if open < 0 || close < 0 || close < open {
		t.Fatal(`.koord-gitter-Regel konnte nicht abgegrenzt werden`)
	}
	rumpf := css[i+open : i+close]
	if strings.Contains(rumpf, "grid-template-columns: 1fr;") {
		t.Error(`.koord-gitter hat "grid-template-columns: 1fr;" als Grundregel — das ist der behobene Fehler (eine Spalte als Regelfall statt als Ausnahme)`)
	}
}

// TestKoordGitterKeineMediaAbfrageMehr belegt, dass .koord-gitter nicht
// mehr über eine @media(min-width)-Abfrage auf zwei Spalten umschaltet:
// eine Abfrage nach der FENSTERBREITE hilft nicht, wenn die Bedienform in
// einer schmalen Spalte innerhalb eines breiten Fensters steht — genau der
// Einwand aus der Nachbesserung. repeat(auto-fit, ...) bemisst sich am
// verfügbaren Platz des Elements selbst und braucht deshalb keine eigene
// @media-Regel mehr für .koord-gitter.
func TestKoordGitterKeineMediaAbfrageMehr(t *testing.T) {
	css := removeComments(string(BaseCSS()))
	i := strings.Index(css, "@media (min-width: 640px)")
	if i < 0 {
		t.Fatal(`base.css enthält "@media (min-width: 640px)" nicht mehr — Testannahme verletzt`)
	}
	open := strings.Index(css[i:], "{")
	if open < 0 {
		t.Fatal("konnte den @media-Block nicht abgrenzen")
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
	block := css[i:ende]
	if strings.Contains(block, "koord-gitter") {
		t.Error(`der @media (min-width: 640px)-Block erwähnt ".koord-gitter" noch — die Spaltenzahl soll sich am verfügbaren Platz orientieren, nicht an der Fensterbreite`)
	}
}
