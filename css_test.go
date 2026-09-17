package designsystem

import (
	"regexp"
	"strings"
	"testing"
)

var (
	farbliteralRe = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b|\brgba?\(|\bhsla?\(`)
	varNutzungRe  = regexp.MustCompile(`var\((--[a-z-]+)`)
	varDefRe      = regexp.MustCompile(`(--[a-z-]+):`)
)

func TestBaseCSSEnthaeltKeineFarbliterale(t *testing.T) {
	// Eine Farbe im Klartext in base.css umgeht die Kontrastprüfung aus
	// tokens_test.go vollständig — sie stünde nirgends, wo der Test sie
	// fände. Ausgenommen sind Schwarz- und Weißwerte in Schatten, die
	// keine Textfarbe sind; die stehen in tokens.css.
	for i, zeile := range strings.Split(string(BaseCSS()), "\n") {
		if strings.HasPrefix(strings.TrimSpace(zeile), "/*") || strings.HasPrefix(strings.TrimSpace(zeile), "*") {
			continue
		}
		if m := farbliteralRe.FindString(zeile); m != "" {
			t.Errorf("base.css:%d enthält den Farbwert %q — gehört nach tokens.css:\n  %s",
				i+1, m, strings.TrimSpace(zeile))
		}
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
