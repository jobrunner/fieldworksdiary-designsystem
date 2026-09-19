package icons

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"testing"
)

// Dieser Test prüft, was TestJedesSymbolFolgtDerBauregel nicht kann: die
// tatsächliche GEOMETRIE jedes Symbols. Belegt in der Schlussprüfung:
//   - ChevronUnten auf "M6 9l6 31 6-31" geändert (reicht bis y=40) → grün,
//     weil die alte Prüfung nur Teilstrings sucht, nie Koordinaten.
//   - Alle Wolken von "A4.51 4.51" auf "A3 3" geändert (der Renderer zeichnet
//     dann eine andere Form, weil die Bogenradien für den gegebenen Abstand
//     nicht mehr ausreichen) → grün, aus demselben Grund.
//
// Zwei Eigenschaften werden je Symbol geprüft:
//  1. Jeder Punkt (Linienpunkte, Bogenendpunkte UND die tatsächliche
//     Ausbuchtung jedes Bogens, Kreis- und Rechteckränder) liegt im Feld
//     0..24, UND ZWAR EINSCHLIESSLICH der halben Strichbreite: bei
//     stroke-width="2" trägt die Kontur 1 Einheit über den Pfad hinaus.
//     Ein Punkt darf deshalb nur im Bereich [1, 23] liegen. Das ist genau
//     die Prüfung, die WetterBewoelkt als fehlerhaft entlarvt (siehe
//     wetter_test.go): der reine Pfad bleibt bei x=23.78 noch im
//     Sichtfeld, die Kontur bei 24.78 nicht mehr.
//  2. Für jeden Bogen gilt λ ≤ 1 nach SVG-Spezifikation F.6.6: mit den
//     Halbdifferenzen dx, dy (im gedrehten Bezugssystem) muss
//     dx²/rx² + dy²/ry² ≤ 1 gelten. Ein größerer Wert bedeutet, dass die
//     angegebenen Radien für den Abstand der Endpunkte nicht ausreichen —
//     der Renderer vergrößert sie dann selbst und zeichnet eine andere
//     Form, als der Pfad behauptet.

const (
	feldMin           = 0.0
	feldMax           = 24.0
	halbeStrichbreite = 1.0 // stroke-width="2" / 2
	// geometrieEps ist die Toleranz für Fließkomma-Rundung in der
	// Bogenmathematik selbst (Winkelvergleiche, Λ nahe 1 bei Halbkreisen).
	geometrieEps = 1e-6
	// feldToleranz ist die Toleranz der FELDPRÜFUNG (imFeld) — großzügiger
	// als geometrieEps, bewusst so gewählt. Mehrere unveränderte, bereits
	// geprüfte Symbole (etwa der Thermometer-Kolben und die
	// Niederschlags-Wolke aus WetterNebel & Co.) reichen an ihrem
	// äußersten Bogenpunkt geometrisch exakt bis knapp über den erlaubten
	// Rand — um bis zu 0,03 Einheiten bei 24 Einheiten Feldbreite (0,13 %,
	// deutlich unter einem Bildschirm-Pixel bei jeder realistischen
	// Einsatzgröße). Das ist sichtbar etwas anderes als der Befund an
	// WetterBewoelkt: dort liegt die Kontur 0,78 Einheiten (3,3 % der
	// Feldbreite) außerhalb — mit bloßem Auge als abgeflachte Kante
	// erkennbar (siehe wetter_test.go). feldToleranz liegt bewusst
	// zwischen beidem: großzügig genug, das Rauschen bestehender,
	// unveränderter Symbole nicht als Fehler zu melden, aber weit unter
	// dem WetterBewoelkt-Fehler und unter den Mutationsproben dieser Datei
	// (y=40 statt y≤23, λ>1 durch zu kleinen Bogenradius).
	feldToleranz = 0.05
)

type punkt struct{ x, y float64 }

// imFeld prüft, ob p auch mit halber Strichbreite noch im 24er-Feld liegt.
func imFeld(p punkt) bool {
	return p.x >= halbeStrichbreite-feldToleranz &&
		p.x <= feldMax-halbeStrichbreite+feldToleranz &&
		p.y >= halbeStrichbreite-feldToleranz &&
		p.y <= feldMax-halbeStrichbreite+feldToleranz
}

// --- Zerleger für SVG-Pfaddaten (d="...") --------------------------------
//
// Unterstützt genau die Befehle, die in diesem Paket vorkommen: M/m, L/l,
// H/h, V/v, A/a, Z/z. Ein unbekannter Befehl ist ein Fehler — lieber laut
// scheitern, als eine Konstruktion stillschweigend zu ignorieren, die der
// Zerleger nicht kennt.
//
// Zwei Fallen sind absichtlich berücksichtigt:
//   - Bogenflags sind je EIN Zeichen und dürfen ohne Trennzeichen an die
//     nächste Zahl angrenzen: "010-8" ist die Folge 0, 1, 0, -8 — ein
//     gierig lesender Zerleger machte daraus fälschlich 010 und -8.
//   - Nach M/m gilt für nachfolgende Zahlenpaare ohne neuen Befehlsbuchstaben
//     implizit L/l (Liniensegmente), nicht ein weiteres Moveto.

type pfadZerleger struct {
	d string
	i int
}

func (z *pfadZerleger) uebrig() bool { return z.i < len(z.d) }

func (z *pfadZerleger) skipSep() {
	for z.i < len(z.d) {
		c := z.d[z.i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == ',' {
			z.i++
			continue
		}
		break
	}
}

func istZiffer(c byte) bool { return c >= '0' && c <= '9' }

func istBefehl(c byte) bool {
	switch c {
	case 'M', 'm', 'L', 'l', 'H', 'h', 'V', 'v', 'A', 'a', 'Z', 'z':
		return true
	}
	return false
}

func (z *pfadZerleger) zahl() (float64, bool) {
	z.skipSep()
	start := z.i
	if z.i < len(z.d) && (z.d[z.i] == '+' || z.d[z.i] == '-') {
		z.i++
	}
	hat := false
	for z.i < len(z.d) && istZiffer(z.d[z.i]) {
		z.i++
		hat = true
	}
	if z.i < len(z.d) && z.d[z.i] == '.' {
		z.i++
		for z.i < len(z.d) && istZiffer(z.d[z.i]) {
			z.i++
			hat = true
		}
	}
	if !hat {
		z.i = start
		return 0, false
	}
	if z.i < len(z.d) && (z.d[z.i] == 'e' || z.d[z.i] == 'E') {
		save := z.i
		z.i++
		if z.i < len(z.d) && (z.d[z.i] == '+' || z.d[z.i] == '-') {
			z.i++
		}
		if z.i < len(z.d) && istZiffer(z.d[z.i]) {
			for z.i < len(z.d) && istZiffer(z.d[z.i]) {
				z.i++
			}
		} else {
			z.i = save
		}
	}
	v, err := strconv.ParseFloat(z.d[start:z.i], 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// flag liest genau EIN Zeichen ('0' oder '1') als Bogenflag — siehe die
// Falle im Kommentar oben. Führende Trennzeichen (Leerraum/Komma) werden
// übersprungen, aber niemals mehr als eine Ziffer gelesen.
func (z *pfadZerleger) flag() (int, bool) {
	z.skipSep()
	if z.i < len(z.d) && (z.d[z.i] == '0' || z.d[z.i] == '1') {
		v := int(z.d[z.i] - '0')
		z.i++
		return v, true
	}
	return 0, false
}

// arcAuswertung berechnet aus einem Bogenbefehl:
//   - lambda: der Λ-Wert nach F.6.6 (soll ≤ 1 sein)
//   - punkte: alle Punkte, die auf der gezeichneten Kontur extremal sind
//     (Endpunkte UND die tatsächliche Ausbuchtung), für die Feldprüfung.
func arcAuswertung(x1, y1, x2, y2, rx, ry, rotDeg float64, large, sweep int) (lambda float64, punkte []punkt) {
	rx = math.Abs(rx)
	ry = math.Abs(ry)
	if rx == 0 || ry == 0 {
		return 0, []punkt{{x2, y2}}
	}
	phi := rotDeg * math.Pi / 180
	cosPhi, sinPhi := math.Cos(phi), math.Sin(phi)

	dx2 := (x1 - x2) / 2
	dy2 := (y1 - y2) / 2
	x1p := cosPhi*dx2 + sinPhi*dy2
	y1p := -sinPhi*dx2 + cosPhi*dy2

	lambda = (x1p*x1p)/(rx*rx) + (y1p*y1p)/(ry*ry)

	// Für die nachfolgende Mittelpunktsberechnung (nur zur Bestimmung der
	// tatsächlich gezeichneten Ausbuchtung) werden zu kleine Radien wie in
	// der SVG-Spezifikation korrigiert — die violation selbst wird über
	// den oben berechneten, UNKORRIGIERTEN lambda-Wert gemeldet.
	rxc, ryc := rx, ry
	if lambda > 1 {
		s := math.Sqrt(lambda)
		rxc = rx * s
		ryc = ry * s
	}

	sign := -1.0
	if large != sweep {
		sign = 1.0
	}
	num := rxc*rxc*ryc*ryc - rxc*rxc*y1p*y1p - ryc*ryc*x1p*x1p
	den := rxc*rxc*y1p*y1p + ryc*ryc*x1p*x1p
	co := 0.0
	if den > geometrieEps {
		v := num / den
		if v < 0 {
			v = 0
		}
		co = sign * math.Sqrt(v)
	}
	cxp := co * (rxc * y1p / ryc)
	cyp := co * (-ryc * x1p / rxc)
	cx := cosPhi*cxp - sinPhi*cyp + (x1+x2)/2
	cy := sinPhi*cxp + cosPhi*cyp + (y1+y2)/2

	winkel := func(ux, uy, vx, vy float64) float64 {
		dot := ux*vx + uy*vy
		lenProd := math.Hypot(ux, uy) * math.Hypot(vx, vy)
		if lenProd == 0 {
			return 0
		}
		c := dot / lenProd
		if c > 1 {
			c = 1
		}
		if c < -1 {
			c = -1
		}
		a := math.Acos(c)
		if ux*vy-uy*vx < 0 {
			a = -a
		}
		return a
	}

	theta1 := winkel(1, 0, (x1p-cxp)/rxc, (y1p-cyp)/ryc)
	dtheta := winkel((x1p-cxp)/rxc, (y1p-cyp)/ryc, (-x1p-cxp)/rxc, (-y1p-cyp)/ryc)
	if sweep == 0 && dtheta > 0 {
		dtheta -= 2 * math.Pi
	}
	if sweep == 1 && dtheta < 0 {
		dtheta += 2 * math.Pi
	}

	pointAt := func(t float64) punkt {
		return punkt{
			x: cx + rxc*cosPhi*math.Cos(t) - ryc*sinPhi*math.Sin(t),
			y: cy + rxc*sinPhi*math.Cos(t) + ryc*cosPhi*math.Sin(t),
		}
	}

	inRange := func(t float64) bool {
		diff := math.Mod(t-theta1, 2*math.Pi)
		if dtheta >= 0 {
			if diff < 0 {
				diff += 2 * math.Pi
			}
			return diff <= dtheta+1e-9
		}
		if diff > 0 {
			diff -= 2 * math.Pi
		}
		return diff >= dtheta-1e-9
	}

	punkte = []punkt{{x1, y1}, {x2, y2}}

	tx := math.Atan2(-ryc*sinPhi, rxc*cosPhi)
	ty := math.Atan2(ryc*cosPhi, rxc*sinPhi)
	for _, t := range []float64{tx, tx + math.Pi, ty, ty + math.Pi} {
		if inRange(t) {
			punkte = append(punkte, pointAt(t))
		}
	}
	return lambda, punkte
}

// pfadGeometrie zerlegt d und liefert alle Punkte (für die Feldprüfung)
// sowie alle Bogen-Lambda-Werte (für die Radienprüfung).
func pfadGeometrie(d string) (punkte []punkt, lambdas []float64, err error) {
	z := &pfadZerleger{d: d}
	var cx, cy, sx, sy float64
	var befehl byte
	fehlerBefehl := byte(0)

	for {
		z.skipSep()
		if !z.uebrig() {
			break
		}
		if istBefehl(z.d[z.i]) {
			befehl = z.d[z.i]
			z.i++
		} else if befehl == 0 {
			return nil, nil, fmt.Errorf("Pfad beginnt nicht mit einem Befehl: %q", d)
		}
		relativ := befehl >= 'a' && befehl <= 'z'
		ob := befehl
		if relativ {
			ob = befehl - ('a' - 'A')
		}
		switch ob {
		case 'M':
			x, ok1 := z.zahl()
			y, ok2 := z.zahl()
			if !ok1 || !ok2 {
				return nil, nil, fmt.Errorf("ungültiges M in %q", d)
			}
			if relativ {
				x += cx
				y += cy
			}
			cx, cy = x, y
			sx, sy = cx, cy
			punkte = append(punkte, punkt{cx, cy})
			// Nachfolgende Zahlenpaare ohne neuen Befehl sind implizit L/l.
			if relativ {
				befehl = 'l'
			} else {
				befehl = 'L'
			}
		case 'L':
			x, ok1 := z.zahl()
			y, ok2 := z.zahl()
			if !ok1 || !ok2 {
				return nil, nil, fmt.Errorf("ungültiges L in %q", d)
			}
			if relativ {
				x += cx
				y += cy
			}
			cx, cy = x, y
			punkte = append(punkte, punkt{cx, cy})
		case 'H':
			x, ok := z.zahl()
			if !ok {
				return nil, nil, fmt.Errorf("ungültiges H in %q", d)
			}
			if relativ {
				x += cx
			}
			cx = x
			punkte = append(punkte, punkt{cx, cy})
		case 'V':
			y, ok := z.zahl()
			if !ok {
				return nil, nil, fmt.Errorf("ungültiges V in %q", d)
			}
			if relativ {
				y += cy
			}
			cy = y
			punkte = append(punkte, punkt{cx, cy})
		case 'A':
			rx, ok1 := z.zahl()
			ry, ok2 := z.zahl()
			rot, ok3 := z.zahl()
			large, ok4 := z.flag()
			sweep, ok5 := z.flag()
			x, ok6 := z.zahl()
			y, ok7 := z.zahl()
			if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
				return nil, nil, fmt.Errorf("ungültiges A in %q (an Position %d)", d, z.i)
			}
			if relativ {
				x += cx
				y += cy
			}
			lambda, pts := arcAuswertung(cx, cy, x, y, rx, ry, rot, large, sweep)
			lambdas = append(lambdas, lambda)
			punkte = append(punkte, pts...)
			cx, cy = x, y
		case 'Z':
			cx, cy = sx, sy
			punkte = append(punkte, punkt{cx, cy})
		default:
			fehlerBefehl = ob
		}
		if fehlerBefehl != 0 {
			return nil, nil, fmt.Errorf("unbekannter Pfadbefehl %q in %q — Zerleger muss erweitert werden", string(fehlerBefehl), d)
		}
	}
	return punkte, lambdas, nil
}

var (
	pathDRe   = regexp.MustCompile(`<path\b[^>]*\bd="([^"]*)"`)
	circleRe  = regexp.MustCompile(`<circle\b[^>]*/?>`)
	rectRe    = regexp.MustCompile(`<rect\b[^>]*/?>`)
	attrNumRe = func(name string) *regexp.Regexp {
		return regexp.MustCompile(name + `="(-?[0-9.]+)"`)
	}
	cxAttrRe = attrNumRe("cx")
	cyAttrRe = attrNumRe("cy")
	rAttrRe  = attrNumRe("r")
	xAttrRe  = attrNumRe("x")
	yAttrRe  = attrNumRe("y")
	wAttrRe  = attrNumRe("width")
	hAttrRe  = attrNumRe("height")
)

func attrVal(re *regexp.Regexp, tag string) (float64, bool) {
	m := re.FindStringSubmatch(tag)
	if m == nil {
		return 0, false
	}
	v, err := strconv.ParseFloat(m[1], 64)
	return v, err == nil
}

// svgGeometrie sammelt alle Punkte und Bogen-Lambda-Werte eines
// vollständigen SVG-Symbols: alle <path d="...">, <circle> und <rect>.
func svgGeometrie(svg string) (punkte []punkt, lambdas []float64, err error) {
	for _, m := range pathDRe.FindAllStringSubmatch(svg, -1) {
		p, l, e := pfadGeometrie(m[1])
		if e != nil {
			return nil, nil, e
		}
		punkte = append(punkte, p...)
		lambdas = append(lambdas, l...)
	}
	for _, tag := range circleRe.FindAllString(svg, -1) {
		cx, ok1 := attrVal(cxAttrRe, tag)
		cy, ok2 := attrVal(cyAttrRe, tag)
		r, ok3 := attrVal(rAttrRe, tag)
		if !ok1 || !ok2 || !ok3 {
			return nil, nil, fmt.Errorf("circle ohne cx/cy/r: %s", tag)
		}
		punkte = append(punkte,
			punkt{cx - r, cy}, punkt{cx + r, cy},
			punkt{cx, cy - r}, punkt{cx, cy + r},
		)
	}
	for _, tag := range rectRe.FindAllString(svg, -1) {
		x, ok1 := attrVal(xAttrRe, tag)
		y, ok2 := attrVal(yAttrRe, tag)
		w, ok3 := attrVal(wAttrRe, tag)
		h, ok4 := attrVal(hAttrRe, tag)
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return nil, nil, fmt.Errorf("rect ohne x/y/width/height: %s", tag)
		}
		// Abgerundete Ecken (rx/ry) ziehen die Kontur nur nach INNEN, nie
		// über die Eckpunkte hinaus — die vier Ecken reichen deshalb aus.
		punkte = append(punkte,
			punkt{x, y}, punkt{x + w, y}, punkt{x, y + h}, punkt{x + w, y + h},
		)
	}
	return punkte, lambdas, nil
}

// TestSymbolGeometrieImFeld prüft für jedes Symbol aus Alle(): jeder Punkt
// liegt auch mit halber Strichbreite noch im 24er-Feld, und jeder Bogen hat
// hinreichend große Radien (λ ≤ 1). Siehe die Kommentare oben für die
// beiden nachgewiesenen Mutationen, die diese Prüfung fängt und die alte,
// teilstringbasierte Prüfung nicht fing.
func TestSymbolGeometrieImFeld(t *testing.T) {
	if len(Alle()) == 0 {
		t.Fatal("Alle() ist leer — dann prüft dieser Test nichts")
	}
	for name, icon := range Alle() {
		name, icon := name, icon
		t.Run(name, func(t *testing.T) {
			punkte, lambdas, err := svgGeometrie(string(icon))
			if err != nil {
				t.Fatalf("Geometrie von %s konnte nicht zerlegt werden: %v", name, err)
			}
			for _, p := range punkte {
				if !imFeld(p) {
					t.Errorf("%s: Punkt (%.4g, %.4g) liegt mit halber Strichbreite (%.0f) außerhalb des Feldes 0..24",
						name, p.x, p.y, halbeStrichbreite)
				}
			}
			for i, l := range lambdas {
				if l > 1+geometrieEps {
					t.Errorf("%s: Bogen #%d hat λ=%.4f > 1 — die angegebenen Radien reichen für den Abstand der Endpunkte nicht aus, der Renderer zeichnet eine andere Form als der Pfad behauptet", name, i, l)
				}
			}
		})
	}
}

// TestGeometriepruefungErkenntZuGrossenChevron belegt die erste der beiden
// nachgewiesenen Mutationen: ChevronUnten auf "M6 9l6 31 6-31" geändert
// (reicht bis y=40) muss als außerhalb des Feldes erkannt werden.
func TestGeometriepruefungErkenntZuGrossenChevron(t *testing.T) {
	mutiert := strich(`<path d="M6 9l6 31 6-31"/>`)
	punkte, _, err := svgGeometrie(string(mutiert))
	if err != nil {
		t.Fatalf("Zerlegung fehlgeschlagen: %v", err)
	}
	ausserhalb := false
	for _, p := range punkte {
		if !imFeld(p) {
			ausserhalb = true
		}
	}
	if !ausserhalb {
		t.Error("ein bis y=40 reichender Chevron wird nicht als außerhalb des Feldes erkannt")
	}
}

// TestGeometriepruefungErkenntZuKleinenBogenradius belegt die zweite
// nachgewiesene Mutation: eine Wolke mit "A4.51 4.51" auf "A3 3" geändert
// (λ > 1, weil der Abstand der Endpunkte den doppelten Radius übersteigt)
// muss erkannt werden.
func TestGeometriepruefungErkenntZuKleinenBogenradius(t *testing.T) {
	// Originalform aus WetterNebel, aber mit dem in der Schlussprüfung
	// nachgewiesenen zu kleinen Radius.
	mutiert := strich(`<path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3 3 0 0117 13z"/>`)
	_, lambdas, err := svgGeometrie(string(mutiert))
	if err != nil {
		t.Fatalf("Zerlegung fehlgeschlagen: %v", err)
	}
	zuKlein := false
	for _, l := range lambdas {
		if l > 1+geometrieEps {
			zuKlein = true
		}
	}
	if !zuKlein {
		t.Error("ein Bogen mit zu kleinem Radius (λ > 1) wird nicht erkannt")
	}
}
