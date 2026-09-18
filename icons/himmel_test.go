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

func TestMondphasenDuerfenFlaechenNutzen(t *testing.T) {
	// Die einzige zugelassene Ausnahme von der Strichzeichnung. Geprüft
	// wird, dass sie auch wirklich nur hier auftritt.
	for name, icon := range Alle() {
		gefuellt := contains(string(icon), `fill="currentColor"`)
		istMond := len(name) >= 5 && name[:5] == "mond-"
		if gefuellt && !istMond {
			t.Errorf("%s nutzt Flächen, ist aber keine Mondphase — die Regelform ist die Strichzeichnung", name)
		}
	}
}

func TestSonnenAufUndUntergangUnterscheidenSich(t *testing.T) {
	if Alle()["sonnenaufgang"] == Alle()["sonnenuntergang"] {
		t.Error("Sonnenauf- und -untergang sind identisch gezeichnet")
	}
}
