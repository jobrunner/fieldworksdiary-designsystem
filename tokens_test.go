package designsystem

import (
	"regexp"
	"strings"
	"testing"
)

// tokenRe findet "--name: wert;" — für Farbwerte genügt das, weil tokens.css
// bewusst keine Funktionsaufrufe oder mehrteiligen Farbangaben enthält.
var tokenRe = regexp.MustCompile(`(--[a-z-]+):\s*([^;]+);`)

// keineTextfarbe listet Token, die keine Textfarbe sind und deshalb nicht
// der 7:1-Pflicht aus TestAlleTextfarbenErreichenAAA unterliegen. Jedes
// andere Token in tokens.css — auch ein künftig hinzugefügtes — muss diese
// Pflicht erfüllen. Ein neues Token kostet damit eine bewusste Aufnahme in
// diese Liste statt stillschweigend ungeprüft zu bleiben.
//
// Diese Liste ist die Gefahrenstelle des Moduls: jedes hier eingetragene
// Token entgeht der allgemeinen Prüfung. Deshalb steht hinter jedem Eintrag
// ein Verweis auf die Prüfung, die es STATTDESSEN abdeckt — ein Eintrag ohne
// eigene gezielte Prüfung wäre schlicht ungeprüft.
var keineTextfarbe = map[string]bool{
	"--bg":           true,
	"--card":         true,
	"--border":       true,
	"--control-line": true, // eigene 3:1-Prüfung, siehe TestControlLineErreicht3zu1
	"--shadow":       true,
	"--radius":       true,
	"--radius-sm":    true,
	"--font":         true,
	"--font-mono":    true,
	// --accent und --accent-hover sind Flächen, keine Textfarben: der
	// primäre Knopf trägt --accent-on darauf, nicht --card. --accent als
	// Text gegen dunklen Grund gelesen erreicht nur 1.91:1 und würde diesen
	// Test zu Recht brechen. Eigene Prüfung: TestKnopfbeschriftungAufAkzent
	// (--accent-on auf --accent UND --accent-hover, beide Themen).
	"--accent":       true,
	"--accent-hover": true,
	// --accent-on ist Text auf der Akzentfläche, nicht auf --bg/--card —
	// dieselbe TestKnopfbeschriftungAufAkzent deckt es ab.
	"--accent-on": true,
}

// textTokenNamen liefert die Namen aller Token in einem Themen-Satz, die
// laut keineTextfarbe nicht ausgenommen sind — also als Textfarbe gelten
// und 7:1 gegen --bg und --card erfüllen müssen.
func textTokenNamen(tokens map[string]string) []string {
	var namen []string
	for name := range tokens {
		if !keineTextfarbe[name] {
			namen = append(namen, name)
		}
	}
	return namen
}

// parseThemesAusText ist der eigentliche Zerlegungsschritt, unabhängig von
// tokens.css — parseThemes (für die echten Tests) und
// TestParseThemesIgnoriertKommentare (für die Kommentarfälle) rufen ihn auf
// unterschiedlichem CSS-Text auf.
//
// Kommentare werden VOR jeder weiteren Auswertung entfernt (removeComments,
// siehe css_test.go) — sowohl bevor die @media-Grenze per strings.Index
// bestimmt wird, als auch bevor tokenRe darauf läuft. Ohne das könnte ein
// Kommentar mit einem eigenen "--token: wert;"-Muster den echten Wert
// überschreiben (der zuletzt gefundene Treffer gewinnt), oder ein "@media"
// im Kommentartext die Aufteilung zwischen hellem und dunklem Thema
// verschieben. TestParseThemesIgnoriertKommentare belegt beide Fälle.
func parseThemesAusText(css string) (hell, dunkel map[string]string, mediaGefunden bool) {
	css = removeComments(css)
	i := strings.Index(css, "@media")
	if i < 0 {
		return nil, nil, false
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
	return hell, dunkel, true
}

// parseThemes zerlegt tokens.css in das helle und das dunkle Thema. Der
// dunkle Satz steht im @media-Block; alles davor gehört zum hellen.
// Die Annahme, dass GENAU EIN @media-Block existiert und dieser das dunkle
// Thema ist, wird in TestGenauEinMediaBlockFuerDunkel gesondert geprüft.
func parseThemes(t *testing.T) (hell, dunkel map[string]string) {
	t.Helper()
	hell, dunkel, ok := parseThemesAusText(string(TokensCSS()))
	if !ok {
		t.Fatal("tokens.css enthält keinen @media-Block für das dunkle Thema")
	}
	return hell, dunkel
}

// TestParseThemesIgnoriertKommentare belegt den Befund: ein Kommentar mit
// einem eigenen "--token: wert;"-Muster darf den echten Wert weder
// überschreiben noch vortäuschen, dass er es täte. Geprüft wird die
// Zerlegung direkt an einem eigenen CSS-Text, nicht an tokens.css selbst —
// tokens.css bleibt dabei unangetastet und der Test bleibt unabhängig von
// seinem jeweiligen Inhalt.
func TestParseThemesIgnoriertKommentare(t *testing.T) {
	t.Run("Kommentar nach echtem Wert wird nicht übernommen", func(t *testing.T) {
		// Genau der nachgewiesene Fall: ein Kommentar STEHT NACH der echten
		// Definition, im selben :root-Block, und nennt selbst einen Wert für
		// dasselbe Token.
		css := `:root {
  --error: #991b1b;
  /* Fruehrerer Entwurf war --error: #00ff00; zu grell */
}
@media (prefers-color-scheme: dark) {
  :root {
    --error: #fca5a5;
  }
}`
		hell, dunkel, ok := parseThemesAusText(css)
		if !ok {
			t.Fatal("kein @media-Block gefunden")
		}
		if hell["--error"] != "#991b1b" {
			t.Errorf("--error (hell) = %q, erwartet den echten Wert #991b1b — der Kommentarwert #00ff00 hätte ihn überschrieben", hell["--error"])
		}
		if dunkel["--error"] != "#fca5a5" {
			t.Errorf("--error (dunkel) = %q, erwartet #fca5a5", dunkel["--error"])
		}
	})

	t.Run("gefährliche Richtung: schlechter echter Wert bleibt trotz gutem Kommentarwert bestehen", func(t *testing.T) {
		// Die von der Prüfung eigentlich gefürchtete Richtung: der ECHTE
		// Wert ist schlecht, ein NACHFOLGENDER Kommentar nennt einen guten
		// Wert desselben Tokens. Bliebe der Kommentarwert maßgeblich, würde
		// die Kontrastprüfung einen Text prüfen, den niemand ausliefert, und
		// eine tatsächlich zusagenverletzende Palette für grün erklären.
		css := `:root {
  --error: #00ff00;
  /* besser waere --error: #991b1b; */
}
@media (prefers-color-scheme: dark) {
  :root {}
}`
		hell, _, ok := parseThemesAusText(css)
		if !ok {
			t.Fatal("kein @media-Block gefunden")
		}
		if hell["--error"] != "#00ff00" {
			t.Errorf("--error (hell) = %q, erwartet den echten (schlechten) Wert #00ff00 — sonst prüfte der Test einen Wert, der nicht ausgeliefert wird", hell["--error"])
		}
	})

	t.Run("@media in einem Kommentar verschiebt die Themengrenze nicht", func(t *testing.T) {
		// Ein "@media" im Kommentartext, VOR dem echten dunklen Block, darf
		// strings.Index nicht vorzeitig zuschlagen lassen — sonst rechnete
		// parseThemesAusText einen Teil des hellen Themas dem dunklen zu.
		css := `:root {
  --text: #1e293b;
  /* frueher ohne @media (prefers-color-scheme: dark) Unterstuetzung */
  --error: #991b1b;
}
@media (prefers-color-scheme: dark) {
  :root {
    --text: #f1f5f9;
  }
}`
		hell, dunkel, ok := parseThemesAusText(css)
		if !ok {
			t.Fatal("kein @media-Block gefunden")
		}
		if hell["--text"] != "#1e293b" || hell["--error"] != "#991b1b" {
			t.Errorf("helles Thema unvollständig: %#v", hell)
		}
		if dunkel["--text"] != "#f1f5f9" {
			t.Errorf("--text (dunkel) = %q, erwartet #f1f5f9 — das @media im Kommentar hat die Grenze verschoben", dunkel["--text"])
		}
	})
}

// TestGenauEinMediaBlockFuerDunkel hält die Annahme aus parseThemes fest,
// statt sie stillschweigend vorauszusetzen: stünde je ein @media
// (min-width: …) vor dem Dunkelblock, würde parseThemes die halbe helle
// Palette dem dunklen Thema zurechnen und gegen die falschen Flächen
// prüfen.
func TestGenauEinMediaBlockFuerDunkel(t *testing.T) {
	// removeComments aus demselben Grund wie in parseThemes: ein "@media" im
	// Kommentartext darf diese Zählung nicht verfälschen.
	css := removeComments(string(TokensCSS()))
	anzahl := strings.Count(css, "@media")
	if anzahl != 1 {
		t.Fatalf("tokens.css enthält %d @media-Blöcke, parseThemes setzt genau einen voraus (das dunkle Thema)", anzahl)
	}
	rest := css[strings.Index(css, "@media"):]
	if j := strings.Index(rest, "{"); j >= 0 {
		rest = rest[:j]
	}
	if !strings.Contains(rest, "prefers-color-scheme: dark") {
		t.Errorf("der einzige @media-Block lautet %q, erwartet prefers-color-scheme: dark", strings.TrimSpace(rest))
	}
}

// TestAlleTextfarbenErreichenAAA prüft jedes Token, das keine Ausnahme in
// keineTextfarbe ist, gegen BEIDE Flächen: --card ist im dunklen Thema
// heller als --bg und damit die strengere Bezugsfläche. Ein Wert, der nur
// gegen den Seitenhintergrund geprüft ist, versagt auf der Karte.
//
// Die Liste der geprüften Namen steht damit nicht im Testcode, sondern
// ergibt sich aus tokens.css selbst — ein neues Textfarb-Token wird ohne
// weiteres Zutun erfasst.
func TestAlleTextfarbenErreichenAAA(t *testing.T) {
	hell, dunkel := parseThemes(t)
	for _, c := range []struct {
		thema  string
		tokens map[string]string
	}{{"hell", hell}, {"dunkel", dunkel}} {
		for _, name := range textTokenNamen(c.tokens) {
			for _, flaeche := range []string{"--bg", "--card"} {
				pruefe(t, c.thema, c.tokens, name, flaeche, 7.0)
			}
		}
	}
}

// TestAuswahllogikErfasstNeuesToken belegt anhand eines eigenen CSS-Texts
// (nicht tokens.css selbst), dass textTokenNamen ein hinzugefügtes
// Farbtoken tatsächlich in die Prüfung aufnimmt und eine Ausnahme wie --bg
// nicht fälschlich mit aufnimmt.
func TestAuswahllogikErfasstNeuesToken(t *testing.T) {
	css := `:root {
		--bg: #ffffff;
		--text-ganz-neu: #112233;
	}`
	tokens := map[string]string{}
	for _, m := range tokenRe.FindAllStringSubmatch(css, -1) {
		tokens[m[1]] = strings.TrimSpace(m[2])
	}
	namen := textTokenNamen(tokens)
	gefunden := false
	for _, n := range namen {
		if n == "--text-ganz-neu" {
			gefunden = true
		}
		if n == "--bg" {
			t.Error("--bg ist keine Textfarbe und sollte von der Auswahllogik nicht erfasst werden")
		}
	}
	if !gefunden {
		t.Error("ein neu hinzugefügtes Token --text-ganz-neu wird von der Auswahllogik nicht erfasst — genau das soll dieser Test verhindern")
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

// TestKnopfbeschriftungAufAkzent prüft die Knopfbeschriftung gegen ihre
// tatsächliche Fläche in beiden Zuständen und beiden Themen. base.css setzt
// die Beschriftung auf var(--accent-on) (nicht auf --card) — dieser Test
// behauptet deshalb dasselbe wie das CSS, statt einen eigenen, davon
// abweichenden Literalwert zu prüfen. --card taugt hier nicht mehr: --accent
// ist als Fläche in beiden Themen derselbe dunkle Markenton, --card wäre im
// dunklen Thema dunkler Text auf dunkler Fläche.
//
// Dieser Test ist die gezielte Prüfung, die --accent, --accent-hover und
// --accent-on in tokens_test.go von TestAlleTextfarbenErreichenAAA ausnimmt
// (siehe keineTextfarbe) — ohne ihn wären alle drei schlicht ungeprüft.
func TestKnopfbeschriftungAufAkzent(t *testing.T) {
	hell, dunkel := parseThemes(t)
	for _, c := range []struct {
		thema  string
		tokens map[string]string
	}{{"hell", hell}, {"dunkel", dunkel}} {
		for _, flaeche := range []string{"--accent", "--accent-hover"} {
			pruefe(t, c.thema, c.tokens, "--accent-on", flaeche, 7.0)
		}
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
