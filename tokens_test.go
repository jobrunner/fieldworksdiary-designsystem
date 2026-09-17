package designsystem

import (
	"regexp"
	"strings"
	"testing"
)

// tokenRe findet "--name: wert;" — für Farbwerte genügt das, weil tokens.css
// bewusst keine Funktionsaufrufe oder mehrteiligen Farbangaben enthält.
var tokenRe = regexp.MustCompile(`(--[a-z-]+):\s*([^;]+);`)

// parseThemes zerlegt tokens.css in das helle und das dunkle Thema. Der
// dunkle Satz steht im @media-Block; alles davor gehört zum hellen.
func parseThemes(t *testing.T) (hell, dunkel map[string]string) {
	t.Helper()
	css := string(TokensCSS())
	i := strings.Index(css, "@media")
	if i < 0 {
		t.Fatal("tokens.css enthält keinen @media-Block für das dunkle Thema")
	}
	hell = map[string]string{}
	for _, m := range tokenRe.FindAllStringSubmatch(css[:i], -1) {
		hell[m[1]] = strings.TrimSpace(m[2])
	}
	// Das dunkle Thema erbt alles, was es nicht selbst überschreibt.
	dunkel = map[string]string{}
	for k, v := range hell {
		dunkel[k] = v
	}
	for _, m := range tokenRe.FindAllStringSubmatch(css[i:], -1) {
		dunkel[m[1]] = strings.TrimSpace(m[2])
	}
	return hell, dunkel
}

func TestTextTokensErreichenAAA(t *testing.T) {
	// Geprüft wird gegen BEIDE Flächen: --card ist im dunklen Thema heller
	// als --bg und damit die strengere Bezugsfläche. Ein Wert, der nur
	// gegen den Seitenhintergrund geprüft ist, versagt auf der Karte.
	textTokens := []string{"--text", "--text-muted", "--text-disabled"}

	hell, dunkel := parseThemes(t)
	for _, c := range []struct {
		thema  string
		tokens map[string]string
	}{{"hell", hell}, {"dunkel", dunkel}} {
		for _, name := range textTokens {
			for _, flaeche := range []string{"--bg", "--card"} {
				pruefe(t, c.thema, c.tokens, name, flaeche, 7.0)
			}
		}
	}
}

func TestStatusfarbenErreichenAAA(t *testing.T) {
	hell, dunkel := parseThemes(t)
	for _, c := range []struct {
		thema  string
		tokens map[string]string
	}{{"hell", hell}, {"dunkel", dunkel}} {
		for _, name := range []string{"--success", "--error", "--warning"} {
			for _, flaeche := range []string{"--bg", "--card"} {
				pruefe(t, c.thema, c.tokens, name, flaeche, 7.0)
			}
		}
	}
}

func TestControlLineErreicht3zu1(t *testing.T) {
	// WCAG 1.4.11 fordert 3:1 für die Begrenzung, die ein BEDIENELEMENT
	// erkennbar macht. --border ist davon ausgenommen und wird hier
	// bewusst nicht geprüft.
	hell, dunkel := parseThemes(t)
	for _, c := range []struct {
		thema  string
		tokens map[string]string
	}{{"hell", hell}, {"dunkel", dunkel}} {
		for _, flaeche := range []string{"--bg", "--card"} {
			pruefe(t, c.thema, c.tokens, "--control-line", flaeche, 3.0)
		}
	}
}

func TestWeisserTextAufAkzent(t *testing.T) {
	// Der primäre Knopf trägt weiße Schrift auf --accent. Im dunklen Thema
	// ist --accent eine Textfarbe, keine Fläche — dort gilt die Prüfung aus
	// TestTextTokensErreichenAAA sinngemäß über --accent selbst.
	hell, _ := parseThemes(t)
	akzent, err := ParseHex(hell["--accent"])
	if err != nil {
		t.Fatalf("--accent: %v", err)
	}
	weiss, _ := ParseHex("#ffffff")
	if got := ContrastRatio(akzent, weiss); got < 7.0 {
		t.Errorf("weißer Text auf --accent (%s): %.2f:1, gefordert >= 7:1", hell["--accent"], got)
	}
}

func TestAkzentImDunklenThemaIstLesbar(t *testing.T) {
	_, dunkel := parseThemes(t)
	for _, flaeche := range []string{"--bg", "--card"} {
		pruefe(t, "dunkel", dunkel, "--accent", flaeche, 7.0)
	}
}

func pruefe(t *testing.T, thema string, tokens map[string]string, vorn, hinten string, min float64) {
	t.Helper()
	v, ok := tokens[vorn]
	if !ok {
		t.Fatalf("[%s] Token %s fehlt in tokens.css", thema, vorn)
	}
	h, ok := tokens[hinten]
	if !ok {
		t.Fatalf("[%s] Token %s fehlt in tokens.css", thema, hinten)
	}
	a, err := ParseHex(v)
	if err != nil {
		t.Fatalf("[%s] %s: %v", thema, vorn, err)
	}
	b, err := ParseHex(h)
	if err != nil {
		t.Fatalf("[%s] %s: %v", thema, hinten, err)
	}
	if got := ContrastRatio(a, b); got < min {
		t.Errorf("[%s] %s (%s) auf %s (%s): %.2f:1, gefordert >= %.1f:1",
			thema, vorn, v, hinten, h, got, min)
	}
}
