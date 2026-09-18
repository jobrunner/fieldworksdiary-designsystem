package icons

// Symbole für Messgrößen — die Zahlen, die Tempus neben Wetter- und
// Himmelssymbole stellt (Taupunkt, Temperatur, Wind, Luftdruck) sowie die
// großräumigen Angaben in Ortus (Globus) und zusammengefasste Reihen
// (Diagramm).

// Tropfen — Taupunkt und Feuchte. Ersetzt das Emoji in Tempus.
func Tropfen() Icon {
	return strich(`<path d="M12 2.7l5.3 5.3a7.5 7.5 0 11-10.6 0z"/>`)
}

// Globus — großräumige Angaben wie Bioklima.
func Globus() Icon {
	return strich(`<circle cx="12" cy="12" r="10"/><path d="M2 12h20M12 2a15 15 0 010 20 15 15 0 010-20z"/>`)
}

// Diagramm — zusammengefasste Werte über einen Zeitraum.
func Diagramm() Icon {
	return strich(`<path d="M3 3v18h18"/><path d="M7 15l4-4 3 3 5-6"/>`)
}

// Thermometer — Lufttemperatur.
func Thermometer() Icon {
	return strich(`<path d="M14 14.76V3.5a2.5 2.5 0 00-5 0v11.26a4.5 4.5 0 105 0z"/>`)
}

// Wind — Windgeschwindigkeit und -richtung.
func Wind() Icon {
	return strich(`<path d="M9.6 4.6A2 2 0 1111 8H2M12.6 19.4A2 2 0 1014 16H2M17.7 7.7A2.5 2.5 0 1119.5 12H2"/>`)
}

// Druck — Luftdruck.
func Druck() Icon {
	return strich(`<circle cx="12" cy="12" r="9"/><path d="M12 12l4-3M12 7v1"/>`)
}

// messwertsymbole sammelt die Symbole dieser Datei für Alle().
func messwertsymbole() map[string]Icon {
	return map[string]Icon{
		"tropfen":     Tropfen(),
		"globus":      Globus(),
		"diagramm":    Diagramm(),
		"thermometer": Thermometer(),
		"wind":        Wind(),
		"druck":       Druck(),
	}
}
