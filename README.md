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
Marken, Tabellen, Zustände, Akkordeon, Combobox. `icons/` — 40 Symbole als
fertiges SVG (siehe unten). Das Wurzelpaket liefert außerdem Kopf- und
Fußzeile (`Kopfzeile()`, `Fusszeile()`) sowie das Skript der Combobox
(`JS()`).

Fachliche Komponenten gehören nicht hierher: die Wetterkarte bleibt in Tempus,
die Quellenliste in Ortus.

## Herkunft der Palette

Die Akzentfarbe ist Grün und stammt aus der iOS-App `fieldworksdiary-ng`, die
Grün als Markenfarbe führt — Web und App sollen zusammenpassen. Web und App
weisen ihr unterschiedliche Rollen zu (App: Fläche mit weißer Schrift; Web:
zusätzlich Textfarbe für Links, Reiter und Erfolg), deshalb gibt es zwei
Token: `--accent` (Fläche, in beiden Themen gleich) und `--accent-text`
(Textfarbe, themenabhängig). Details und die gemessenen Kontraste stehen in
[docs/design.md](docs/design.md#herkunft-der-palette).

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

## Symbole

    import "github.com/jobrunner/fieldworksdiary-designsystem/icons"

    icons.Standort()      // <svg …> für den Standort-Knopf
    icons.WetterNebel()   // Zeichnung; welcher WMO-Code das ist, weiß der Dienst
    icons.Alle()          // alle 40 Symbole unter ihrem Namen ("standort", "wetter-nebel", …)
    icons.Gruppen()       // dieselben Symbole, nach Bedeutung geordnet statt alphabetisch

Jedes Symbol folgt derselben Bauregel: 24er-Raster, keine festen Maße,
`currentColor`, `aria-hidden`. Ein Symbol nimmt damit die Farbe seines
Umfelds an — es folgt Thema und Zustand, ohne eine eigene Farbe zu kennen.
Die Größe bestimmt das CSS über die Klasse `.icon`.

Einzige Ausnahme sind die acht Mondphasen: sie nutzen `fill="currentColor"`
statt einer reinen Strichzeichnung, weil dort die Flächenaufteilung
zwischen beleuchtetem und unbeleuchtetem Teil die Aussage trägt.

**Achtung bei `html/template`:** `icons.Icon` ist ein `string`, keine
Auszeichnung — `{{.Symbol}}` in einer `html/template`-Vorlage maskiert das
SVG und zeigt sichtbaren Quelltext statt eines Symbols. Immer
`{{.Symbol.HTML}}` schreiben (liefert `template.HTML`), niemals `{{.Symbol}}`
direkt. `Kopfzeile()`/`Fusszeile()` liefern zum Vergleich bereits
`template.HTML` — bei Icon ist der zusätzliche Schritt bewusst nötig, damit
das Paket `icons` selbst ohne die Abhängigkeit auf `html/template` bleibt.

    icons.MitKlasse(icons.Standort(), "icon")   // setzt class="icon" aufs <svg>
    icons.MitBeschriftung(icons.Suche(), "Suchen") // role="img" + aria-label
                                                     // statt aria-hidden, für
                                                     // ein Symbol OHNE
                                                     // begleitenden Text
                                                     // (z. B. in einem Knopf)

`icons.Alle()` liefert die Symbole unsortiert (eine `map`); `icons.Gruppen()`
liefert dieselbe Menge als geordnete Liste von Gruppen — Bedienung, Wetter,
Mondphasen (in ihrer natürlichen Abfolge: neu, zunehmende Sichel, erstes
Viertel, zunehmend gibbous, voll, abnehmend gibbous, letztes Viertel,
abnehmende Sichel), Sonnenstände, Messwerte. Ein Dienst, der eine
Symbolauswahl anbietet, kann diese Reihenfolge übernehmen, statt sie
nachzubauen.

Das Modul liefert Zeichnungen, keine Fachlogik: welcher Zustand welches
Symbol bekommt, entscheidet der Dienst.

## Seitengerüst

    designsystem.Kopfzeile(designsystem.KopfDaten{Name: "Ortus", Untertitel: "…"})
    designsystem.Fusszeile(designsystem.FussDaten{
        Verweise: []designsystem.Verweis{{Text: "Health Status", Ziel: "/health"}},
        Name:     "ortus",
        Fassung:  "1.4.2",
    })

Name, Untertitel und Verweise übergibt der Dienst; Auszeichnung und Klassen
(`ds-kopf`, `ds-fuss`, …) kommen aus dem Modul. Beide Funktionen geben
`template.HTML` zurück; alle Werte laufen vorher durch `html/template` und
werden maskiert.

## Akkordeon

    <details class="akkordeon">
      <summary><span>Optionen</span><!-- Chevron-Symbol --></summary>
      <div class="akkordeon-inhalt">…</div>
    </details>

Kein eigener Baustein in Go — `class="akkordeon"` auf einem nativen
`<details>`/`<summary>`-Paar reicht. Auf- und Zuklappen, Tastaturbedienung
und die Ansage für Screenreader liefert der Browser; das Modul steuert nur
die Gestaltung, unter anderem die Drehung des Chevron-Symbols im
geöffneten Zustand.

## Combobox mit Vorschlagsliste

    import (
        designsystem "github.com/jobrunner/fieldworksdiary-designsystem"
    )

    mux.HandleFunc("/designsystem.js", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
        w.Write(designsystem.JS())
    })

Im HTML als ES-Modul einbinden und die Logik über `mountCombobox` an ein
vorhandenes Eingabefeld samt Vorschlagsliste hängen:

    <div class="form-group combobox">
      <label for="art">Art</label>
      <input type="text" id="art" role="combobox" aria-expanded="false"
             aria-autocomplete="list" aria-controls="art-liste" autocomplete="off" />
      <ul id="art-liste" class="combobox-liste" role="listbox" hidden></ul>
    </div>

    <script type="module">
      import { mountCombobox } from "/designsystem.js";

      mountCombobox({
        input: document.getElementById("art"),
        listbox: document.getElementById("art-liste"),
        suggest: async (query) => { /* Vorschläge zu query liefern */ },
        onPick: (eintrag) => { /* Auswahl verarbeiten */ },
      });
    </script>

`JS()` liefert ein einziges, abhängigkeitsfreies ES-Modul — kein Bündler
nötig, ebenso wie `CSS()` ein einziges Stylesheet liefert. Die Fachlogik
(woher die Vorschläge kommen) bleibt beim Dienst; das Modul liefert nur
Tastaturbedienung, ARIA-Zustände und Gestaltung der Liste.

## Referenzseite

    go run ./cmd/demo

Zeigt jede Komponente unter <http://127.0.0.1:5180>. Das Thema folgt der
Systemeinstellung.

## Entwurf

[docs/design.md](docs/design.md) — warum das System so aussieht und wie die
fünf Dienste darauf umgestellt werden.
