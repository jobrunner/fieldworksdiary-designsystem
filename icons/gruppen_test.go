package icons

import "testing"

// TestGruppenEnthaeltJedesSymbolGenauEinmal leitet die erwartete Menge aus
// Alle() ab statt aus einer Aufzählung im Testcode — sonst würde ein neues
// Symbol still weder als fehlend noch als doppelt gemeldet.
func TestGruppenEnthaeltJedesSymbolGenauEinmal(t *testing.T) {
	alle := Alle()
	gesehen := map[string]int{}

	for _, gruppe := range Gruppen() {
		if gruppe.Titel == "" {
			t.Error("eine Gruppe hat keinen Titel")
		}
		for _, name := range gruppe.Namen {
			gesehen[name]++
			if _, da := alle[name]; !da {
				t.Errorf("Gruppe %q nennt %q, das Alle() nicht kennt", gruppe.Titel, name)
			}
		}
	}

	for name := range alle {
		switch gesehen[name] {
		case 0:
			t.Errorf("Symbol %q aus Alle() steht in keiner Gruppe", name)
		case 1:
			// genau richtig
		default:
			t.Errorf("Symbol %q steht in %d Gruppen statt in genau einer", name, gesehen[name])
		}
	}
}

// TestMondphasenGruppeHatDieNatuerlicheAbfolge belegt die Reihenfolge, um
// die es bei den Mondphasen eigentlich geht: alphabetisch stünde
// "mond-abnehmend-gibbous" vor "mond-neu" — die Abfolge der Beleuchtung
// wäre zerstört.
func TestMondphasenGruppeHatDieNatuerlicheAbfolge(t *testing.T) {
	erwartet := []string{
		"mond-neu", "mond-zunehmende-sichel", "mond-erstes-viertel",
		"mond-zunehmend-gibbous", "mond-voll", "mond-abnehmend-gibbous",
		"mond-letztes-viertel", "mond-abnehmende-sichel",
	}

	var mondphasen []string
	for _, gruppe := range Gruppen() {
		if gruppe.Titel == "Mondphasen" {
			mondphasen = gruppe.Namen
		}
	}

	if len(mondphasen) != len(erwartet) {
		t.Fatalf("Gruppe %q hat %d Symbole, erwartet %d", "Mondphasen", len(mondphasen), len(erwartet))
	}
	for i, name := range erwartet {
		if mondphasen[i] != name {
			t.Errorf("Position %d: %q, erwartet %q — die Abfolge ist nicht die natürliche", i, mondphasen[i], name)
		}
	}
}
