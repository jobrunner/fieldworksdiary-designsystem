package designsystem

import (
	"regexp"
	"strings"
	"testing"
)

// fund beschreibt eine gefundene Farbe im Klartext.
type fund struct {
	zeile int
	wert  string
}

var (
	hexRe         = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	funcColorRe   = regexp.MustCompile(`\brgba?\(|\bhsla?\(`)
	namedColorRe  = regexp.MustCompile(`(?:^|[^\w-])(red|green|blue|white|black|gray|grey|orange|yellow|purple|silver|navy|teal|olive|maroon|aqua|fuchsia|lime)(?:[^\w-]|$)`)
	varNutzungRe  = regexp.MustCompile(`var\((--[a-z-]+)`)
	varDefRe      = regexp.MustCompile(`(--[a-z-]+):`)
)

// farbliterale prüft CSS-Text auf Farbliterale und gibt Fundstellen mit
// Zeilennummern zurück. Kommentare werden ignoriert. transparent ist erlaubt,
// alle anderen benannten Farben und Hex-Werte gelten als Verstoß.
func farbliterale(css string) []fund {
	// Kommentare entfernen, aber Zeilennummern durch Zeilenumbrüche erhalten
	kommentarlos := removeComments(css)

	var funde []fund
	for i, zeile := range strings.Split(kommentarlos, "\n") {
		lineNum := i + 1

		// Hex-Farben prüfen
		if m := hexRe.FindString(zeile); m != "" {
			funde = append(funde, fund{zeile: lineNum, wert: m})
		}

		// rgb/rgba/hsl/hsla prüfen
		if funcColorRe.MatchString(zeile) {
			funde = append(funde, fund{zeile: lineNum, wert: "rgb/rgba/hsl/hsla"})
		}

		// Benannte Farben prüfen (außer transparent)
		for _, match := range namedColorRe.FindAllStringSubmatch(zeile, -1) {
			if len(match) > 1 {
				farbe := match[1]
				if farbe != "transparent" {
					funde = append(funde, fund{zeile: lineNum, wert: farbe})
				}
			}
		}
	}
	return funde
}

// removeComments entfernt alle /* ... */ Kommentare aus CSS,
// behält aber alle Zeilenumbrüche: damit bleiben die Zeilennummern
// in den Originalzellen erhalten.
func removeComments(css string) string {
	re := regexp.MustCompile(`(?s)/\*.*?\*/`)
	return re.ReplaceAllStringFunc(css, func(s string) string {
		// Ersetze den Kommentar charakterweise: Zeilenumbrüche bleiben,
		// alles andere wird zu Leerzeichen.
		return strings.Map(func(r rune) rune {
			if r == '\n' {
				return '\n'
			}
			return ' '
		}, s)
	})
}

func TestBaseCSSEnthaeltKeineFarbliterale(t *testing.T) {
	// Eine Farbe im Klartext in base.css umgeht die Kontrastprüfung aus
	// tokens_test.go vollständig — sie stünde nirgends, wo der Test sie
	// fände. Ausgenommen sind Schwarz- und Weißwerte in Schatten, die
	// keine Textfarbe sind; die stehen in tokens.css. transparent ist erlaubt
	// — es ist keine Farbe im Sinne der Zusage.
	funde := farbliterale(string(BaseCSS()))
	if len(funde) > 0 {
		for _, f := range funde {
			t.Errorf("base.css:%d enthält den Farbwert %q — gehört nach tokens.css",
				f.zeile, f.wert)
		}
	}
}

// TestUniversalselectorUmweichungVerhindert prüft, dass der Test die
// Umweichung über * { ... } erkennt.
func TestUniversalselectorUmweichungVerhindert(t *testing.T) {
	funde := farbliterale("* { color: #ff0000; }")
	if len(funde) == 0 {
		t.Error("Test lässt Universalselektor mit Farbliteral durch — sollte erkannt werden")
	}
}

// TestBenannteFarbenWerdenErkannt prüft, dass benannte CSS-Farben erkannt
// werden.
func TestBenannteFarbenWerdenErkannt(t *testing.T) {
	funde := farbliterale("body { color: red; }")
	if len(funde) == 0 {
		t.Error("Test erkennt benannte Farbe 'red' nicht — sollte erkannt werden")
	}
}

// TestTransparentIstErlaubt prüft, dass transparent nicht als Verstoß
// gemeldet wird — es ist keine Farbe im Sinne der Zusage.
func TestTransparentIstErlaubt(t *testing.T) {
	funde := farbliterale(".tab { border-bottom-color: transparent; }")
	if len(funde) > 0 {
		t.Error("'transparent' sollte erlaubt sein, wird aber als Verstoß gemeldet")
	}
}

// TestKommentareWerdenIgnoriert prüft, dass Farben in Kommentaren nicht
// gemeldet werden.
func TestKommentareWerdenIgnoriert(t *testing.T) {
	funde := farbliterale("/* rot: #ff0000 */\nbody { color: var(--text); }")
	if len(funde) > 0 {
		t.Error("Farbe in Kommentar sollte ignoriert werden, wird aber gemeldet")
	}
}

// TestMehzeiligKommentareErhaltenZeilennummern prüft, dass die Zeilennummern
// auch nach mehzeiligen Kommentaren noch korrekt sind.
func TestMehzeiligKommentareErhaltenZeilennummern(t *testing.T) {
	css := `/* Kommentar
über mehrere
Zeilen */
body { color: red; }`
	funde := farbliterale(css)
	if len(funde) != 1 {
		t.Fatalf("erwartet 1 Fund, got %d", len(funde))
	}
	// Der Verstoß sollte auf Zeile 4 sein (nach dem 3zeiligen Kommentar)
	if funde[0].zeile != 4 {
		t.Errorf("erwartet Zeile 4, got %d", funde[0].zeile)
	}
	if funde[0].wert != "red" {
		t.Errorf("erwartet 'red', got %q", funde[0].wert)
	}
}

func TestBaseCSSNutztNurDefinierteTokens(t *testing.T) {
	definiert := map[string]bool{}
	for _, m := range varDefRe.FindAllStringSubmatch(string(TokensCSS()), -1) {
		definiert[m[1]] = true
	}
	gesehen := map[string]bool{}
	for _, m := range varNutzungRe.FindAllStringSubmatch(string(BaseCSS()), -1) {
		if !definiert[m[1]] && !gesehen[m[1]] {
			gesehen[m[1]] = true
			t.Errorf("base.css nutzt %s, das in tokens.css nicht definiert ist", m[1])
		}
	}
}

func TestCSSLiefertBeideDateienInReihenfolge(t *testing.T) {
	// Die Reihenfolge ist bindend: base.css greift auf Variablen zu, die
	// tokens.css setzt. Umgekehrt eingebunden bliebe jede Farbe leer.
	zusammen := string(CSS())
	iTokens := strings.Index(zusammen, "--control-line:")
	iBase := strings.Index(zusammen, ".card {")
	if iTokens < 0 || iBase < 0 {
		t.Fatalf("CSS() enthält nicht beide Dateien (tokens: %d, base: %d)", iTokens, iBase)
	}
	if iTokens > iBase {
		t.Error("CSS() liefert base.css vor tokens.css — die Variablen wären beim Auswerten noch nicht gesetzt")
	}
}
