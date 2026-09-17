package designsystem

import (
	_ "embed"
	"slices"
)

//go:embed css/tokens.css
var tokensCSS []byte

//go:embed css/base.css
var baseCSS []byte

// TokensCSS liefert ausschließlich die Variablen — Farben für beide Themen,
// Typografie, Abstände, Radien. Wer nur die Werte braucht und seine
// Komponenten selbst mitbringt, bindet allein diese Datei ein.
//
// Liefert eine Kopie, keine Referenz auf den eingebetteten Puffer: ein
// Aufrufer, der das Ergebnis in-place verändert (z. B. anhängt, sortiert,
// mit append umschreibt), würde sonst den Prozesszustand für jeden
// weiteren Aufruf und jeden anderen Aufrufer im selben Prozess beschädigen
// — []byte ist veränderlich, embed.FS gibt hier nur den rohen Puffer her.
func TokensCSS() []byte { return slices.Clone(tokensCSS) }

// BaseCSS liefert Reset und Komponenten. Setzt die Tokens voraus und ist
// ohne sie wirkungslos — jede Farbe darin ist eine Variable.
//
// Liefert wie TokensCSS eine Kopie, aus demselben Grund.
func BaseCSS() []byte { return slices.Clone(baseCSS) }

// CSS liefert beides in der einzig gültigen Reihenfolge: erst die
// Variablen, dann was sie verwendet. Baut ohnehin einen neuen Puffer und
// gibt damit nie eine Referenz auf tokensCSS oder baseCSS heraus.
func CSS() []byte {
	out := make([]byte, 0, len(tokensCSS)+len(baseCSS)+1)
	out = append(out, tokensCSS...)
	out = append(out, '\n')
	return append(out, baseCSS...)
}
