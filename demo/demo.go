// Package demo liefert die Referenzseite des Design-Systems. Sie ist von
// den drei Kernfunktionen (TokensCSS, BaseCSS, CSS) im Wurzelpaket getrennt,
// damit ein Dienst, der nur das CSS einbindet, nicht auch index.html in sein
// Binary bekommt — docs/design.md hält die Go-Schnittstelle des
// Wurzelpakets bewusst auf drei Funktionen und einen Kontrastrechner klein.
package demo

import (
	"bytes"
	_ "embed"
	"net/http"
	"strings"

	designsystem "github.com/jobrunner/fieldworksdiary-designsystem"
	"github.com/jobrunner/fieldworksdiary-designsystem/icons"
)

//go:embed index.html
var indexHTML []byte

// symbolPlatzhalter markiert die Stelle in index.html, an der die
// Symbolübersicht eingesetzt wird. index.html ist statisch, die Symbole
// kommen aber aus dem Paket icons — seite() setzt sie beim Ausliefern (und
// in den Tests) ein, statt sie in HTML zu wiederholen.
const symbolPlatzhalter = "__SYMBOLE__"

// seite liefert die Referenzseite mit eingesetzter Symbolübersicht. Handler
// und Tests rufen dieselbe Funktion auf, damit die ausgelieferte Seite und
// das, was die Tests prüfen, nicht auseinanderlaufen können.
func seite() []byte {
	return bytes.ReplaceAll(indexHTML, []byte(symbolPlatzhalter), []byte(symbolGalerie()))
}

// symbolGalerie baut die Übersicht aller Symbole aus icons.Alle(), geordnet
// nach icons.Gruppen() statt alphabetisch — bei den Mondphasen ist die
// Abfolge selbst die Aussage. Jede Kachel trägt data-icon="<name>", damit
// der Test sie findet und beim Ansehen erkennbar ist, welches Symbol man
// vor sich hat.
func symbolGalerie() string {
	alle := icons.Alle()
	var b strings.Builder
	for _, gruppe := range icons.Gruppen() {
		b.WriteString(`<h3 class="icon-gruppe-titel">`)
		b.WriteString(gruppe.Titel)
		b.WriteString(`</h3><div class="icon-galerie">`)
		for _, name := range gruppe.Namen {
			svg := string(alle[name])
			// Die Icon-Funktionen liefern das SVG ohne class="icon" — die
			// Größe bestimmt der Einsatzort, hier über dieselbe Klasse wie
			// überall sonst im Design-System.
			svg = strings.Replace(svg, "<svg ", `<svg class="icon" `, 1)
			b.WriteString(`<div class="icon-kachel" data-icon="`)
			b.WriteString(name)
			b.WriteString(`">`)
			b.WriteString(svg)
			b.WriteString(`<span class="icon-name">`)
			b.WriteString(name)
			b.WriteString(`</span></div>`)
		}
		b.WriteString(`</div>`)
	}
	return b.String()
}

// Handler liefert die Referenzseite samt Stylesheet. Sie dient der
// visuellen Abnahme des Systems und ist zugleich das kürzeste Beispiel,
// wie ein Dienst das CSS einbindet:
//
//	go run ./cmd/demo   (oder im Test: httptest.NewServer(demo.Handler()))
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/designsystem.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(designsystem.CSS())
	})
	// Eigene Route statt eines eingebetteten <script>-Blocks in index.html:
	// aus demselben Grund, aus dem designsystem.css nicht als <style>
	// eingebettet wird. Ein Skript, dessen Verhalten direkt in der
	// Referenzseite stünde, würde ein kaputtes JS() lokal überdecken — die
	// Demo bindet die Combobox-Logik deshalb als ES-Modul über diese Route
	// ein und wickelt nur die zehn erfundenen Beispieleinträge in einem
	// kurzen Inline-Skript, das JS() importiert statt es nachzubauen.
	mux.HandleFunc("/designsystem.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Write(designsystem.JS())
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(seite())
	})
	return mux
}
