// Package farbe erkennt Farbliterale im Klartext — Hex-Werte, die
// klassischen und modernen CSS-Farbfunktionen, prozentkodierte Hex-Werte in
// data-URIs und die vollständige Liste der CSS-Farbschlüsselwörter
// (148 Namen, CSS Color Module, "extended color keywords"). "transparent"
// ist kein Farbwert im Sinne der Zusage und zählt nicht als Fund.
//
// Zwei Prüfungen im Modul brauchen dieselbe Strenge: css_test.go für
// base.css und icons_test.go für die SVG-Symbole. Eine Symbolfunktion, die
// stroke="black" statt currentColor schreibt, umgeht die Kontrastprüfung
// genauso wie eine Farbe im Klartext in base.css — die Erkennung gehört
// deshalb an eine Stelle statt zweimal gepflegt zu werden.
package farbe

import "regexp"

// NamedColors ist die vollständige Liste der CSS-Farbschlüsselwörter. Eine
// Auswahl weniger, bekannter Namen ließe jede Farbe außerhalb der Auswahl
// unbemerkt durch — nachgewiesen etwa mit "crimson", "tomato" oder
// "darkblue" in css_test.go, und mit "black" in icons_test.go.
const NamedColors = `aliceblue|antiquewhite|aqua|aquamarine|azure|beige|bisque|black|blanchedalmond|blue|blueviolet|brown|burlywood|cadetblue|chartreuse|chocolate|coral|cornflowerblue|cornsilk|crimson|cyan|darkblue|darkcyan|darkgoldenrod|darkgray|darkgreen|darkgrey|darkkhaki|darkmagenta|darkolivegreen|darkorange|darkorchid|darkred|darksalmon|darkseagreen|darkslateblue|darkslategray|darkslategrey|darkturquoise|darkviolet|deeppink|deepskyblue|dimgray|dimgrey|dodgerblue|firebrick|floralwhite|forestgreen|fuchsia|gainsboro|ghostwhite|gold|goldenrod|gray|grey|green|greenyellow|honeydew|hotpink|indianred|indigo|ivory|khaki|lavender|lavenderblush|lawngreen|lemonchiffon|lightblue|lightcoral|lightcyan|lightgoldenrodyellow|lightgray|lightgreen|lightgrey|lightpink|lightsalmon|lightseagreen|lightskyblue|lightslategray|lightslategrey|lightsteelblue|lightyellow|lime|limegreen|linen|magenta|maroon|mediumaquamarine|mediumblue|mediumorchid|mediumpurple|mediumseagreen|mediumslateblue|mediumspringgreen|mediumturquoise|mediumvioletred|midnightblue|mintcream|mistyrose|moccasin|navajowhite|navy|oldlace|olive|olivedrab|orange|orangered|orchid|palegoldenrod|palegreen|paleturquoise|palevioletred|papayawhip|peachpuff|peru|pink|plum|powderblue|purple|rebeccapurple|red|rosybrown|royalblue|saddlebrown|salmon|sandybrown|seagreen|seashell|sienna|silver|skyblue|slateblue|slategray|slategrey|snow|springgreen|steelblue|tan|teal|thistle|tomato|turquoise|violet|wheat|white|whitesmoke|yellow|yellowgreen`

var (
	HexRe = regexp.MustCompile(`#[0-9a-fA-F]{3,8}\b`)
	// FuncColorRe erfasst die klassischen UND die modernen CSS-Farbfunktionen.
	FuncColorRe = regexp.MustCompile(`\brgba?\(|\bhsla?\(|\boklch\(|\boklab\(|\blab\(|\blch\(|\bcolor-mix\(|\bcolor\(`)
	// URIHexRe erfasst ein prozentkodiertes "#" (%23) gefolgt von Hex-Ziffern.
	URIHexRe     = regexp.MustCompile(`%23[0-9a-fA-F]{3,8}`)
	NamedColorRe = regexp.MustCompile(`(?:^|[^\w-])(` + NamedColors + `)(?:[^\w-]|$)`)
)

// Funde gibt jedes Farbliteral zurück, das in text vorkommt — Hex-Werte,
// prozentkodierte Hex-Werte, Farbfunktionen (als ein gemeinsamer Eintrag,
// da mehrere Funktionsformen auf derselben Stelle zuschlagen können) und
// benannte Farben außer "transparent". Ruft man es zeilenweise auf, bleiben
// Zeilennummern erhalten (siehe css_test.go); ruft man es auf einem ganzen
// SVG-String auf, reicht die bloße Anwesenheit (siehe icons_test.go).
func Funde(text string) []string {
	var funde []string
	for _, m := range HexRe.FindAllString(text, -1) {
		funde = append(funde, m)
	}
	for _, m := range URIHexRe.FindAllString(text, -1) {
		funde = append(funde, m)
	}
	if FuncColorRe.MatchString(text) {
		funde = append(funde, "rgb/rgba/hsl/hsla/oklch/lab/lch/color-mix/color")
	}
	for _, match := range NamedColorRe.FindAllStringSubmatch(text, -1) {
		if len(match) > 1 && match[1] != "transparent" {
			funde = append(funde, match[1])
		}
	}
	return funde
}
