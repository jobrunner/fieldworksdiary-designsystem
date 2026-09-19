package icons

// Bedien-Symbole: was in jedem Dienst vorkommt, unabhängig vom Fach.

// Standort — für den Knopf „Aktuellen Standort verwenden", der in Expertus
// prominent vorkommt.
//
// Radius 7 statt der ursprünglichen 3: in der Symbolgalerie wirkte das
// Symbol neben Information() und Fehler() (beide r=10) deutlich leichter —
// ein kleiner Kreis auf demselben 24er-Raster nimmt weniger optisches
// Gewicht ein als sein Nachbar. r=7 gleicht das an, ohne die vier Striche an
// den Feldrand zu drängen.
//
// Der Mittelpunkt ist bewusst eine gefüllte Fläche (fill="currentColor"
// statt eines weiteren Kreis-Strichs): das Bild ist ein Fadenkreuz für die
// AKTUELLE Position, keine Stecknadel für EINEN Ort auf der Karte — dieser
// Unterschied ist in einer Feld-Anwendung wichtig. Ein Punkt als Fläche
// bringt genau die optische Betonung, die eine Positionsmarkierung
// braucht; ein hohler Punkt sähe wie ein weiterer Ring aus. Damit nutzt
// dieses Symbol als einziges Nicht-Mond-Symbol eine Fläche — siehe
// darfFlaechenNutzen() in icons.go für die ausdrücklich zugelassene
// Ausnahme.
//
// Übernommen aus Ortus, wo dieses Symbol bereits in Gebrauch ist.
func Standort() Icon {
	return strich(`<circle cx="12" cy="12" r="7"/><circle cx="12" cy="12" r="2" fill="currentColor" stroke="none"/><path d="M12 1v3m0 16v3M1 12h3m16 0h3"/>`)
}

// ChevronUnten — Aufklappen. In Ortus steht derselbe Pfad dreimal in einer
// Datei; das ist der Grund, warum es dieses Paket gibt.
func ChevronUnten() Icon {
	return strich(`<path d="M6 9l6 6 6-6"/>`)
}

// Kalender — Datums- und Zeitauswahl. Übernommen aus Tempus.
func Kalender() Icon {
	return strich(`<rect x="3" y="4" width="18" height="18" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/>`)
}

func Schliessen() Icon {
	return strich(`<path d="M18 6L6 18M6 6l12 12"/>`)
}

func Suche() Icon {
	return strich(`<circle cx="11" cy="11" r="8"/><path d="M21 21l-4.35-4.35"/>`)
}

func Herunterladen() Icon {
	return strich(`<path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3"/>`)
}

func Kopieren() Icon {
	return strich(`<rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>`)
}

func Haken() Icon {
	return strich(`<path d="M20 6L9 17l-5-5"/>`)
}

func Warnung() Icon {
	return strich(`<path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><path d="M12 9v4M12 17h.01"/>`)
}

func Information() Icon {
	return strich(`<circle cx="12" cy="12" r="10"/><path d="M12 16v-4M12 8h.01"/>`)
}

func Fehler() Icon {
	return strich(`<circle cx="12" cy="12" r="10"/><path d="M15 9l-6 6M9 9l6 6"/>`)
}

func Menue() Icon {
	return strich(`<path d="M3 12h18M3 6h18M3 18h18"/>`)
}

// bediensymbole ist die Sammlung dieser Datei. Alle() setzt daraus und aus
// den übrigen Dateien die Gesamtmenge zusammen; so muss niemand beim
// Hinzufügen eines Symbols an eine zweite Stelle denken.
func bediensymbole() map[string]Icon {
	return map[string]Icon{
		"standort":      Standort(),
		"chevron-unten": ChevronUnten(),
		"kalender":      Kalender(),
		"schliessen":    Schliessen(),
		"suche":         Suche(),
		"herunterladen": Herunterladen(),
		"kopieren":      Kopieren(),
		"haken":         Haken(),
		"warnung":       Warnung(),
		"information":   Information(),
		"fehler":        Fehler(),
		"menue":         Menue(),
	}
}

// Alle liefert jedes Symbol des Pakets unter seinem Namen. Prüfung und
// Referenzseite leiten sich daraus ab, damit ein neues Symbol nicht still
// ungeprüft und ungezeigt bleibt.
func Alle() map[string]Icon {
	out := map[string]Icon{}
	for _, teil := range []map[string]Icon{bediensymbole(), wettersymbole(), himmelssymbole(), messwertsymbole()} {
		for k, v := range teil {
			out[k] = v
		}
	}
	return out
}
