// Package icons liefert die Symbole der fieldworksdiary-Dienste als fertiges
// SVG. Je Symbol eine Funktion; die gemeinsame Hülle steht an genau einer
// Stelle, damit die Bauregel nicht 38-mal wiederholt und dabei verletzt wird.
//
// Das Modul liefert Zeichnungen, keine Fachlogik: welcher Zustand welches
// Symbol bekommt — etwa welcher WMO-Wettercode Nebel bedeutet — entscheidet
// der Dienst.
package icons

import (
	"html"
	"html/template"
	"strings"
)

// Icon ist fertiges SVG-Markup. Eigener Typ statt string, damit an der
// Aufrufstelle erkennbar bleibt, dass hier Markup eingesetzt wird und keine
// Beschriftung.
//
// BEWUSST kein Alias auf template.HTML (type Icon = template.HTML): das
// böte html/template zwar von sich aus Schutz vor Maskierung, würde dieses
// abhängigkeitsfreie Zeichnungspaket aber dauerhaft an html/template binden
// — und einen Bruch für jeden Aufrufer bedeuten, der Icon als eigenen,
// methodentragenden Typ nutzt (icon.String(), Vergleich in Tests). Stattdessen
// liefert HTML() den sicheren Übergang an der Stelle, an der er gebraucht
// wird: ein Dienst, der {{.Symbol}} statt {{.Symbol.HTML}} in eine
// html/template-Vorlage schreibt, bekommt sichtbaren SVG-Quelltext statt
// eines Symbols — Icon ist für html/template gewöhnlicher Text, keine
// Auszeichnung. Kopfzeile()/Fusszeile() im Wurzelpaket liefern deshalb
// bereits template.HTML; icons.Icon.HTML() zieht mit derselben Zusage nach.
// Siehe README.md, Abschnitt "Symbole", für den Hinweis an Dienste.
type Icon string

func (i Icon) String() string { return string(i) }

// HTML liefert das Symbol als template.HTML, damit html/template es als
// Markup einsetzt statt es zu maskieren. Aufrufstelle in einer Vorlage:
// {{ .Symbol.HTML }} statt des gefährlichen {{ .Symbol }}.
func (i Icon) HTML() template.HTML { return template.HTML(i) }

// MitKlasse setzt eine CSS-Klasse auf das <svg>-Element eines Symbols, ohne
// dass ein Dienst dafür Zeichenketten-Chirurgie am Markup selbst betreiben
// muss. Vorher musste jeder Dienst etwas wie
// strings.Replace(svg, "<svg ", `<svg class="icon" `, 1) nachbauen — die
// Referenzseite (demo/demo.go) tat das bereits, jeder weitere Dienst hätte
// es wiederholt.
func MitKlasse(icon Icon, klasse string) Icon {
	return Icon(strings.Replace(string(icon), "<svg ", `<svg class="`+klasse+`" `, 1))
}

// MitBeschriftung macht ein Symbol ansagbar, das ALLEIN steht — ohne
// begleitenden Text, etwa in einem Knopf ohne Beschriftung. huelle() setzt
// aria-hidden="true", weil ein Symbol NEBEN einer Beschriftung Schmuck ist
// und sonst doppelt vorgelesen würde (siehe icons.go); steht es allein,
// trägt es die Beschriftung selbst und darf nicht länger vor Screenreadern
// versteckt sein. MitBeschriftung ersetzt deshalb aria-hidden="true" durch
// role="img" und aria-label mit dem übergebenen Text (maskiert für den
// Attributkontext).
func MitBeschriftung(icon Icon, beschriftung string) Icon {
	ersetzt := `role="img" aria-label="` + html.EscapeString(beschriftung) + `"`
	return Icon(strings.Replace(string(icon), `aria-hidden="true"`, ersetzt, 1))
}

// huelle setzt ein Symbol aus seinem Inhalt zusammen.
//
// currentColor ist der Kern der Bauregel: ein Symbol nimmt damit die Farbe
// seines Umfelds an und folgt Thema und Zustand, ohne eine eigene Farbe zu
// kennen. Ein Symbol mit eigener Farbe stünde außerhalb der Kontrastprüfung
// des Moduls.
//
// Ohne width und height: die Größe bestimmt der Einsatzort über CSS. Feste
// Maße würden ein Symbol im Knopf und in der Tabelle gleich groß machen.
//
// aria-hidden, weil ein Symbol neben einer Beschriftung Schmuck ist und
// sonst doppelt vorgelesen würde. Steht es allein — etwa in einem Knopf ohne
// Text —, trägt das umgebende Bedienelement die Beschriftung.
//
// gefuellt kehrt die Zeichenart um: Flächen statt Striche. Das ist allein
// für die Mondphasen vorgesehen, bei denen die Aufteilung zwischen
// beleuchtetem und unbeleuchtetem Teil die Aussage trägt — eine reine
// Strichzeichnung könnte zunehmenden und abnehmenden Halbmond nicht
// unterscheiden.
func huelle(inhalt string, gefuellt bool) Icon {
	var b strings.Builder
	b.WriteString(`<svg viewBox="0 0 24 24" `)
	if gefuellt {
		b.WriteString(`fill="currentColor" stroke="none" `)
	} else {
		b.WriteString(`fill="none" stroke="currentColor" stroke-width="2" `)
		b.WriteString(`stroke-linecap="round" stroke-linejoin="round" `)
	}
	b.WriteString(`aria-hidden="true" focusable="false">`)
	b.WriteString(inhalt)
	b.WriteString(`</svg>`)
	return Icon(b.String())
}

// strich baut ein Symbol in der Regelform: Strichzeichnung.
func strich(inhalt string) Icon { return huelle(inhalt, false) }

// flaeche baut ein Symbol als Flächenzeichnung. Nur für die Mondphasen —
// siehe die Begründung an huelle.
func flaeche(inhalt string) Icon { return huelle(inhalt, true) }

// darfFlaechenNutzen zählt ausdrücklich auf, welche Symbole eine Fläche
// (fill="currentColor") tragen dürfen, statt sich auf die Strichzeichnung
// zu beschränken. Diese Liste steht bewusst hier im Produktivcode statt nur
// im Test: die Zusage soll dort sichtbar sein, wo Symbole entstehen, und
// nicht nur dort, wo sie geprüft wird — wer ein neues Symbol schreibt, das
// eine Fläche braucht, findet die Stelle, die dafür freigeschaltet werden
// muss, statt sie erst im Test zu erraten.
//
// Zwei Gruppen sind zugelassen:
//   - die acht Mondphasen, deren Aufteilung zwischen beleuchtetem und
//     unbeleuchtetem Teil die Aussage trägt (siehe huelle());
//   - "standort", dessen gefüllter Mittelpunkt das Symbol als Fadenkreuz
//     für die AKTUELLE Position lesbar macht statt als hohle Stecknadel
//     (siehe Standort() in bedienung.go).
//
// Jedes weitere Symbol, das ohne Absicht eine Fläche nutzt, soll weiterhin
// auffallen — deshalb eine benannte Liste statt einer aufgeweichten
// Bedingung.
func darfFlaechenNutzen(name string) bool {
	if name == "standort" {
		return true
	}
	return len(name) >= 5 && name[:5] == "mond-"
}
