# fieldworksdiary-designsystem

Gemeinsame Gestaltungsgrundlage der fieldworksdiary-Dienste — Ortus, Tempus,
Situs, Hostus und Expertus.

## Einbinden

    go get github.com/jobrunner/fieldworksdiary-designsystem

Das CSS unter einer eigenen Route ausliefern:

    import designsystem "github.com/jobrunner/fieldworksdiary-designsystem"

    mux.HandleFunc("/assets/designsystem.css", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/css; charset=utf-8")
        w.Write(designsystem.CSS())
    })

Im HTML verlinken:

    <link rel="stylesheet" href="/assets/designsystem.css">

`TokensCSS()` und `BaseCSS()` liefern die Teile einzeln, falls ein Dienst nur
die Variablen braucht.

## Was drin ist

`css/tokens.css` — Farben für helles und dunkles Thema, Typografie, Abstände,
Radien. `css/base.css` — Reset, Formularelemente, Knöpfe, Karten, Reiter,
Marken, Tabellen, Zustände.

Fachliche Komponenten gehören nicht hierher: die Wetterkarte bleibt in Tempus,
die Quellenliste in Ortus.

## Zusagen

Alle Textfarben erreichen 7:1 gegen Seitenhintergrund und Kartenfläche, in
beiden Themen. Ränder von Bedienelementen erreichen 3:1. `go test ./...`
rechnet das nach — die Werte werden aus `css/tokens.css` gelesen, nicht im
Testcode wiederholt.

`base.css` darf keine Farbe im Klartext enthalten; auch das prüft der Test.
Eine Farbe außerhalb von `tokens.css` stünde außerhalb der Kontrastprüfung.

Geprüft wird außerdem jede Regel in `base.css`, die `color: var(--…)` UND
`background`/`background-color: var(--…)` **im selben Block** setzt: die
beiden Token müssen gegeneinander ebenfalls 7:1 erreichen — zwei einzeln
geprüfte Token lassen sich sonst zu einem unlesbaren Paar kombinieren (etwa
`color: var(--text-muted)` auf `background: var(--accent)`, 1.15:1). Diese
Prüfung hat eine bewusste Grenze: sie sieht nur Paare, die *in derselben
Regel* stehen. Eine Farbe, die eine Regel wie `.btn:hover` aus einer anderen
Regel (`.btn`) erbt, weil sie dort nicht neu gesetzt wird, entgeht ihr — das
verlangte den vollen CSS-Kaskadenalgorithmus, den diese einfache
Textzerlegung nicht nachbildet. Mit anderen Worten: geprüft ist, dass jedes
Token für sich gegen die Flächen besteht, und dass jede Regel, die Text- und
Flächenfarbe gemeinsam benennt, zueinander passt — nicht, dass *jede* im CSS
tatsächlich zustande kommende Kombination aus Kaskade und Vererbung
zueinander passt.

## Referenzseite

    go run ./cmd/demo

Zeigt jede Komponente unter <http://127.0.0.1:5180>. Das Thema folgt der
Systemeinstellung.

## Entwurf

[docs/design.md](docs/design.md) — warum das System so aussieht und wie die
fünf Dienste darauf umgestellt werden.
