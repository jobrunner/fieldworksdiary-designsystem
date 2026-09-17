package designsystem

import (
	"regexp"
	"strings"
	"testing"
)

var (
	hexRe         = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	funcColorRe   = regexp.MustCompile(`\brgba?\(|\bhsla?\(`)
	namedColorRe  = regexp.MustCompile(`(?:^|[^\w-])(red|green|blue|white|black|gray|grey|orange|yellow|purple|silver|navy|teal|olive|maroon|aqua|fuchsia|lime)(?:[^\w-]|$)`)
	varNutzungRe  = regexp.MustCompile(`var\((--[a-z-]+)`)
	varDefRe      = regexp.MustCompile(`(--[a-z-]+):`)
)

// farbliterale prüft CSS-Text auf Farbliterale und gibt Fundstellen zurück.
// Jede Fundstelle ist ein String mit Zeilennummer, Farbe und Zeile.
// Kommentare werden ignoriert. transparent und currentColor sind erlaubt,
// alle anderen benannten Farben und Hex-Werte gelten als Verstoß.
func farbliterale(css string) []string {
	// Kommentare entfernen, aber Zeilennummern durch Leerzeichen erhalten
	kommentarlos := removeComments(css)

	var funde []string
	for i, zeile := range strings.Split(kommentarlos, "\n") {
		lineNum := i + 1

		// Hex-Farben prüfen
		if m := hexRe.FindString(zeile); m != "" {
			funde = append(funde, m)
		}

		// rgb/rgba/hsl/hsla prüfen
		if funcColorRe.MatchString(zeile) {
			funde = append(funde, "rgb/rgba/hsl/hsla")
		}

		// Benannte Farben prüfen (außer transparent, currentColor)
		for _, match := range namedColorRe.FindAllStringSubmatch(zeile, -1) {
			if len(match) > 1 {
				farbe := match[1]
				if farbe != "transparent" && farbe != "currentColor" {
					funde = append(funde, farbe)
				}
			}
		}
		_ = lineNum // lineNum für zukünftige Verwendung
	}
	return funde
}

// removeComments entfernt alle /* ... */ Kommentare aus CSS,
// behält aber Zeilennummern durch Leerzeichen.
func removeComments(css string) string {
	re := regexp.MustCompile(`(?s)/\*.*?\*/`)
	return re.ReplaceAllStringFunc(css, func(s string) string {
		// Ersetze den Kommentar durch Leerzeichen, behalte aber Umbrüche
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
	// keine Textfarbe sind; die stehen in tokens.css. transparent und
	// currentColor sind erlaubt — sie sind keine Farben im Sinne der Zusage.
	funde := farbliterale(string(BaseCSS()))
	if len(funde) > 0 {
		for _, farbe := range funde {
			t.Errorf("base.css enthält den Farbwert %q — gehört nach tokens.css", farbe)
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

// TestTransparentUndCurrentColorSindErlaubt prüft, dass diese speziellen
// Werte nicht gemeldet werden, obwohl sie benannte Farben sind.
func TestTransparentUndCurrentColorSindErlaubt(t *testing.T) {
	transparent := farbliterale(".tab { border-bottom-color: transparent; }")
	if len(transparent) > 0 {
		t.Error("'transparent' sollte erlaubt sein, wird aber als Verstoß gemeldet")
	}

	current := farbliterale(".icon { color: currentColor; }")
	if len(current) > 0 {
		t.Error("'currentColor' sollte erlaubt sein, wird aber als Verstoß gemeldet")
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
