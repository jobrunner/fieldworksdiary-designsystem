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

// cssNamedColors ist die vollständige Liste der CSS-Farbschlüsselwörter
// (CSS Color Module, "extended color keywords", 148 Namen). Eine Auswahl
// wie zuvor (19 Namen) lässt jede Farbe außerhalb der Auswahl durch —
// nachgewiesen etwa mit "crimson", "tomato" oder "darkblue". transparent
// ist kein Farbwert im Sinne der Zusage und wird gesondert behandelt.
const cssNamedColors = `aliceblue|antiquewhite|aqua|aquamarine|azure|beige|bisque|black|blanchedalmond|blue|blueviolet|brown|burlywood|cadetblue|chartreuse|chocolate|coral|cornflowerblue|cornsilk|crimson|cyan|darkblue|darkcyan|darkgoldenrod|darkgray|darkgreen|darkgrey|darkkhaki|darkmagenta|darkolivegreen|darkorange|darkorchid|darkred|darksalmon|darkseagreen|darkslateblue|darkslategray|darkslategrey|darkturquoise|darkviolet|deeppink|deepskyblue|dimgray|dimgrey|dodgerblue|firebrick|floralwhite|forestgreen|fuchsia|gainsboro|ghostwhite|gold|goldenrod|gray|grey|green|greenyellow|honeydew|hotpink|indianred|indigo|ivory|khaki|lavender|lavenderblush|lawngreen|lemonchiffon|lightblue|lightcoral|lightcyan|lightgoldenrodyellow|lightgray|lightgreen|lightgrey|lightpink|lightsalmon|lightseagreen|lightskyblue|lightslategray|lightslategrey|lightsteelblue|lightyellow|lime|limegreen|linen|magenta|maroon|mediumaquamarine|mediumblue|mediumorchid|mediumpurple|mediumseagreen|mediumslateblue|mediumspringgreen|mediumturquoise|mediumvioletred|midnightblue|mintcream|mistyrose|moccasin|navajowhite|navy|oldlace|olive|olivedrab|orange|orangered|orchid|palegoldenrod|palegreen|paleturquoise|palevioletred|papayawhip|peachpuff|peru|pink|plum|powderblue|purple|rebeccapurple|red|rosybrown|royalblue|saddlebrown|salmon|sandybrown|seagreen|seashell|sienna|silver|skyblue|slateblue|slategray|slategrey|snow|springgreen|steelblue|tan|teal|thistle|tomato|turquoise|violet|wheat|white|whitesmoke|yellow|yellowgreen`

var (
	hexRe = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	// funcColorRe erfasst die klassischen UND die modernen CSS-Farbfunktionen.
	// oklch(), lab(), lch(), color-mix() und color() umgehen die
	// Kontrastprüfung genauso wie rgb()/hsl() — nachgewiesen kamen sie
	// bislang unbemerkt durch.
	funcColorRe = regexp.MustCompile(`\brgba?\(|\bhsla?\(|\boklch\(|\boklab\(|\blab\(|\blch\(|\bcolor-mix\(|\bcolor\(`)
	// uriHexRe erfasst ein prozentkodiertes "#" (%23) gefolgt von Hex-Ziffern
	// — der übliche Weg, einem <select> per data-URI (etwa
	// background-image: url("data:image/svg+xml,...fill='%23ff0000'..."))
	// einen eigenen Pfeil mit fest kodierter Farbe zu geben. Ein rohes "#"
	// kommt darin nicht vor, hexRe findet diesen Fall deshalb nicht.
	uriHexRe     = regexp.MustCompile(`%23[0-9a-fA-F]{3,8}`)
	namedColorRe = regexp.MustCompile(`(?:^|[^\w-])(` + cssNamedColors + `)(?:[^\w-]|$)`)
	varNutzungRe = regexp.MustCompile(`var\((--[a-z-]+)`)
	varDefRe     = regexp.MustCompile(`(--[a-z-]+):`)
)

// farbliterale prüft CSS-Text auf Farbliterale und gibt Fundstellen mit
// Zeilennummern zurück. Kommentare werden ignoriert. transparent ist erlaubt,
// alle anderen benannten Farben, Hex-Werte, Farbfunktionen (klassisch und
// modern) sowie prozentkodierte Hex-Werte in data-URIs gelten als Verstoß.
func farbliterale(css string) []fund {
	// Kommentare entfernen, aber Zeilennummern durch Zeilenumbrüche erhalten
	kommentarlos := removeComments(css)

	var funde []fund
	for i, zeile := range strings.Split(kommentarlos, "\n") {
		lineNum := i + 1

		// Hex-Farben prüfen — ALLE auf der Zeile, nicht nur die erste.
		for _, m := range hexRe.FindAllString(zeile, -1) {
			funde = append(funde, fund{zeile: lineNum, wert: m})
		}

		// Prozentkodierte Hex-Farben in data-URIs prüfen (%23... statt #...).
		for _, m := range uriHexRe.FindAllString(zeile, -1) {
			funde = append(funde, fund{zeile: lineNum, wert: m})
		}

		// rgb/rgba/hsl/hsla sowie moderne Farbfunktionen prüfen
		if funcColorRe.MatchString(zeile) {
			funde = append(funde, fund{zeile: lineNum, wert: "rgb/rgba/hsl/hsla/oklch/lab/lch/color-mix/color"})
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

// TestUnbekannteBenannteFarbenWerdenErkannt belegt, dass die vollständige
// Farbwortliste greift, nicht nur eine Auswahl der bekanntesten Namen —
// nachgewiesen durchgelassen wurden zuvor u. a. "crimson", "tomato" und
// "darkblue".
func TestUnbekannteBenannteFarbenWerdenErkannt(t *testing.T) {
	for _, farbe := range []string{"crimson", "tomato", "slategray", "darkblue"} {
		funde := farbliterale("body { color: " + farbe + "; }")
		if len(funde) == 0 {
			t.Errorf("Test erkennt benannte Farbe %q nicht — sollte erkannt werden", farbe)
		}
	}
}

// TestModerneFarbfunktionenWerdenErkannt belegt, dass oklch(), lab(), lch()
// und color-mix() erkannt werden — die alte Prüfung kannte nur rgb()/hsl()
// und ließ diese moderneren Schreibweisen unbemerkt durch.
func TestModerneFarbfunktionenWerdenErkannt(t *testing.T) {
	for _, wert := range []string{
		"oklch(59% 0.15 250)",
		"lab(50% 40 20)",
		"lch(50% 40 20)",
		"color-mix(in srgb, red 50%, blue 50%)",
		"color(display-p3 1 0 0)",
	} {
		funde := farbliterale("body { color: " + wert + "; }")
		if len(funde) == 0 {
			t.Errorf("Test erkennt Farbfunktion %q nicht — sollte erkannt werden", wert)
		}
	}
}

// TestProzentkodierteHexFarbeInDataUriWirdErkannt belegt, dass ein
// prozentkodiertes "#" in einer data-URI erkannt wird — der übliche Weg,
// einem <select> einen eigenen Pfeil mit fest kodierter Farbe zu geben, und
// damit der wahrscheinlichste reale Fall.
func TestProzentkodierteHexFarbeInDataUriWirdErkannt(t *testing.T) {
	css := `select { background-image: url("data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg'><path fill='%23ff0000'/></svg>"); }`
	funde := farbliterale(css)
	if len(funde) == 0 {
		t.Error("Test erkennt prozentkodierte Hex-Farbe (%23...) in data-URI nicht — sollte erkannt werden")
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

// TestAlleHexwerteEinerZeileWerdenGemeldet belegt, dass die Meldung nicht
// nur den ERSTEN Hex-Wert je Zeile nennt — ein zweiter Verstoß in derselben
// Zeile blieb zuvor unerwähnt.
func TestAlleHexwerteEinerZeileWerdenGemeldet(t *testing.T) {
	funde := farbliterale("body { border: 1px solid #ff0000; background: #00ff00; }")
	if len(funde) != 2 {
		t.Fatalf("erwartet 2 Funde (beide Hex-Werte der Zeile), got %d", len(funde))
	}
	werte := map[string]bool{funde[0].wert: true, funde[1].wert: true}
	for _, erwartet := range []string{"#ff0000", "#00ff00"} {
		if !werte[erwartet] {
			t.Errorf("erwartet Fund für %q, nicht dabei: %v", erwartet, funde)
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
