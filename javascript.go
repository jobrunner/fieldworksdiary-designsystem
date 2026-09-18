package designsystem

import (
	_ "embed"
	"slices"
)

//go:embed js/combobox.js
var comboboxJS []byte

// JS liefert das Verhalten der Komponenten, die ein Skript brauchen —
// bislang allein die Combobox mit Vorschlagsliste. Ein einziges ES-Modul
// ohne Bündler und ohne Abhängigkeiten, genau wie CSS() ein einziges
// Stylesheet liefert.
//
// Liefert eine Kopie, keine Referenz auf den eingebetteten Puffer — aus
// demselben Grund wie TokensCSS() und BaseCSS(): []byte ist veränderlich,
// ein Aufrufer, der das Ergebnis in-place verändert, würde sonst den
// Prozesszustand für jeden weiteren Aufruf und jeden anderen Aufrufer im
// selben Prozess beschädigen.
func JS() []byte { return slices.Clone(comboboxJS) }
