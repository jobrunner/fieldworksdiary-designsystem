package designsystem

import (
	"html/template"
	"strings"
)

// KopfDaten trägt, was einen Dienst im Seitenkopf unterscheidet. Alles
// andere — Auszeichnung, Klassen, Aufbau — kommt aus dem Modul, damit Ortus
// und Tempus nicht länger dieselbe Kopfzeile getrennt pflegen.
type KopfDaten struct {
	Name       string
	Untertitel string
}

// Verweis ist ein Eintrag der Fußzeile.
type Verweis struct {
	Text string
	Ziel string
}

// FussDaten trägt die Verweise und die Fassungsangabe.
type FussDaten struct {
	Verweise []Verweis
	Name     string
	Fassung  string
}

// Alle Werte laufen durch die Maskierung von html/template: Name, Untertitel
// und Verweisziele stammen aus der Konfiguration eines Dienstes und sind
// damit Fremdtext. Ohne Maskierung wäre das eine Einschleusstelle.
var (
	kopfTpl = template.Must(template.New("kopf").Parse(
		`<header class="ds-kopf">` +
			`<div class="ds-kopf-titel"><h1>{{.Name}}</h1>` +
			`{{if .Untertitel}}<p class="muted">{{.Untertitel}}</p>{{end}}</div>` +
			`</header>`))

	fussTpl = template.Must(template.New("fuss").Parse(
		`<footer class="ds-fuss">` +
			`{{if .Verweise}}<p class="ds-fuss-verweise">` +
			`{{range $i, $v := .Verweise}}{{if $i}} &middot; {{end}}` +
			`<a href="{{$v.Ziel}}">{{$v.Text}}</a>{{end}}</p>{{end}}` +
			`<p class="ds-fuss-fassung muted">{{.Name}} {{.Fassung}}</p>` +
			`</footer>`))
)

// Kopfzeile liefert den Seitenkopf.
func Kopfzeile(k KopfDaten) template.HTML {
	var b strings.Builder
	if err := kopfTpl.Execute(&b, k); err != nil {
		// Die Vorlage ist fest eingebaut und wurde beim Start übersetzt;
		// ein Fehler hier wäre ein Programmierfehler, kein Laufzeitfall.
		return template.HTML("")
	}
	return template.HTML(b.String())
}

// Fusszeile liefert den Seitenfuß.
func Fusszeile(f FussDaten) template.HTML {
	var b strings.Builder
	if err := fussTpl.Execute(&b, f); err != nil {
		return template.HTML("")
	}
	return template.HTML(b.String())
}
