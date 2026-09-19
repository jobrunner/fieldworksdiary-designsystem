package designsystem

import (
	"os"
	"os/exec"
	"testing"
)

// TestJSHatGueltigeSyntax prüft js/combobox.js mit "node --check" — kein
// Go-Test parst JavaScript, ein Tippfehler ("function kaputt( {") ginge
// sonst unbemerkt in fünf Dienste. Der eigentliche, verbindliche Schritt
// steht in .github/workflows/ci.yml; dieser Test macht dieselbe Prüfung
// zusätzlich lokal verfügbar und überspringt sich selbst, wenn "node"
// fehlt (t.Skip), statt den Build von einer lokal fehlenden Abhängigkeit
// abhängig zu machen.
//
// Geprüft wird eine .mjs-Kopie, nicht die Datei direkt: "node --check"
// erkennt js/combobox.js zwar am führenden "export" als ES-Modul, nutzt
// dafür bei einer .js-Datei aber einen laxeren Erkennungspfad, der ein
// abgeschnittenes Dateiende nachweislich NICHT als Syntaxfehler meldet
// (siehe .github/workflows/ci.yml für die Gegenprobe). Als .mjs erzwingt
// derselbe Aufruf den vollständigen ESM-Parser.
func TestJSHatGueltigeSyntax(t *testing.T) {
	nodePath, err := exec.LookPath("node")
	if err != nil {
		t.Skip(`"node" ist nicht installiert — diese Prüfung läuft verbindlich in .github/workflows/ci.yml`)
	}

	tmp, err := os.CreateTemp(t.TempDir(), "combobox-*.mjs")
	if err != nil {
		t.Fatalf("konnte Temp-Datei nicht anlegen: %v", err)
	}
	if _, err := tmp.Write(JS()); err != nil {
		t.Fatalf("konnte JS() nicht in Temp-Datei schreiben: %v", err)
	}
	if err := tmp.Close(); err != nil {
		t.Fatalf("konnte Temp-Datei nicht schließen: %v", err)
	}

	out, err := exec.Command(nodePath, "--check", tmp.Name()).CombinedOutput()
	if err != nil {
		t.Errorf("node --check meldet einen Syntaxfehler in js/combobox.js:\n%s", out)
	}
}

// TestJSHatGueltigeSyntaxErkenntKaputteDatei belegt, dass
// TestJSHatGueltigeSyntax tatsächlich anschlägt — derselbe reale Fall aus
// dem Befund: "function kaputt( {" ans Modul angehängt.
func TestJSHatGueltigeSyntaxErkenntKaputteDatei(t *testing.T) {
	nodePath, err := exec.LookPath("node")
	if err != nil {
		t.Skip(`"node" ist nicht installiert`)
	}

	kaputt := append(JS(), []byte("\nfunction kaputt( {\n")...)

	tmp, err := os.CreateTemp(t.TempDir(), "combobox-kaputt-*.mjs")
	if err != nil {
		t.Fatalf("konnte Temp-Datei nicht anlegen: %v", err)
	}
	if _, err := tmp.Write(kaputt); err != nil {
		t.Fatalf("konnte kaputte Datei nicht schreiben: %v", err)
	}
	if err := tmp.Close(); err != nil {
		t.Fatalf("konnte Temp-Datei nicht schließen: %v", err)
	}

	out, err := exec.Command(nodePath, "--check", tmp.Name()).CombinedOutput()
	if err == nil {
		t.Errorf("node --check meldet die angehängte kaputte Funktion nicht als Syntaxfehler:\n%s", out)
	}
}
