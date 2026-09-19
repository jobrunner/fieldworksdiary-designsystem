package icons

import "testing"

func TestAlleAchtMondphasen(t *testing.T) {
	for _, name := range []string{
		"mond-neu", "mond-zunehmende-sichel", "mond-erstes-viertel",
		"mond-zunehmend-gibbous", "mond-voll", "mond-abnehmend-gibbous",
		"mond-letztes-viertel", "mond-abnehmende-sichel",
	} {
		if _, da := Alle()[name]; !da {
			t.Errorf("Alle() kennt %q nicht", name)
		}
	}
}

func TestJedeMondphaseIstEigenstaendigGezeichnet(t *testing.T) {
	// Der Anlass für dieses Paket: als Emoji waren die Phasen auf vielen
	// Systemen nicht auseinanderzuhalten. Zwei gleich gezeichnete Phasen
	// wären derselbe Mangel in neuer Form.
	phasen := []string{
		"mond-neu", "mond-zunehmende-sichel", "mond-erstes-viertel",
		"mond-zunehmend-gibbous", "mond-voll", "mond-abnehmend-gibbous",
		"mond-letztes-viertel", "mond-abnehmende-sichel",
	}
	gesehen := map[string]string{}
	for _, name := range phasen {
		s := string(Alle()[name])
		if vorher, da := gesehen[s]; da {
			t.Errorf("%s ist identisch gezeichnet wie %s", name, vorher)
		}
		gesehen[s] = name
	}
}

// TestNurZugelasseneSymboleDuerfenFlaechenNutzen prüft die Zusage aus
// darfFlaechenNutzen() (icons.go): eine Fläche (fill="currentColor") darf
// nur bei den dort benannten Ausnahmen auftreten — den acht Mondphasen und
// "standort". Der Test hieß früher TestMondphasenDuerfenFlaechenNutzen, als
// die Mondphasen die einzige Ausnahme waren; er prüft inzwischen eine
// allgemeinere Zusage und trägt deshalb den allgemeineren Namen.
func TestNurZugelasseneSymboleDuerfenFlaechenNutzen(t *testing.T) {
	for name, icon := range Alle() {
		gefuellt := contains(string(icon), `fill="currentColor"`)
		if gefuellt && !darfFlaechenNutzen(name) {
			t.Errorf("%s nutzt Flächen, ist aber nicht in darfFlaechenNutzen() zugelassen — die Regelform ist die Strichzeichnung", name)
		}
	}
}

func TestSonnenAufUndUntergangUnterscheidenSich(t *testing.T) {
	if Alle()["sonnenaufgang"] == Alle()["sonnenuntergang"] {
		t.Error("Sonnenauf- und -untergang sind identisch gezeichnet")
	}
}
