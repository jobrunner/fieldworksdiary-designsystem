package designsystem

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

// klassenRe findet die Klassennamen, die base.css definiert.
var klassenRe = regexp.MustCompile(`^\.([a-z][a-z0-9-]*)`)

func TestDemoZeigtJedeKomponente(t *testing.T) {
	// Welche Klassen base.css anbietet, steht in base.css — nicht in einer
	// Liste hier, die beim nächsten Zuwachs vergessen würde.
	angeboten := map[string]bool{}
	for _, zeile := range strings.Split(string(BaseCSS()), "\n") {
		if m := klassenRe.FindStringSubmatch(strings.TrimSpace(zeile)); m != nil {
			angeboten[m[1]] = true
		}
	}
	// Zustandsklassen erscheinen nur zusammen mit ihrer Grundklasse.
	for _, nur := range []string{"active"} {
		delete(angeboten, nur)
	}

	demo := string(demoHTML)
	for klasse := range angeboten {
		if !strings.Contains(demo, `"`+klasse+`"`) && !strings.Contains(demo, `"`+klasse+` `) &&
			!strings.Contains(demo, ` `+klasse+`"`) && !strings.Contains(demo, ` `+klasse+` `) {
			t.Errorf("Die Demo-Seite zeigt .%s nicht — die Klasse wäre nirgends vor dem Einsatz zu sehen", klasse)
		}
	}
}

func TestDemoHandlerLiefertSeiteUndCSS(t *testing.T) {
	srv := httptest.NewServer(DemoHandler())
	defer srv.Close()

	for _, f := range []struct {
		pfad string
		typ  string
		teil string
	}{
		{"/", "text/html; charset=utf-8", "<!doctype html>"},
		{"/designsystem.css", "text/css; charset=utf-8", "--control-line:"},
	} {
		res, err := http.Get(srv.URL + f.pfad)
		if err != nil {
			t.Fatalf("GET %s: %v", f.pfad, err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Errorf("GET %s: Status %d, erwartet 200", f.pfad, res.StatusCode)
		}
		if got := res.Header.Get("Content-Type"); got != f.typ {
			t.Errorf("GET %s: Content-Type %q, erwartet %q", f.pfad, got, f.typ)
		}
		buf := make([]byte, 4096)
		n, _ := res.Body.Read(buf)
		if !strings.Contains(strings.ToLower(string(buf[:n])), strings.ToLower(f.teil)) {
			t.Errorf("GET %s: %q kommt im Anfang der Antwort nicht vor", f.pfad, f.teil)
		}
	}
}
