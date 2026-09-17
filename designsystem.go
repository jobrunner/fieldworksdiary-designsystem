package designsystem

import _ "embed"

//go:embed css/tokens.css
var tokensCSS []byte

//go:embed css/base.css
var baseCSS []byte

// TokensCSS liefert ausschließlich die Variablen — Farben für beide Themen,
// Typografie, Abstände, Radien. Wer nur die Werte braucht und seine
// Komponenten selbst mitbringt, bindet allein diese Datei ein.
func TokensCSS() []byte { return tokensCSS }

// BaseCSS liefert Reset und Komponenten. Setzt die Tokens voraus und ist
// ohne sie wirkungslos — jede Farbe darin ist eine Variable.
func BaseCSS() []byte { return baseCSS }

// CSS liefert beides in der einzig gültigen Reihenfolge: erst die
// Variablen, dann was sie verwendet.
func CSS() []byte {
	out := make([]byte, 0, len(tokensCSS)+len(baseCSS)+1)
	out = append(out, tokensCSS...)
	out = append(out, '\n')
	return append(out, baseCSS...)
}
