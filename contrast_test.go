package designsystem

import (
	"math"
	"testing"
)

func TestContrastRatio(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want float64
	}{
		{"Schwarz auf Weiß ist das Maximum", "#000000", "#ffffff", 21},
		{"eine Farbe gegen sich selbst", "#2563eb", "#2563eb", 1},
		{"Reihenfolge ist ohne Belang", "#ffffff", "#000000", 21},
		{"Kurzschreibweise wird verstanden", "#fff", "#000", 21},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := ParseHex(tt.a)
			if err != nil {
				t.Fatalf("ParseHex(%q): %v", tt.a, err)
			}
			b, err := ParseHex(tt.b)
			if err != nil {
				t.Fatalf("ParseHex(%q): %v", tt.b, err)
			}
			got := ContrastRatio(a, b)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("ContrastRatio = %.2f, erwartet %.2f", got, tt.want)
			}
		})
	}
}

func TestParseHexLehntUngueltigesAb(t *testing.T) {
	for _, s := range []string{"", "#", "#12345", "#gggggg", "2563eb "} {
		if _, err := ParseHex(s); err == nil {
			t.Errorf("ParseHex(%q) hat keinen Fehler geliefert", s)
		}
	}
}
