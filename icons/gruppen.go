package icons

// Gruppe fasst zusammengehörige Symbole in ihrer für Menschen sinnvollen
// Reihenfolge zusammen. Eine alphabetische Sortierung würde ausgerechnet bei
// den Mondphasen die Abfolge zerreißen, die dort die eigentliche Aussage
// trägt — neu, zunehmende Sichel, erstes Viertel, zunehmend gibbous, voll,
// abnehmend gibbous, letztes Viertel, abnehmende Sichel.
type Gruppe struct {
	Titel string
	Namen []string
}

// Gruppen liefert alle Symbole aus Alle(), geordnet nach Bedeutungszusammenhang
// statt alphabetisch. Die Zuordnung steht hier statt nur auf der
// Referenzseite, damit ein künftiger Dienst — etwa eine Symbolauswahl in
// Tempus oder Ortus — dieselbe Reihenfolge nutzen kann, ohne sie erneut zu
// erfinden. gruppen_test.go stellt sicher, dass jedes Symbol aus Alle() in
// genau einer Gruppe steht.
func Gruppen() []Gruppe {
	return []Gruppe{
		{
			Titel: "Bedienung",
			Namen: []string{
				"standort", "chevron-unten", "kalender", "schliessen", "suche",
				"herunterladen", "kopieren", "haken", "warnung", "information",
				"fehler", "menue",
			},
		},
		{
			Titel: "Wetter",
			Namen: []string{
				"wetter-klar-tag", "wetter-klar-nacht",
				"wetter-leicht-bewoelkt-tag", "wetter-leicht-bewoelkt-nacht",
				"wetter-bewoelkt", "wetter-nebel", "wetter-regen",
				"wetter-schauer", "wetter-schnee", "wetter-schneeschauer",
				"wetter-gewitter", "wetter-hagel",
			},
		},
		{
			// Die natürliche Abfolge der Mondphasen — nicht alphabetisch,
			// siehe Begründung an Gruppe.
			Titel: "Mondphasen",
			Namen: []string{
				"mond-neu", "mond-zunehmende-sichel", "mond-erstes-viertel",
				"mond-zunehmend-gibbous", "mond-voll", "mond-abnehmend-gibbous",
				"mond-letztes-viertel", "mond-abnehmende-sichel",
			},
		},
		{
			Titel: "Sonnenstände",
			Namen: []string{"sonnenaufgang", "sonnenuntergang"},
		},
		{
			Titel: "Messwerte",
			Namen: []string{
				"tropfen", "globus", "diagramm", "thermometer", "wind", "druck",
			},
		},
	}
}
