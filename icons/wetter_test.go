package icons

import "testing"

func TestWettersymboleSindVollstaendig(t *testing.T) {
	// Die Namen sind eine Zusage: Tempus bildet seine WMO-Codes darauf ab.
	// Die Zuordnung selbst bleibt in Tempus — hier liegt nur die Zeichnung.
	for _, name := range []string{
		"wetter-klar-tag", "wetter-klar-nacht",
		"wetter-leicht-bewoelkt-tag", "wetter-leicht-bewoelkt-nacht",
		"wetter-bewoelkt", "wetter-nebel", "wetter-regen", "wetter-schauer",
		"wetter-schnee", "wetter-schneeschauer", "wetter-gewitter", "wetter-hagel",
	} {
		if _, da := Alle()[name]; !da {
			t.Errorf("Alle() kennt %q nicht", name)
		}
	}
}

func TestTagUndNachtUnterscheidenSich(t *testing.T) {
	// Sonne und Mond müssen verschieden gezeichnet sein — sonst trägt die
	// Unterscheidung zwischen Tag und Nacht keine Information.
	paare := [][2]string{
		{"wetter-klar-tag", "wetter-klar-nacht"},
		{"wetter-leicht-bewoelkt-tag", "wetter-leicht-bewoelkt-nacht"},
	}
	for _, p := range paare {
		if Alle()[p[0]] == Alle()[p[1]] {
			t.Errorf("%s und %s sind identisch gezeichnet", p[0], p[1])
		}
	}
}
