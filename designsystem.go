package designsystem

import _ "embed"

//go:embed css/tokens.css
var tokensCSS []byte

// TokensCSS liefert ausschließlich die Variablen — Farben für beide Themen,
// Typografie, Abstände, Radien. Wer nur die Werte braucht und seine
// Komponenten selbst mitbringt, bindet allein diese Datei ein.
func TokensCSS() []byte { return tokensCSS }
