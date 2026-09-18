package icons

// Wettersymbole. Das Paket liefert die Zeichnungen; welcher WMO-Code welches
// Symbol bekommt, entscheidet der Dienst — diese Zuordnung ist Wetterlogik
// und ändert sich mit der Schnittstelle des Datenanbieters.

func WetterKlarTag() Icon {
	return strich(`<circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>`)
}

func WetterKlarNacht() Icon {
	return strich(`<path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>`)
}

func WetterLeichtBewoelktTag() Icon {
	// Wolke sitzt höher als in Niederschlagssymbolen (y=12-20 statt y=5-13), um Platz
	// für die Sonne oben zu haben. Die letzte Arc nutzt Radius 5 statt 3.5, weil die
	// Distanz von (16.5,11) zu (17,20) etwa 9 ist — größer als 2*3.5=7.
	return strich(`<path d="M12 4V2M5.6 5.6L4.2 4.2M4 12H2M18.4 5.6l1.4-1.4"/><path d="M8.5 12a3.5 3.5 0 015.9-2.5"/><path d="M17 20H7a4 4 0 010-8 5 5 0 019.5-1A5 5 0 0117 20z"/>`)
}

func WetterLeichtBewoelktNacht() Icon {
	// Mondsichel: zwei konzentrische Bögen (wie WetterKlarNacht, aber verkleinert
	// und versetzt). Großer Bogen (Radius 4) von (10,6) zu (7,2), dann kleiner Bogen
	// (Radius 3) zurück zu (10,6). Das Fehlen eines z am Ende des ersten Bogens ist
	// beabsichtigt — die beiden Bögen sind nahtlos verbunden.
	// Wolke sitzt höher als in Niederschlagssymbolen (y=12-20 statt y=5-13), um Platz
	// für den Mond oben zu haben. Die letzte Arc nutzt Radius 5 statt 3.5, weil die
	// Distanz von (16.5,11) zu (17,20) etwa 9 ist — größer als 2*3.5=7.
	return strich(`<path d="M10 6A4 4 0 1107 2A3 3 0 0010 6z"/><path d="M17 20H7a4 4 0 010-8 5 5 0 019.5-1A5 5 0 0117 20z"/>`)
}

// WetterBewoelkt nutzt eine größere Wolke (a5 5 statt a4 4), weil dieses Symbol
// nur die Wolke zeigt und das Feld ausfüllen kann. Alle übrigen Symbole mit Wolke
// brauchen Platz für Niederschlag darunter (Regen, Schnee, Hagel) oder Details
// (Nebel, Gewitter), weshalb dort die kleinere Wolkenform verwendet wird.
func WetterBewoelkt() Icon {
	return strich(`<path d="M18 17H7a5 5 0 010-10 6 6 0 0111.5-1.5A4.5 4.5 0 0118 17z"/>`)
}

func WetterNebel() Icon {
	return strich(`<path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M4 17h16M7 21h13"/>`)
}

func WetterRegen() Icon {
	return strich(`<path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M8 17v3M12 17v4M16 17v3"/>`)
}

func WetterSchauer() Icon {
	return strich(`<path d="M12 5V3M6.5 6.5L5 5M20 13h-1"/><path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M9 17v3M15 17v3"/>`)
}

func WetterSchnee() Icon {
	return strich(`<path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M8 18h.01M12 20h.01M16 18h.01M10 21h.01M14 17h.01"/>`)
}

func WetterSchneeschauer() Icon {
	return strich(`<path d="M12 5V3M6.5 6.5L5 5"/><path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M9 18h.01M13 20h.01M16 18h.01"/>`)
}

func WetterGewitter() Icon {
	// Blitz startet bei y=16, Wolke endet bei y=13. Der Blitz soll unterhalb der Wolke
	// sichtbar sein, aber innerhalb des viewBox (0-24) bleiben. Endpunkt daher bei y=23
	// statt y=25.
	return strich(`<path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M13 16l-3 5h4l-3 2"/>`)
}

func WetterHagel() Icon {
	return strich(`<path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M8 17v1M12 17v1M16 17v1M10 20v1M14 20v1"/>`)
}

func wettersymbole() map[string]Icon {
	return map[string]Icon{
		"wetter-klar-tag":              WetterKlarTag(),
		"wetter-klar-nacht":            WetterKlarNacht(),
		"wetter-leicht-bewoelkt-tag":   WetterLeichtBewoelktTag(),
		"wetter-leicht-bewoelkt-nacht": WetterLeichtBewoelktNacht(),
		"wetter-bewoelkt":              WetterBewoelkt(),
		"wetter-nebel":                 WetterNebel(),
		"wetter-regen":                 WetterRegen(),
		"wetter-schauer":               WetterSchauer(),
		"wetter-schnee":                WetterSchnee(),
		"wetter-schneeschauer":         WetterSchneeschauer(),
		"wetter-gewitter":              WetterGewitter(),
		"wetter-hagel":                 WetterHagel(),
	}
}
