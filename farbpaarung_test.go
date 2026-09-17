package designsystem

import (
	"regexp"
	"strings"
	"testing"
)

// blockRe zerlegt CSS in flache "Selektor { Deklarationen }"-Blöcke.
// base.css enthält keine verschachtelten Regeln — bei @media-Blöcken liefert
// das Muster deshalb genau die innere Regel (der äußere @media-Rahmen bleibt
// unvollständig und damit unberücksichtigt, was hier gewollt ist: uns
// interessiert nur die tatsächliche Deklaration, nicht die Media-Bedingung).
var blockRe = regexp.MustCompile(`([^{}]+)\{([^{}]*)\}`)

// colorVarRe und bgVarRe finden "color: var(--x)" bzw.
// "background[-color]: var(--y)" INNERHALB EINES Blocks.
// colorVarRe verlangt Leerraum, ";" oder Blockanfang vor "color:" — ein
// vorangehender Bindestrich (wie in "border-color:") ist kein Wortzeichen
// und würde \b sonst fälschlich als Wortgrenze durchgehen lassen.
var (
	colorVarRe = regexp.MustCompile(`(?:^|[\s;])color:\s*var\((--[a-z-]+)\)`)
	bgVarRe    = regexp.MustCompile(`(?:^|[\s;])background(?:-color)?:\s*var\((--[a-z-]+)\)`)
)

// TestFarbpaarungenInBaseCSSErreichenAAA prüft jeden Regelblock in base.css,
// der SOWOHL color: var(--…) ALS AUCH background/background-color: var(--…)
// setzt, gegeneinander: beide Tokens sind einzeln gegen --bg und --card
// geprüft (tokens_test.go), das beweist aber nicht, dass eine tatsächliche
// CSS-Regel sie sinnvoll kombiniert. Nachgewiesen: eine neue Regel mit
// color: var(--text-muted); background: var(--accent); kombiniert zwei
// einzeln geprüfte Tokens zu 1.15:1, während alle bisherigen Tests grün
// bleiben.
//
// Diese Prüfung erfasst absichtlich nur Paare, die INNERHALB DESSELBEN
// Blocks gesetzt werden (wie es bei body, .btn, input/select/textarea,
// .btn-secondary usw. der Fall ist) — das ist mit einfacher Textzerlegung
// robust machbar, weil base.css keine verschachtelten Regeln enthält. Eine
// Farbe, die aus einem ANDEREN Block geerbt wird (etwa .btn:hover, das nur
// background überschreibt und color von .btn übernimmt), erfasst sie NICHT
// — das verlangte den vollen CSS-Kaskadenalgorithmus, den diese einfache
// Zerlegung bewusst nicht nachbildet. Diese Grenze ist in README.md und
// docs/design.md dokumentiert.
func TestFarbpaarungenInBaseCSSErreichenAAA(t *testing.T) {
	hell, dunkel := parseThemes(t)
	kommentarlos := removeComments(string(BaseCSS()))
	for _, block := range blockRe.FindAllStringSubmatch(kommentarlos, -1) {
		selektor := strings.TrimSpace(collapseWhitespace(block[1]))
		decls := block[2]
		cm := colorVarRe.FindStringSubmatch(decls)
		bm := bgVarRe.FindStringSubmatch(decls)
		if cm == nil || bm == nil {
			continue
		}
		vorn, hinten := cm[1], bm[1]
		for _, c := range []struct {
			thema  string
			tokens map[string]string
		}{{"hell", hell}, {"dunkel", dunkel}} {
			pruefeFarbpaar(t, selektor, c.thema, c.tokens, vorn, hinten)
		}
	}
}

func pruefeFarbpaar(t *testing.T, selektor, thema string, tokens map[string]string, vorn, hinten string) {
	t.Helper()
	v, ok := tokens[vorn]
	if !ok {
		t.Fatalf("[%s] %s: Token %s fehlt in tokens.css", thema, selektor, vorn)
	}
	h, ok := tokens[hinten]
	if !ok {
		t.Fatalf("[%s] %s: Token %s fehlt in tokens.css", thema, selektor, hinten)
	}
	a, err := ParseHex(v)
	if err != nil {
		t.Fatalf("[%s] %s: %s: %v", thema, selektor, vorn, err)
	}
	b, err := ParseHex(h)
	if err != nil {
		t.Fatalf("[%s] %s: %s: %v", thema, selektor, hinten, err)
	}
	if got := ContrastRatio(a, b); got < 7.0 {
		t.Errorf("[%s] %s: %s (%s) auf %s (%s): %.2f:1, gefordert >= 7:1",
			thema, selektor, vorn, v, hinten, h, got)
	}
}

func collapseWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// TestFarbpaarungErkennungFindetSchlechtesPaar belegt anhand eines eigenen
// CSS-Texts (nicht base.css selbst), dass die Zerlegung ein color/background
// Paar innerhalb eines Blocks tatsächlich extrahiert — genau der Fall, den
// TestFarbpaarungenInBaseCSSErreichenAAA gegen 7:1 prüft.
func TestFarbpaarungErkennungFindetSchlechtesPaar(t *testing.T) {
	css := `.beispiel { color: var(--text-muted); background: var(--accent); }`
	bloecke := blockRe.FindAllStringSubmatch(removeComments(css), -1)
	if len(bloecke) != 1 {
		t.Fatalf("erwartet 1 Block, gefunden %d", len(bloecke))
	}
	cm := colorVarRe.FindStringSubmatch(bloecke[0][2])
	bm := bgVarRe.FindStringSubmatch(bloecke[0][2])
	if cm == nil || bm == nil {
		t.Fatal("Farbpaar (color/background) wird aus dem Testblock nicht extrahiert")
	}
	if cm[1] != "--text-muted" || bm[1] != "--accent" {
		t.Errorf("erwartet --text-muted auf --accent, erhalten %s auf %s", cm[1], bm[1])
	}
}
