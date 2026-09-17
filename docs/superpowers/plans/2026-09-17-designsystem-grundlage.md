# Design-System-Grundlage — Implementierungsplan

> **Für agentische Bearbeiter:** ERFORDERLICHE UNTER-SKILL: `superpowers:subagent-driven-development` (empfohlen) oder `superpowers:executing-plans`, um diesen Plan Aufgabe für Aufgabe umzusetzen. Die Schritte nutzen Kontrollkästchen (`- [ ]`) zur Nachverfolgung.

**Ziel:** Das Go-Modul `github.com/jobrunner/fieldworksdiary-designsystem` mit geprüften Farb-Tokens, Basiskomponenten und einer Demo-Seite bauen — lauffähig und abnehmbar, bevor ein einziger Dienst es einbindet.

**Architektur:** Zwei CSS-Dateien, per `go:embed` in ein Go-Modul mit drei Funktionen gepackt. Ein Go-Test liest die Farbwerte aus `tokens.css` und rechnet die WCAG-Kontraste nach — die Werte stehen genau einmal im Repo, nicht zusätzlich im Testcode. Eine statische Demo-Seite zeigt jede Komponente in hellem und dunklem Thema.

**Tech-Stack:** Go 1.26, `embed`, reines CSS ohne Präprozessor. Keine weiteren Abhängigkeiten.

**Spec:** `docs/design.md`

## Globale Randbedingungen

- Modulpfad: `github.com/jobrunner/fieldworksdiary-designsystem`
- Go-Version in `go.mod`: `go 1.26.0` (wie Ortus, Tempus, Situs, Hostus)
- Keine externen Go-Abhängigkeiten. Keine externen CSS-Abhängigkeiten. Kein Bundler, kein npm.
- Text-Tokens erreichen ≥ 7:1 gegen `--bg` **und** `--card`, in beiden Themen.
- `--control-line` erreicht ≥ 3:1 gegen `--bg` und `--card`, in beiden Themen.
- Weißer Text auf `--accent` erreicht ≥ 7:1 im hellen Thema.
- Alle Kommentare in CSS und Go auf Deutsch, im Stil der bestehenden Dienste: sie begründen, warum etwas so ist, nicht was dasteht.
- Commits auf Deutsch, Conventional-Commits-Präfix (`feat:`, `test:`, `docs:`, `chore:`).

---

## Dateistruktur

| Datei | Verantwortung |
|---|---|
| `go.mod` | Modulpfad und Go-Version |
| `designsystem.go` | `embed.FS` und die drei Ausgabefunktionen — sonst nichts |
| `css/tokens.css` | ausschließlich Variablen, hell und dunkel; keine Selektoren mit Wirkung |
| `css/base.css` | Reset und Komponenten; greift nur auf Tokens zu, enthält keine Farbliterale |
| `contrast.go` | Farbwert-Parser und WCAG-Leuchtdichte — von Test **und** Demo nutzbar |
| `contrast_test.go` | die Kontrastanforderungen als Tabelle |
| `css_test.go` | prüft, dass `base.css` kein Farbliteral enthält und nur definierte Tokens nutzt |
| `demo/index.html` | Referenzseite, alle Komponenten, beide Themen |
| `README.md` | Zweck, Einbindung, Status |

`contrast.go` liegt bewusst im Produktivcode, nicht in der Testdatei: die Demo-Seite soll die gemessenen Werte später anzeigen können, und ein Kontrastrechner ist keine Testhilfe, sondern die Durchsetzung einer Zusage des Moduls.

---

### Aufgabe 1: Kontrastrechner

**Dateien:**
- Anlegen: `go.mod`
- Anlegen: `contrast.go`
- Test: `contrast_test.go`

**Schnittstellen:**
- Liefert: `func ParseHex(s string) (RGB, error)`, `func (c RGB) Luminance() float64`, `func ContrastRatio(a, b RGB) float64` — von Aufgabe 2 zur Prüfung der Tokens genutzt.

- [ ] **Schritt 1: Modul anlegen**

```bash
cd /Users/jbrunner/work/projects/fieldworksdiary-designsystem
cat > go.mod <<'EOF'
module github.com/jobrunner/fieldworksdiary-designsystem

go 1.26.0
EOF
```

- [ ] **Schritt 2: Den fehlschlagenden Test schreiben**

Die Prüfwerte stammen aus der WCAG-Definition und sind von Hand nachrechenbar: Schwarz auf Weiß ergibt genau 21, eine Farbe gegen sich selbst genau 1.

`contrast_test.go`:

```go
package designsystem

import (
	"math"
	"testing"
)

func TestContrastRatio(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want float64
	}{
		{"Schwarz auf Weiß ist das Maximum", "#000000", "#ffffff", 21},
		{"eine Farbe gegen sich selbst", "#2563eb", "#2563eb", 1},
		{"Reihenfolge ist ohne Belang", "#ffffff", "#000000", 21},
		{"Kurzschreibweise wird verstanden", "#fff", "#000", 21},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := ParseHex(tt.a)
			if err != nil {
				t.Fatalf("ParseHex(%q): %v", tt.a, err)
			}
			b, err := ParseHex(tt.b)
			if err != nil {
				t.Fatalf("ParseHex(%q): %v", tt.b, err)
			}
			got := ContrastRatio(a, b)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("ContrastRatio = %.2f, erwartet %.2f", got, tt.want)
			}
		})
	}
}

func TestParseHexLehntUngueltigesAb(t *testing.T) {
	for _, s := range []string{"", "#", "#12345", "#gggggg", "2563eb "} {
		if _, err := ParseHex(s); err == nil {
			t.Errorf("ParseHex(%q) hat keinen Fehler geliefert", s)
		}
	}
}
```

- [ ] **Schritt 3: Test laufen lassen, Fehlschlag bestätigen**

Ausführen: `go test ./... -run TestContrast -v`
Erwartet: Übersetzungsfehler — `undefined: ParseHex`.

- [ ] **Schritt 4: Die Umsetzung schreiben**

`contrast.go`:

```go
// Package designsystem liefert die gemeinsame Gestaltungsgrundlage der
// fieldworksdiary-Dienste als CSS und hält die Kontrastzusagen des Systems
// prüfbar.
package designsystem

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// RGB ist eine Farbe im sRGB-Raum, je Kanal 0..255.
type RGB struct{ R, G, B uint8 }

// ParseHex liest #rgb und #rrggbb. Andere Schreibweisen werden abgelehnt
// statt stillschweigend gedeutet: ein Tippfehler in einem Farbwert soll den
// Build brechen, nicht zu einer zufälligen Farbe führen.
func ParseHex(s string) (RGB, error) {
	if !strings.HasPrefix(s, "#") {
		return RGB{}, fmt.Errorf("Farbwert %q beginnt nicht mit #", s)
	}
	h := s[1:]
	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	if len(h) != 6 {
		return RGB{}, fmt.Errorf("Farbwert %q hat weder 3 noch 6 Stellen", s)
	}
	n, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return RGB{}, fmt.Errorf("Farbwert %q ist nicht hexadezimal: %w", s, err)
	}
	return RGB{R: uint8(n >> 16), G: uint8(n >> 8), B: uint8(n)}, nil
}

// Luminance ist die relative Leuchtdichte nach WCAG 2.x.
func (c RGB) Luminance() float64 {
	lin := func(v uint8) float64 {
		f := float64(v) / 255
		if f <= 0.03928 {
			return f / 12.92
		}
		return math.Pow((f+0.055)/1.055, 2.4)
	}
	return 0.2126*lin(c.R) + 0.7152*lin(c.G) + 0.0722*lin(c.B)
}

// ContrastRatio liefert das Kontrastverhältnis zweier Farben, immer >= 1.
// Die Reihenfolge der Argumente spielt keine Rolle.
func ContrastRatio(a, b RGB) float64 {
	l1, l2 := a.Luminance(), b.Luminance()
	if l2 > l1 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}
```

- [ ] **Schritt 5: Test laufen lassen, Erfolg bestätigen**

Ausführen: `go test ./... -v`
Erwartet: PASS für alle Unterfälle.

- [ ] **Schritt 6: Einchecken**

```bash
git add go.mod contrast.go contrast_test.go
git commit -m "feat: Kontrastrechner nach WCAG

Die Kontrastzusagen des Systems sollen im Build nachgerechnet werden,
nicht in einer Tabelle behauptet. Der Rechner liegt im Produktivcode,
weil die Demo-Seite die Werte später anzeigen soll."
```

---

### Aufgabe 2: Farb-Tokens samt Prüfung

**Dateien:**
- Anlegen: `css/tokens.css`
- Anlegen: `designsystem.go`
- Test: `tokens_test.go`

**Schnittstellen:**
- Nutzt: `ParseHex`, `ContrastRatio` aus Aufgabe 1.
- Liefert: `func TokensCSS() []byte`; die Token-Namen `--bg`, `--card`, `--text`, `--text-muted`, `--text-disabled`, `--accent`, `--success`, `--error`, `--warning`, `--border`, `--control-line`, `--radius`, `--radius-sm`, `--shadow`, `--font`, `--font-mono` — Aufgabe 3 und alle Dienste bauen darauf.

- [ ] **Schritt 1: Den fehlschlagenden Test schreiben**

Der Test liest die Werte aus der CSS-Datei. Stünden sie zusätzlich im Test, prüfte er eine Kopie statt der Wahrheit.

`tokens_test.go`:

```go
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
```

- [ ] **Schritt 2: Test laufen lassen, Fehlschlag bestätigen**

Ausführen: `go test ./... -run TestTokens -v`
Erwartet: Übersetzungsfehler — `undefined: TokensCSS`.

- [ ] **Schritt 3: Das Einbetten schreiben**

`designsystem.go`:

```go
package designsystem

import _ "embed"

//go:embed css/tokens.css
var tokensCSS []byte

// TokensCSS liefert ausschließlich die Variablen — Farben für beide Themen,
// Typografie, Abstände, Radien. Wer nur die Werte braucht und seine
// Komponenten selbst mitbringt, bindet allein diese Datei ein.
func TokensCSS() []byte { return tokensCSS }
```

- [ ] **Schritt 4: Die Tokens schreiben**

`css/tokens.css`:

```css
/* Gemeinsame Gestaltungswerte der fieldworksdiary-Dienste.
   Diese Datei enthält ausschließlich Variablen — kein Selektor hier
   verändert die Darstellung. Wer Komponenten will, bindet base.css dazu.

   Alle Textfarben erreichen 7:1 gegen --bg UND --card (WCAG 2.2, SC 1.4.6).
   Die Dienste werden im Gelände auf Handydisplays bei Tageslicht gelesen;
   die 4.5:1 der Stufe AA reichen dafür erfahrungsgemäß nicht.
   contrast_test.go rechnet jeden Wert im Build nach. */

:root {
  /* Flächen */
  --bg: #f8fafc;
  --card: #ffffff;

  /* Text */
  --text: #1e293b;
  --text-muted: #475569;
  /* Gleich dem gedämpften Text: ein gesperrtes Bedienelement darf nicht
     durch schwachen Kontrast erkennbar gemacht werden, sondern durch
     aria-disabled, Mauszeiger und fehlende Reaktion. */
  --text-disabled: #475569;

  /* Akzent — trägt weiße Schrift auf der Fläche des primären Knopfes */
  --accent: #1e40af;
  --accent-hover: #1e3a8a;

  /* Status */
  --success: #166534;
  --error: #991b1b;
  --warning: #92400e;

  /* Linien. --border ist dekorativ und darf zart bleiben; --control-line
     begrenzt Bedienelemente und macht sie überhaupt erst erkennbar.
     WCAG 1.4.11 fordert nur für Letztere einen Kontrast (>= 3:1) —
     dekorative Trennlinien sind ausdrücklich ausgenommen. Getrennte
     Tokens, weil ein gemeinsamer Wert entweder die Linien zu hart oder
     die Bedienelemente zu schwach macht. */
  --border: #e2e8f0;
  --control-line: #767676;

  /* Form */
  --radius: 8px;
  --radius-sm: 4px;
  --shadow: 0 1px 3px rgba(0, 0, 0, 0.1);

  /* Schrift */
  --font: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  --font-mono: "SF Mono", Monaco, ui-monospace, monospace;
}

@media (prefers-color-scheme: dark) {
  :root {
    --bg: #0f172a;
    --card: #1e293b;

    --text: #f1f5f9;
    --text-muted: #b4c0ce;
    --text-disabled: #b4c0ce;

    /* Im dunklen Thema ist der Akzent eine Textfarbe, keine Fläche:
       ein Knopf bekommt hier dunkle Schrift auf hellem Blau. */
    --accent: #93c5fd;
    --accent-hover: #bfdbfe;

    --success: #86efac;
    --error: #fca5a5;
    --warning: #fcd34d;

    --border: #334155;
    --control-line: #8695a8;

    /* Auf dunklem Grund trägt ein Schatten nicht: --card steht gegen --bg
       nur bei 1.22:1. Karten heben sich hier über ihren Rand ab, nicht
       über den Schatten — siehe .card in base.css. */
    --shadow: none;
  }
}
```

- [ ] **Schritt 5: Test laufen lassen, Erfolg bestätigen**

Ausführen: `go test ./... -v`
Erwartet: PASS. Schlägt ein Kontrast fehl, nennt die Meldung Token, Wert, Bezugsfläche und gemessenes Verhältnis — der Wert wird in `tokens.css` korrigiert, nicht der Test.

- [ ] **Schritt 6: Einchecken**

```bash
git add css/tokens.css designsystem.go tokens_test.go
git commit -m "feat: Farb-Tokens für helles und dunkles Thema

Alle Textfarben auf 7:1 gegen Seitenhintergrund und Kartenfläche. Der
Test liest die Werte aus tokens.css statt sie zu wiederholen: sonst
prüfte er eine Kopie und nicht das, was ausgeliefert wird."
```

---

### Aufgabe 3: Basiskomponenten

**Dateien:**
- Anlegen: `css/base.css`
- Ändern: `designsystem.go`
- Test: `css_test.go`

**Schnittstellen:**
- Nutzt: die Token-Namen aus Aufgabe 2.
- Liefert: `func BaseCSS() []byte`, `func CSS() []byte`; die Klassen `.container`, `.card`, `.card-title`, `.form-group`, `.btn`, `.btn-secondary`, `.btn-row`, `.tabs`, `.tab`, `.badge`, `.badge-success`, `.badge-error`, `.table-wrap`, `.spinner`, `.loading`, `.error`, `.sr-only`, `.skip-link`, `.muted` — Expertus und die übrigen Dienste bauen darauf.

- [ ] **Schritt 1: Den fehlschlagenden Test schreiben**

Zwei Zusagen werden geprüft: `base.css` enthält keine Farbe im Klartext, und es verweist auf kein Token, das es nicht gibt. Beides verhindert genau den Rückfall, der die Prüfung aus Aufgabe 2 wertlos machte.

`css_test.go`:

```go
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
```

- [ ] **Schritt 2: Test laufen lassen, Fehlschlag bestätigen**

Ausführen: `go test ./... -run TestBase -v`
Erwartet: Übersetzungsfehler — `undefined: BaseCSS`.

- [ ] **Schritt 3: Die Ausgabefunktionen ergänzen**

`designsystem.go` vollständig:

```go
package designsystem

import _ "embed"

//go:embed css/tokens.css
var tokensCSS []byte

//go:embed css/base.css
var baseCSS []byte

// TokensCSS liefert ausschließlich die Variablen — Farben für beide Themen,
// Typografie, Abstände, Radien. Wer nur die Werte braucht und seine
// Komponenten selbst mitbringt, bindet allein diese Datei ein.
func TokensCSS() []byte { return tokensCSS }

// BaseCSS liefert Reset und Komponenten. Setzt die Tokens voraus und ist
// ohne sie wirkungslos — jede Farbe darin ist eine Variable.
func BaseCSS() []byte { return baseCSS }

// CSS liefert beides in der einzig gültigen Reihenfolge: erst die
// Variablen, dann was sie verwendet.
func CSS() []byte {
	out := make([]byte, 0, len(tokensCSS)+len(baseCSS)+1)
	out = append(out, tokensCSS...)
	out = append(out, '\n')
	return append(out, baseCSS...)
}
```

- [ ] **Schritt 4: Die Komponenten schreiben**

`css/base.css`:

```css
/* Reset und Basiskomponenten der fieldworksdiary-Dienste.
   Setzt tokens.css voraus: jede Farbe hier ist eine Variable, kein Wert.
   Aufgenommen ist, was mehr als ein Dienst braucht — Fachliches wie die
   Wetterkarte in Tempus oder die Quellenliste in Ortus bleibt dort. */

*,
*::before,
*::after {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: var(--font);
  font-size: 1rem;
  line-height: 1.5;
  color: var(--text);
  background: var(--bg);
  min-height: 100vh;
}

.container {
  max-width: 48rem;
  margin-inline: auto;
  padding: 1rem;
}

/* Sichtbarer Fokus ist Pflicht (WCAG 2.4.7). Der Browser-Vorgabering
   verschwindet auf farbigen Flächen, deshalb ein eigener mit Abstand —
   der Abstand hält ihn auch auf einem Knopf in Akzentfarbe sichtbar. */
:focus-visible {
  outline: 3px solid var(--accent);
  outline-offset: 2px;
}

/* --- Text --- */

.muted {
  color: var(--text-muted);
}

/* --- Flächen --- */

.card {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 1.25rem;
  margin-bottom: 1rem;
}

/* Der Rand steht hier immer, nicht nur im dunklen Thema: auf hellem Grund
   trägt der Schatten die Abhebung und der Rand stört nicht, auf dunklem
   ist er die einzige Trennung (--card gegen --bg nur 1.22:1) und
   --shadow ist dort none. Eine Regel für beide Themen statt zweier, die
   auseinanderlaufen können. */

.card-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0 0 1rem;
}

/* --- Formular --- */

.form-group {
  margin-bottom: 1rem;
}

/* Als Block, nicht inline: sonst klebt die Beschriftung unmittelbar am
   Feld und beide verschmelzen optisch zu einem Wort. */
label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 0.375rem;
  color: var(--text);
}

input,
select,
textarea {
  width: 100%;
  /* Ein <select> bemisst sich sonst an seinem längsten Eintrag und sprengt
     bei 200 % Textgröße seine Spalte — das erzeugt waagerechtes Rollen der
     ganzen Seite (WCAG 1.4.4 und 1.4.10). Felder bleiben deshalb an ihren
     Rahmen gebunden, nicht an ihren Inhalt. */
  max-width: 100%;
  min-height: 2.75rem;
  padding: 0.625rem 0.75rem;
  font: inherit;
  color: var(--text);
  background: var(--card);
  border: 1px solid var(--control-line);
  border-radius: var(--radius);
  transition: border-color 0.15s, box-shadow 0.15s;
}

input:focus,
select:focus,
textarea:focus {
  border-color: var(--accent);
}

input::placeholder,
textarea::placeholder {
  color: var(--text-muted);
}

textarea {
  font-family: var(--font-mono);
  font-size: 0.875rem;
  resize: vertical;
  min-height: 10rem;
}

.checkbox-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-height: 2.75rem;
  font-weight: 400;
  cursor: pointer;
}

.checkbox-row input[type="checkbox"] {
  width: 1.05rem;
  height: 1.05rem;
  min-height: 0;
  margin: 0;
  accent-color: var(--accent);
}

/* --- Knöpfe --- */

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  /* 44 px, wo es geht; 24 px sind die harte Untergrenze (WCAG 2.5.8). */
  min-height: 2.75rem;
  padding: 0.75rem 1rem;
  font: inherit;
  font-weight: 500;
  color: var(--card);
  background: var(--accent);
  border: 1px solid var(--accent);
  border-radius: var(--radius);
  cursor: pointer;
}

.btn:hover {
  background: var(--accent-hover);
  border-color: var(--accent-hover);
}

.btn:disabled {
  color: var(--text-disabled);
  background: var(--bg);
  border-color: var(--control-line);
  cursor: not-allowed;
}

.btn-secondary {
  color: var(--text);
  background: var(--card);
  border-color: var(--control-line);
}

.btn-secondary:hover {
  background: var(--bg);
  border-color: var(--text-muted);
}

.btn-row {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

/* --- Reiter --- */

.tabs {
  display: flex;
  gap: 0.25rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid var(--border);
}

.tab {
  min-height: 2.75rem;
  padding: 0.625rem 1rem;
  font: inherit;
  font-weight: 500;
  color: var(--text-muted);
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
}

.tab:hover {
  color: var(--text);
}

.tab[aria-selected="true"] {
  color: var(--accent);
  border-bottom-color: var(--accent);
}

/* --- Marken --- */

.badge {
  display: inline-flex;
  align-items: center;
  white-space: nowrap;
  padding: 0.125rem 0.5rem;
  font-size: 0.75rem;
  font-weight: 500;
  border-radius: 9999px;
  color: var(--text);
  border: 1px solid var(--control-line);
}

.badge-success {
  color: var(--success);
  border-color: var(--success);
}

.badge-error {
  color: var(--error);
  border-color: var(--error);
}

/* --- Tabellen --- */

/* Eine breite Tabelle rollt in ihrem eigenen Kasten; die Seite selbst darf
   nicht waagerecht rollen (WCAG 1.4.10). */
.table-wrap {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  text-align: left;
  padding: 0.5rem;
  border-bottom: 1px solid var(--border);
}

thead th {
  font-weight: 600;
  color: var(--text-muted);
}

/* --- Zustände --- */

.loading {
  display: none;
  text-align: center;
  padding: 2rem;
  color: var(--text-muted);
}

.loading.active {
  display: block;
}

.spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 0.5rem;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error {
  color: var(--error);
  border: 1px solid var(--error);
  border-radius: var(--radius);
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

/* --- Nur für Screenreader --- */

/* display:none oder visibility:hidden wären falsch: sie nehmen das Element
   auch aus dem Accessibility-Baum, die Ansage fiele aus. */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  margin: -1px;
  padding: 0;
  border: 0;
  clip-path: inset(50%);
  overflow: hidden;
  white-space: nowrap;
}

.skip-link {
  position: absolute;
  left: -999px;
}

.skip-link:focus {
  left: 0.5rem;
  top: 0.5rem;
  z-index: 1;
  padding: 0.5rem;
  color: var(--text);
  background: var(--card);
  border: 1px solid var(--control-line);
  border-radius: var(--radius-sm);
}

/* --- Größere Schirme --- */

@media (min-width: 640px) {
  .container {
    padding: 2rem;
  }

  .card {
    padding: 1.5rem;
  }
}

/* --- Weniger Bewegung --- */

/* Wer das im Betriebssystem einstellt, tut das oft wegen eines
   Gleichgewichtsleidens — eine Drehung genügt, um Übelkeit auszulösen. */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

- [ ] **Schritt 5: Test laufen lassen, Erfolg bestätigen**

Ausführen: `go test ./... -v`
Erwartet: PASS. Meldet der Test ein Farbliteral, wandert der Wert nach `tokens.css`; meldet er ein unbekanntes Token, ist es dort zu ergänzen.

- [ ] **Schritt 6: Einchecken**

```bash
git add css/base.css designsystem.go css_test.go
git commit -m "feat: Basiskomponenten

Reset, Formularelemente, Knöpfe, Karten, Reiter, Marken, Tabellen.
Der Test verbietet Farbliterale in base.css: eine Farbe im Klartext
stünde außerhalb der Kontrastprüfung und höhlte sie aus."
```

---

### Aufgabe 4: Demo-Seite

**Dateien:**
- Anlegen: `demo/index.html`
- Ändern: `designsystem.go`
- Test: `demo_test.go`

**Schnittstellen:**
- Nutzt: `CSS()` aus Aufgabe 3.
- Liefert: `func DemoHandler() http.Handler` — dient der visuellen Abnahme und dem Einbindungsbeispiel für die Dienste.

- [ ] **Schritt 1: Den fehlschlagenden Test schreiben**

Geprüft wird, dass die Demo jede Komponente wirklich zeigt. Ohne diese Prüfung veraltet sie still, sobald eine Komponente dazukommt — und sie ist laut Spec für die Mehrzahl der Komponenten die einzige Prüfstelle vor dem produktiven Einsatz.

`demo_test.go`:

```go
package designsystem

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

// klassenRe findet die Klassennamen, die base.css definiert.
var klassenRe = regexp.MustCompile(`^\.([a-z][a-z0-9-]*)`)

func TestDemoZeigtJedeKomponente(t *testing.T) {
	// Welche Klassen base.css anbietet, steht in base.css — nicht in einer
	// Liste hier, die beim nächsten Zuwachs vergessen würde.
	angeboten := map[string]bool{}
	for _, zeile := range strings.Split(string(BaseCSS()), "\n") {
		if m := klassenRe.FindStringSubmatch(strings.TrimSpace(zeile)); m != nil {
			angeboten[m[1]] = true
		}
	}
	// Zustandsklassen erscheinen nur zusammen mit ihrer Grundklasse.
	for _, nur := range []string{"active"} {
		delete(angeboten, nur)
	}

	demo := string(demoHTML)
	for klasse := range angeboten {
		if !strings.Contains(demo, `"`+klasse+`"`) && !strings.Contains(demo, `"`+klasse+` `) &&
			!strings.Contains(demo, ` `+klasse+`"`) && !strings.Contains(demo, ` `+klasse+` `) {
			t.Errorf("Die Demo-Seite zeigt .%s nicht — die Klasse wäre nirgends vor dem Einsatz zu sehen", klasse)
		}
	}
}

func TestDemoHandlerLiefertSeiteUndCSS(t *testing.T) {
	srv := httptest.NewServer(DemoHandler())
	defer srv.Close()

	for _, f := range []struct {
		pfad string
		typ  string
		teil string
	}{
		{"/", "text/html; charset=utf-8", "<!doctype html>"},
		{"/designsystem.css", "text/css; charset=utf-8", "--control-line:"},
	} {
		res, err := http.Get(srv.URL + f.pfad)
		if err != nil {
			t.Fatalf("GET %s: %v", f.pfad, err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Errorf("GET %s: Status %d, erwartet 200", f.pfad, res.StatusCode)
		}
		if got := res.Header.Get("Content-Type"); got != f.typ {
			t.Errorf("GET %s: Content-Type %q, erwartet %q", f.pfad, got, f.typ)
		}
		buf := make([]byte, 4096)
		n, _ := res.Body.Read(buf)
		if !strings.Contains(strings.ToLower(string(buf[:n])), strings.ToLower(f.teil)) {
			t.Errorf("GET %s: %q kommt im Anfang der Antwort nicht vor", f.pfad, f.teil)
		}
	}
}
```

- [ ] **Schritt 2: Test laufen lassen, Fehlschlag bestätigen**

Ausführen: `go test ./... -run TestDemo -v`
Erwartet: Übersetzungsfehler — `undefined: demoHTML`, `undefined: DemoHandler`.

- [ ] **Schritt 3: Den Handler ergänzen**

An `designsystem.go` anhängen:

```go
//go:embed demo/index.html
var demoHTML []byte

// DemoHandler liefert die Referenzseite samt Stylesheet. Sie dient der
// visuellen Abnahme des Systems und ist zugleich das kürzeste Beispiel,
// wie ein Dienst das CSS einbindet:
//
//	go run ./cmd/demo   (oder im Test: httptest.NewServer(DemoHandler()))
func DemoHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/designsystem.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(CSS())
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(demoHTML)
	})
	return mux
}
```

Der Import-Block oben wird zu:

```go
import (
	_ "embed"
	"net/http"
)
```

- [ ] **Schritt 4: Die Demo-Seite schreiben**

`demo/index.html`:

```html
<!doctype html>
<html lang="de">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Design-System — Referenzseite</title>
    <link rel="stylesheet" href="/designsystem.css" />
  </head>
  <body>
    <div class="container">
      <h1>Design-System</h1>
      <p class="muted">
        Jede Komponente in ihrem Normalzustand. Das Thema folgt der
        Systemeinstellung — zum Prüfen des dunklen Themas dort umschalten.
      </p>

      <section class="card">
        <h2 class="card-title">Knöpfe</h2>
        <div class="btn-row">
          <button type="button" class="btn">Abfragen</button>
          <button type="button" class="btn btn-secondary">Leeren</button>
          <button type="button" class="btn" disabled>Gesperrt</button>
        </div>
      </section>

      <section class="card">
        <h2 class="card-title">Formular</h2>
        <div class="form-group">
          <label for="d-text">Beschriftung</label>
          <input type="text" id="d-text" placeholder="Platzhaltertext" />
        </div>
        <div class="form-group">
          <label for="d-select">Auswahl</label>
          <select id="d-select">
            <option>Erster Eintrag</option>
            <option>Ein deutlich längerer zweiter Eintrag</option>
          </select>
        </div>
        <div class="form-group">
          <label for="d-area">Mehrzeilig</label>
          <textarea id="d-area" placeholder="52.52, 13.405"></textarea>
        </div>
        <label class="checkbox-row" for="d-check">
          <input type="checkbox" id="d-check" checked />
          Ankreuzfeld in einer Zeile
        </label>
      </section>

      <section class="card">
        <h2 class="card-title">Reiter</h2>
        <div class="tabs" role="tablist" aria-label="Beispiel">
          <button type="button" class="tab" role="tab" aria-selected="true">Erster</button>
          <button type="button" class="tab" role="tab" aria-selected="false" tabindex="-1">Zweiter</button>
        </div>
      </section>

      <section class="card">
        <h2 class="card-title">Marken</h2>
        <p>
          <span class="badge">neutral</span>
          <span class="badge badge-success">erfolgreich</span>
          <span class="badge badge-error">fehlgeschlagen</span>
        </p>
      </section>

      <section class="card">
        <h2 class="card-title">Tabelle</h2>
        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">Feld</th>
                <th scope="col">Wert</th>
                <th scope="col">Herkunft</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <th scope="row">Country</th>
                <td>Deutschland</td>
                <td class="muted">abgeleitet</td>
              </tr>
              <tr>
                <th scope="row">Altitude (m)</th>
                <td>412</td>
                <td class="muted">abgeleitet</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section class="card">
        <h2 class="card-title">Zustände</h2>
        <p class="error" role="alert">Die Abfrage ist fehlgeschlagen.</p>
        <div class="loading active" role="status">
          <div class="spinner"></div>
          <p>Abfrage wird ausgeführt …</p>
        </div>
        <p class="muted">Gedämpfter Nebentext.</p>
      </section>

      <a class="skip-link" href="#ende">Sprunglink (nur bei Fokus sichtbar)</a>
      <p class="sr-only">Dieser Satz ist nur für Screenreader da.</p>
      <p id="ende" class="muted">Ende der Referenzseite.</p>
    </div>
  </body>
</html>
```

- [ ] **Schritt 5: Test laufen lassen, Erfolg bestätigen**

Ausführen: `go test ./... -v`
Erwartet: PASS. Meldet der Test eine fehlende Klasse, wird sie in die Demo-Seite aufgenommen — nicht aus der Prüfung ausgenommen.

- [ ] **Schritt 6: Die Seite mit eigenen Augen ansehen**

```bash
mkdir -p cmd/demo
cat > cmd/demo/main.go <<'EOF'
// Startet die Referenzseite des Design-Systems zur visuellen Abnahme.
package main

import (
	"log"
	"net/http"

	designsystem "github.com/jobrunner/fieldworksdiary-designsystem"
)

func main() {
	log.Println("Referenzseite auf http://127.0.0.1:5180")
	log.Fatal(http.ListenAndServe("127.0.0.1:5180", designsystem.DemoHandler()))
}
EOF
go run ./cmd/demo
```

Im Browser `http://127.0.0.1:5180` öffnen, einmal im hellen und einmal im dunklen Systemthema durchsehen. **Das ist der Abnahmepunkt aus der Spec** — für die Mehrzahl der Komponenten die einzige Prüfstelle, bevor sie produktiven Code berühren. Auffälligkeiten werden hier behoben, nicht später in einem Dienst.

- [ ] **Schritt 7: Einchecken**

```bash
git add demo/ cmd/ designsystem.go demo_test.go
git commit -m "feat: Referenzseite zur visuellen Abnahme

Zeigt jede Komponente in beiden Themen. Der Test leitet die erwartete
Liste aus base.css ab statt sie zu wiederholen: eine neue Komponente
fehlt sonst still auf der Seite, die sie eigentlich prüfen soll."
```

---

### Aufgabe 5: Einbindung dokumentieren und veröffentlichen

**Dateien:**
- Ändern: `README.md`
- Anlegen: `.github/workflows/ci.yml`

**Schnittstellen:**
- Nutzt: alles Vorherige. Liefert: die Fassung `v0.1.0`, auf die Expertus sich bezieht.

- [ ] **Schritt 1: Den Ablauf für die Fortlaufende Integration anlegen**

`.github/workflows/ci.yml`:

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:

permissions:
  contents: read

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - name: Tests
        run: go test ./... -v
      - name: Formatierung
        run: test -z "$(gofmt -l .)"
```

- [ ] **Schritt 2: Das README fertigschreiben**

```markdown
# fieldworksdiary-designsystem

Gemeinsame Gestaltungsgrundlage der fieldworksdiary-Dienste — Ortus, Tempus,
Situs, Hostus und Expertus.

## Einbinden

    go get github.com/jobrunner/fieldworksdiary-designsystem

Das CSS unter einer eigenen Route ausliefern:

    import designsystem "github.com/jobrunner/fieldworksdiary-designsystem"

    mux.HandleFunc("/assets/designsystem.css", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/css; charset=utf-8")
        w.Write(designsystem.CSS())
    })

Im HTML verlinken:

    <link rel="stylesheet" href="/assets/designsystem.css">

`TokensCSS()` und `BaseCSS()` liefern die Teile einzeln, falls ein Dienst nur
die Variablen braucht.

## Was drin ist

`css/tokens.css` — Farben für helles und dunkles Thema, Typografie, Abstände,
Radien. `css/base.css` — Reset, Formularelemente, Knöpfe, Karten, Reiter,
Marken, Tabellen, Zustände.

Fachliche Komponenten gehören nicht hierher: die Wetterkarte bleibt in Tempus,
die Quellenliste in Ortus.

## Zusagen

Alle Textfarben erreichen 7:1 gegen Seitenhintergrund und Kartenfläche, in
beiden Themen. Ränder von Bedienelementen erreichen 3:1. `go test ./...`
rechnet das nach — die Werte werden aus `css/tokens.css` gelesen, nicht im
Testcode wiederholt.

`base.css` darf keine Farbe im Klartext enthalten; auch das prüft der Test.
Eine Farbe außerhalb von `tokens.css` stünde außerhalb der Kontrastprüfung.

## Referenzseite

    go run ./cmd/demo

Zeigt jede Komponente unter <http://127.0.0.1:5180>. Das Thema folgt der
Systemeinstellung.

## Entwurf

[docs/design.md](docs/design.md) — warum das System so aussieht und wie die
fünf Dienste darauf umgestellt werden.
```

- [ ] **Schritt 3: Alles prüfen**

```bash
go test ./... -v
gofmt -l .
```

Erwartet: alle Tests PASS, `gofmt -l` gibt nichts aus.

- [ ] **Schritt 4: Einchecken und veröffentlichen**

```bash
git add README.md .github/
git commit -m "docs: Einbindung beschreiben, CI ergänzen"
git push
git tag v0.1.0
git push origin v0.1.0
```

Die Marke `v0.1.0` ist die Fassung, auf die Expertus sich im zweiten Plan bezieht. Ohne sie zieht `go get` einen Pseudo-Versionsstring aus dem letzten Commit — das funktioniert, macht die Abhängigkeit aber unlesbar.

---

## Abschluss

Nach Aufgabe 5 steht ein Modul, das für sich lauffähig und geprüft ist:

- `go test ./...` rechnet jeden Farbwert nach
- die Referenzseite zeigt jede Komponente in beiden Themen
- `v0.1.0` ist veröffentlicht

Erst danach beginnt der zweite Plan — Expertus.
