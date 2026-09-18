// Package icons liefert die Symbole der fieldworksdiary-Dienste als fertiges
// SVG. Je Symbol eine Funktion; die gemeinsame Hülle steht an genau einer
// Stelle, damit die Bauregel nicht 38-mal wiederholt und dabei verletzt wird.
//
// Das Modul liefert Zeichnungen, keine Fachlogik: welcher Zustand welches
// Symbol bekommt — etwa welcher WMO-Wettercode Nebel bedeutet — entscheidet
// der Dienst.
package icons

import "strings"

// Icon ist fertiges SVG-Markup. Eigener Typ statt string, damit an der
// Aufrufstelle erkennbar bleibt, dass hier Markup eingesetzt wird und keine
// Beschriftung.
type Icon string

func (i Icon) String() string { return string(i) }

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
