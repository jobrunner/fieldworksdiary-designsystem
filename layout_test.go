package designsystem

import (
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
