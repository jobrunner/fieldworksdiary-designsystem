package designsystem

import (
	_ "embed"
)

//go:embed js/combobox.js
var comboboxJS []byte

//go:embed js/koordinaten.js
var koordinatenJS []byte

// JS liefert das Verhalten der Komponenten, die ein Skript brauchen: die
// Combobox mit Vorschlagsliste und die Koordinateneingabe. Ein einziges
// ES-Modul ohne Bündler und ohne Abhängigkeiten, genau wie CSS() ein
// einziges Stylesheet liefert — beide Quelldateien sind für sich
// eigenständige ES-Module (keine importiert die andere) und werden hier
// aneinandergehängt, damit ein Dienst weiterhin nur eine einzige Datei
// unter /designsystem.js einbinden muss.
//
// Das Ergebnis entsteht als frisch angelegter Puffer (append auf ein mit
// make() erzeugtes []byte, nicht auf einen der eingebetteten Puffer
// selbst) — aus demselben Grund wie slices.Clone() bei TokensCSS() und
// BaseCSS(): ein Aufrufer, der das Ergebnis in-place verändert, darf die
// eingebetteten Puffer nicht für jeden weiteren Aufruf beschädigen.
func JS() []byte {
	kombiniert := make([]byte, 0, len(comboboxJS)+1+len(koordinatenJS))
	kombiniert = append(kombiniert, comboboxJS...)
	kombiniert = append(kombiniert, '\n')
	kombiniert = append(kombiniert, koordinatenJS...)
	return kombiniert
}
