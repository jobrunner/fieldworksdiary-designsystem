package icons

import "testing"

func TestHuelleFolgtDerBauregel(t *testing.T) {
	got := string(huelle(`<path d="M6 9l6 6 6-6"/>`, false))

	for _, pflicht := range []string{
		`viewBox="0 0 24 24"`,
		`fill="none"`,
		`stroke="currentColor"`,
		`stroke-width="2"`,
		`stroke-linecap="round"`,
		`stroke-linejoin="round"`,
		`aria-hidden="true"`,
		`focusable="false"`,
		`<path d="M6 9l6 6 6-6"/>`,
	} {
		if !contains(got, pflicht) {
			t.Errorf("die Hülle enthält %q nicht:\n%s", pflicht, got)
		}
	}
	// Feste Maße würden die Größe am Einsatzort festnageln; sie soll über
	// CSS bestimmt werden.
	for _, verboten := range []string{` width="`, ` height="`} {
		if contains(got, verboten) {
			t.Errorf("die Hülle enthält %q — die Größe gehört ins CSS", verboten)
		}
	}
}

func TestGefuellteHuelleFuerFlaechensymbole(t *testing.T) {
	// Die Mondphasen sind die einzige Ausnahme von der Strichzeichnung:
	// bei ihnen trägt die Flächenaufteilung die Aussage.
	got := string(huelle(`<circle cx="12" cy="12" r="9"/>`, true))
	if !contains(got, `fill="currentColor"`) {
		t.Errorf("gefüllte Hülle nutzt kein fill=currentColor:\n%s", got)
	}
	if contains(got, `fill="none"`) {
		t.Errorf("gefüllte Hülle enthält noch fill=none:\n%s", got)
	}
}

func contains(h, n string) bool {
	return len(n) <= len(h) && indexOf(h, n) >= 0
}

func indexOf(h, n string) int {
	for i := 0; i+len(n) <= len(h); i++ {
		if h[i:i+len(n)] == n {
			return i
		}
	}
	return -1
}
