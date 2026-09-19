package designsystem

import (
	"html/template"
	"strings"
	"testing"
)

func TestKopfzeileTraegtNameUndUntertitel(t *testing.T) {
	got := string(Kopfzeile(KopfDaten{Name: "Ortus", Untertitel: "Point-in-Polygon Abfrage über Datenquellen"}))
	for _, teil := range []string{"<header", "<h1>Ortus</h1>", "Point-in-Polygon Abfrage über Datenquellen"} {
		if !strings.Contains(got, teil) {
			t.Errorf("Kopfzeile enthält %q nicht:\n%s", teil, got)
		}
	}
}

func TestKopfzeileMaskiertFremdtext(t *testing.T) {
	// Name und Untertitel kommen aus der Konfiguration eines Dienstes.
	// Ohne Maskierung wäre das eine Einschleusstelle für Markup.
	got := string(Kopfzeile(KopfDaten{Name: `<script>alert(1)</script>`, Untertitel: "x"}))
	if strings.Contains(got, "<script>") {
		t.Errorf("Kopfzeile reicht Markup ungefiltert durch:\n%s", got)
	}
}

func TestFusszeileZeigtVerweiseUndFassung(t *testing.T) {
	got := string(Fusszeile(FussDaten{
		Verweise: []Verweis{
			{Text: "API Dokumentation", Ziel: "/docs"},
			{Text: "OpenAPI Spec", Ziel: "/openapi.json"},
			{Text: "Health Status", Ziel: "/health"},
		},
		Name:    "ortus",
		Fassung: "1.4.2",
	}))
	for _, teil := range []string{
		`<footer`, `href="/docs"`, "API Dokumentation",
		`href="/health"`, "ortus", "1.4.2",
	} {
		if !strings.Contains(got, teil) {
			t.Errorf("Fußzeile enthält %q nicht:\n%s", teil, got)
		}
	}
}

// TestFusszeileMaskiertFremdtext prüft die Maskierung der Fußzeile — nicht
// nur der Verweistext, sondern insbesondere das VerweisZIEL im
// href-Attribut: dort liegt der gefährlichere Kontext, weil ein
// unmaskiertes Ziel nicht nur sichtbares Markup einschleust (wie im Text),
// sondern eine ausführbare javascript:-URL. layout_test.go prüfte bislang
// nur die Kopfzeile; die Fußzeile blieb ungeprüft. Belegt in der
// Schlussprüfung: Verweis.Ziel von string auf template.URL geändert →
// grün, obwohl damit javascript:alert(1) ungefiltert ins href liefe — denn
// html/template vertraut einem template.URL-Wert als bereits geprüft und
// lässt jedes Schema durch, während ein einfacher string durch den
// URL-Filter läuft.
func TestFusszeileMaskiertFremdtext(t *testing.T) {
	got := string(Fusszeile(FussDaten{
		Verweise: []Verweis{
			{Text: `<script>alert(1)</script>`, Ziel: "javascript:alert(1)"},
		},
		Name:    "ortus",
		Fassung: "1.0.0",
	}))
	if strings.Contains(got, "<script>") {
		t.Errorf("Fußzeile reicht Text ungefiltert durch:\n%s", got)
	}
	if strings.Contains(got, `href="javascript:`) {
		t.Errorf("Fußzeile reicht das Verweisziel ungefiltert ins href — javascript:-URL käme unmaskiert durch:\n%s", got)
	}
}

// TestFusszeileVerweisZielAlsTemplateURLWuerdeEinschleusung erlauben belegt
// TestFusszeileMaskiertFremdtext: ein Verweis.Ziel vom Typ template.URL
// (statt string) gilt html/template als bereits geprüft und läuft NICHT
// mehr durch den URL-Filter — genau die Mutation aus dem Befund. Ein Test,
// der nie rot war, beweist nichts: dieser Test baut denselben Fall direkt
// mit einer eigenen, zu Verweis strukturgleichen Vorlage nach, um zu
// zeigen, dass die Prüfung oben diesen Unterschied tatsächlich sieht.
func TestFusszeileVerweisZielAlsTemplateURLWuerdeEinschleusungErlauben(t *testing.T) {
	unsicher := template.Must(template.New("unsicher").Parse(
		`<a href="{{.Ziel}}">{{.Text}}</a>`))
	var b strings.Builder
	if err := unsicher.Execute(&b, struct {
		Text string
		Ziel template.URL // die Mutation aus dem Befund
	}{Text: "x", Ziel: template.URL("javascript:alert(1)")}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), `href="javascript:`) {
		t.Fatalf("Testaufbau fehlerhaft: template.URL sollte den Filter umgehen und das javascript:-Schema durchlassen, tat es hier nicht: %s", b.String())
	}
	// Und zur Gegenprobe: dieselbe Vorlage mit Ziel als string (wie
	// Verweis es tatsächlich deklariert) filtert das javascript:-Schema
	// vollständig heraus (href wird zu "#ZgotmplZ").
	sicher := template.Must(template.New("sicher").Parse(
		`<a href="{{.Ziel}}">{{.Text}}</a>`))
	var b2 strings.Builder
	if err := sicher.Execute(&b2, struct{ Text, Ziel string }{Text: "x", Ziel: "javascript:alert(1)"}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(b2.String(), `href="javascript:`) {
		t.Errorf("Ziel als string sollte das javascript:-Schema herausfiltern, tat es hier nicht: %s", b2.String())
	}
}

func TestFusszeileOhneVerweiseBleibtGueltig(t *testing.T) {
	// Expertus hat andere Ziele als die vier Go-Dienste und kommt
	// möglicherweise ganz ohne Verweisliste aus.
	got := string(Fusszeile(FussDaten{Name: "expertus", Fassung: "0.1.0"}))
	if !strings.Contains(got, "<footer") || !strings.Contains(got, "expertus") {
		t.Errorf("Fußzeile ohne Verweise ist unbrauchbar:\n%s", got)
	}
	// Kein Trennzeichen ohne etwas zu trennen.
	if strings.Contains(got, "·") {
		t.Errorf("Fußzeile ohne Verweise enthält ein Trennzeichen:\n%s", got)
	}
}
