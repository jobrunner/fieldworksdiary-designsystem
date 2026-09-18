package icons

// Mondphasen und Sonnenstände.
//
// Die Mondphasen sind die einzige Ausnahme von der Strichzeichnung: bei ihnen
// trägt die Aufteilung zwischen beleuchtetem und unbeleuchtetem Teil die
// Aussage. Eine reine Strichzeichnung könnte zunehmenden und abnehmenden
// Halbmond nicht unterscheiden — und genau diese Unterscheidbarkeit war der
// Anlass, die vorherigen Emoji zu ersetzen.
//
// Konstruktion, abweichend vom Ausgangsentwurf im Plan: Alle acht Phasen
// nutzen denselben Vollkreis (Mittelpunkt 12,12, Radius 9) als Grundfläche.
// Die Trennlinie zwischen hell und dunkel ist eine Ellipsen-Arc mit
// ry=9 (senkrechter Radius, gleich dem Kreisradius) und variablem rx
// (waagerechter Radius, bestimmt die Wölbung):
//
//   - rx=0 (eine gerade Linie statt Arc): die Trennlinie ist der Durchmesser
//     selbst — das ist bei den Vierteln richtig, denn dort steht der
//     Betrachter exakt seitlich zur Licht-Schatten-Grenze der Kugel.
//   - rx=4,5 bzw. rx=6: die Trennlinie wölbt sich — bei Sichel und Gibbous
//     ist sie nie eine Gerade, weil man die gekrümmte Kugelfläche schräg
//     sieht. Der Radius bestimmt dabei den beleuchteten Flächenanteil (siehe
//     Fix-Runde 1 unten).
//   - rx=9: die Trennlinie fällt auf den Kreisumfang — Voll- bzw. Neumond.
//
// Flächenanteil einer Sichel/Gibbous-Form in Abhängigkeit von rx=r (bei
// festem ry=9, Kreisradius 9): die Fläche zwischen dem äußeren Halbkreis
// und der inneren Ellipsen-Arc ist eine Integration von
// (9-r)·sqrt(1-((y-12)/9)²) über y von 3 bis 21, das ergibt 9(9-r)(π/2).
// Geteilt durch die Kreisfläche 81π bleibt der einfache Bruch (9-r)/18 für
// die Sichel bzw. 0,5+r/18 für Gibbous (siehe Fix-Runde 1 im Bericht für
// die Herleitung und die Kontrolle per Rasterung).
//
// Der Plan hatte für Sichel und Gibbous zwei Radien "6 6" verwendet
// (rx=ry=6). Das unterschreitet den nötigen Radius: Start- und Endpunkt der
// Arc (12,3) und (12,21) liegen 18 Einheiten auseinander, mehr als der
// doppelte angegebene Radius (12). SVG vergrößert den Radius dann still, bis
// er reicht — gezeichnet wird eine andere Ellipse als angegeben. Mit ry=9
// ist die senkrechte Ausdehnung exakt gleich dem halben Abstand (9), der
// Radius passt exakt und wird nicht nachkorrigiert; rx bleibt frei wählbar
// und bestimmt allein die Wölbungstiefe. Diese Konstruktion ist von
// icons/wetter.go übernommen (siehe die Kommentare zu WetterNebel).
//
// Wohin eine Arc sich wölbt, folgt der Reiserichtung: "sweep=1" ist
// Uhrzeigersinn. Von oben (12,3) nach unten (12,21) wölbt sweep=1 nach
// rechts, sweep=0 nach links; von unten nach oben ist es umgekehrt. Ein
// zunehmender (rechts beleuchteter) Halbmond durchläuft daher stets
// "A9 9 0 0 1 12 21" (Außenrand rechts) und ein abnehmender stets
// "A9 9 0 0 0 12 21" (Außenrand links); die zweite Arc jeder Sichel- oder
// Gibbous-Form schließt mit passendem sweep zurück zum Ausgangspunkt.
//
// Zusätzlich zeigt jede Phase außer Voll- und Neumond den Kreisumriss als
// Strich (fill="none" stroke="currentColor" innerhalb der gefüllten Hülle —
// das enthält keine Farbe im Klartext und ist deshalb zulässig), damit auch
// der unbeleuchtete Teil als Teil derselben Mondscheibe erkennbar bleibt und
// alle acht Symbole in einer Reihe gleich groß wirken.

// MondNeu — unbeleuchtet: nur der Umriss. Braucht keine Fläche, deshalb
// strich() statt flaeche().
func MondNeu() Icon {
	return strich(`<circle cx="12" cy="12" r="9"/>`)
}

// MondZunehmendeSichel — beleuchteter Streifen am rechten Rand: die Fläche
// zwischen dem äußeren rechten Halbkreis (rx=9) und einer rechten Wölbung
// weiter innen.
//
// Fix-Runde 1: rx=6 (≈17 % Flächenanteil) war bei 24 px kaum von Neumond zu
// unterscheiden und ließ die Seitenlage nicht erkennen. Mit rx=4,5 ergibt
// die Formel (9-r)/18 = 4,5/18 = 25 % — deutlich mehr Fläche, aber klar
// unter dem Viertel (50 %), wie gefordert.
func MondZunehmendeSichel() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 3 A9 9 0 0 1 12 21 A4.5 9 0 0 0 12 3 Z"/>`)
}

// MondErstesViertel — rechte Hälfte beleuchtet, Trennlinie ist der
// Durchmesser (rx=0, also eine gerade Linie statt Arc).
func MondErstesViertel() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 3 A9 9 0 0 1 12 21 Z"/>`)
}

// MondZunehmendGibbous — mehr als die Hälfte, rechts: der rechte Halbkreis
// plus eine zusätzliche Wölbung, die über den Durchmesser hinaus in die
// linke (dunkle) Seite reicht.
func MondZunehmendGibbous() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 3 A9 9 0 0 1 12 21 A6 9 0 0 1 12 3 Z"/>`)
}

// MondVoll — vollständig beleuchtet: die Fläche füllt den ganzen Kreis, ein
// zusätzlicher Umriss wäre hier ohne Wirkung.
func MondVoll() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9"/>`)
}

// MondAbnehmendGibbous — Spiegelbild von MondZunehmendGibbous: mehr als die
// Hälfte, links.
func MondAbnehmendGibbous() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 3 A9 9 0 0 0 12 21 A6 9 0 0 0 12 3 Z"/>`)
}

// MondLetztesViertel — Spiegelbild von MondErstesViertel: linke Hälfte
// beleuchtet, gerade Trennlinie.
func MondLetztesViertel() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 3 A9 9 0 0 0 12 21 Z"/>`)
}

// MondAbnehmendeSichel — Spiegelbild von MondZunehmendeSichel: beleuchteter
// Streifen am linken Rand, gleiches rx=4,5 wie dort (siehe Fix-Runde 1).
func MondAbnehmendeSichel() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 3 A9 9 0 0 0 12 21 A4.5 9 0 0 1 12 3 Z"/>`)
}

// SonnenAufgang — Horizont, Sonnenhalbkreis, vier Strahlen und ein nach oben
// weisender Pfeil. Der Pfeil ersetzt den fünften (oberen) Strahl, damit an
// dieser Stelle nicht zwei Linien übereinanderliegen.
func SonnenAufgang() Icon {
	return strich(`<path d="M17 18a5 5 0 00-10 0"/><path d="M4.22 10.22l1.42 1.42M1 18h2M21 18h2M18.36 11.64l1.42-1.42M23 22H1"/><path d="M8 6l4-4 4 4"/>`)
}

// SonnenUntergang — dieselbe Grundform wie SonnenAufgang, einziger
// Unterschied ist der Pfeil: er weist nach unten statt nach oben. Das ist
// bewusst die einzige Abweichung — Auf- und Untergang sollen sich auf den
// ersten Blick nur durch die Richtung unterscheiden, nicht durch die Szene.
func SonnenUntergang() Icon {
	return strich(`<path d="M17 18a5 5 0 00-10 0"/><path d="M4.22 10.22l1.42 1.42M1 18h2M21 18h2M18.36 11.64l1.42-1.42M23 22H1"/><path d="M8 4l4 4 4-4"/>`)
}

// himmelssymbole ist die Sammlung dieser Datei. Alle() setzt daraus und aus
// den übrigen Dateien die Gesamtmenge zusammen.
func himmelssymbole() map[string]Icon {
	return map[string]Icon{
		"mond-neu":               MondNeu(),
		"mond-zunehmende-sichel": MondZunehmendeSichel(),
		"mond-erstes-viertel":    MondErstesViertel(),
		"mond-zunehmend-gibbous": MondZunehmendGibbous(),
		"mond-voll":              MondVoll(),
		"mond-abnehmend-gibbous": MondAbnehmendGibbous(),
		"mond-letztes-viertel":   MondLetztesViertel(),
		"mond-abnehmende-sichel": MondAbnehmendeSichel(),
		"sonnenaufgang":          SonnenAufgang(),
		"sonnenuntergang":        SonnenUntergang(),
	}
}
