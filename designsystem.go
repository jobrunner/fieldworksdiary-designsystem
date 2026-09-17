package designsystem

import (
	_ "embed"
	"net/http"
)

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

//go:embed demo/index.html
var demoHTML []byte

// DemoHandler liefert die Referenzseite samt Stylesheet. Sie dient der
// visuellen Abnahme des Systems und ist zugleich das kürzeste Beispiel,
// wie ein Dienst das CSS einbindet:
//
//	go run ./cmd/demo   (oder im Test: httptest.NewServer(DemoHandler()))
func DemoHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/designsystem.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(CSS())
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(demoHTML)
	})
	return mux
}
