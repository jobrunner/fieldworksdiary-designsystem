// Package demo liefert die Referenzseite des Design-Systems. Sie ist von
// den drei Kernfunktionen (TokensCSS, BaseCSS, CSS) im Wurzelpaket getrennt,
// damit ein Dienst, der nur das CSS einbindet, nicht auch index.html in sein
// Binary bekommt — docs/design.md hält die Go-Schnittstelle des
// Wurzelpakets bewusst auf drei Funktionen und einen Kontrastrechner klein.
package demo

import (
	_ "embed"
	"net/http"

	designsystem "github.com/jobrunner/fieldworksdiary-designsystem"
)

//go:embed index.html
var indexHTML []byte

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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})
	return mux
}
