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

// classAttrRe findet alle class="..."-Attribute in HTML.
var classAttrRe = regexp.MustCompile(`class="([^"]*)"`)

// benutzteKlassen extrahiert die tatsächlich vergebenen Klassennamen aus HTML,
// indem es nur class="..."-Attribute berücksichtigt. So werden id-Attribute,
// Klassenerwähnungen im Fließtext oder per display:none versteckte Elemente
// nicht fälschlicherweise als „gezeigt" gewertet.
func benutzteKlassen(html string) map[string]bool {
	klassen := make(map[string]bool)
	for _, m := range classAttrRe.FindAllStringSubmatch(html, -1) {
		if len(m) > 1 {
			for _, k := range strings.Fields(m[1]) {
				klassen[k] = true
			}
		}
	}
	return klassen
}

func TestBenutzteKlassenExtrahiertNurAusClassAttributen(t *testing.T) {
	tests := []struct {
		html         string
		erwartet     map[string]bool
		beschreibung string
	}{
		{
			`<span id="btn">Text</span>`,
			map[string]bool{},
			"id-Attribut ist keine Klasse",
		},
		{
			`<div class="card btn-row">Text</div>`,
			map[string]bool{"card": true, "btn-row": true},
			"mehrere Klassen im Attribut",
		},
		{
			`<p>Die Klasse btn ist blau.</p>`,
			map[string]bool{},
			"Klassenname im Text ergibt keine Klasse",
		},
		{
			`<div class="card">A</div><div class="btn">B</div>`,
			map[string]bool{"card": true, "btn": true},
			"mehrere Attribute zusammen",
		},
	}

	for _, test := range tests {
		got := benutzteKlassen(test.html)
		if len(got) != len(test.erwartet) {
			t.Errorf("%s: Länge %d, erwartet %d", test.beschreibung, len(got), len(test.erwartet))
		}
		for k := range test.erwartet {
			if !got[k] {
				t.Errorf("%s: .%s fehlt", test.beschreibung, k)
			}
		}
		for k := range got {
			if !test.erwartet[k] {
				t.Errorf("%s: unerwartete .%s", test.beschreibung, k)
			}
		}
	}
}

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
	benutzt := benutzteKlassen(demo)
	for klasse := range angeboten {
		if !benutzt[klasse] {
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
