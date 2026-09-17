// Package designsystem liefert die gemeinsame Gestaltungsgrundlage der
// fieldworksdiary-Dienste als CSS und hält die Kontrastzusagen des Systems
// prüfbar.
package designsystem

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// RGB ist eine Farbe im sRGB-Raum, je Kanal 0..255.
type RGB struct{ R, G, B uint8 }

// ParseHex liest #rgb und #rrggbb. Andere Schreibweisen werden abgelehnt
// statt stillschweigend gedeutet: ein Tippfehler in einem Farbwert soll den
// Build brechen, nicht zu einer zufälligen Farbe führen.
func ParseHex(s string) (RGB, error) {
	if !strings.HasPrefix(s, "#") {
		return RGB{}, fmt.Errorf("Farbwert %q beginnt nicht mit #", s)
	}
	h := s[1:]
	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	if len(h) != 6 {
		return RGB{}, fmt.Errorf("Farbwert %q hat weder 3 noch 6 Stellen", s)
	}
	n, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return RGB{}, fmt.Errorf("Farbwert %q ist nicht hexadezimal: %w", s, err)
	}
	return RGB{R: uint8(n >> 16), G: uint8(n >> 8), B: uint8(n)}, nil
}

// Luminance ist die relative Leuchtdichte nach WCAG 2.x.
func (c RGB) Luminance() float64 {
	lin := func(v uint8) float64 {
		f := float64(v) / 255
		if f <= 0.03928 {
			return f / 12.92
		}
		return math.Pow((f+0.055)/1.055, 2.4)
	}
	return 0.2126*lin(c.R) + 0.7152*lin(c.G) + 0.0722*lin(c.B)
}

// ContrastRatio liefert das Kontrastverhältnis zweier Farben, immer >= 1.
// Die Reihenfolge der Argumente spielt keine Rolle.
func ContrastRatio(a, b RGB) float64 {
	l1, l2 := a.Luminance(), b.Luminance()
	if l2 > l1 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}
