# Design-System für die fieldworksdiary-Dienste

Stand: 2026-09-17

## Ausgangslage

Fünf Dienste bringen je ein eigenes Frontend mit, keine zwei davon teilen
Code:

| Projekt | Form | Umfang | Palette | Dark Mode |
|---|---|---|---|---|
| Ortus | Go-String-Konstante in `frontend.go` | 2007 Zeilen | Blau `#2563eb` | nein |
| Tempus | `index.html` per `go:embed` | 2344 Zeilen | Blau `#2563eb` | nein |
| Situs | `explorer.html` per `go:embed` | 487 Zeilen | Grün `#2d6a4f` | ja |
| Hostus | `assets/style.css` per `go:embed` | 189 Zeilen | neutral | nein |
| Expertus | `styles.css`, von nginx ausgeliefert | 75 Zeilen | keine | nein |

Ortus und Tempus führen dasselbe System zweimal: identische Token-Namen,
identische Werte, identische Komponenten — entstanden durch Kopieren, seither
getrennt gepflegt. Tempus hat inzwischen ein `--text-disabled`, das Ortus
fehlt; Ortus hat einen Accordion, den Tempus nicht kennt. Die Schere geht mit
jeder Änderung weiter auf.

Gleichzeitig arbeiten Hostus und Expertus auf einem höheren
Barrierefreiheitsniveau als Ortus und Tempus. Expertus dokumentiert 7:1 für
Fließtext, Hostus trennt `--control-line` (≥ 3:1, WCAG 1.4.11) von der
dekorativen `--line`. Ortus' und Tempus' `--text-muted: #64748b` erreicht auf
dem Seitenhintergrund nur 4.55:1.

## Ziel

Ein gemeinsames Design-System, das die Dienste als Familie erkennbar macht,
das höhere Barrierefreiheitsniveau zum Standard erhebt und die Duplizierung
zwischen Ortus und Tempus auflöst.

## Entscheidungen

| Frage | Entscheidung |
|---|---|
| Auslieferung | Go-Modul `github.com/jobrunner/fieldworksdiary-designsystem` |
| Umfang | Tokens + Basiskomponenten; Domänenspezifisches bleibt lokal |
| Kontrast | AAA für Fließtext (7:1); Bedienelement-Ränder ≥ 3:1 |
| Markenbild | eine gemeinsame Akzentfarbe für alle fünf Dienste |
| Dark Mode | ja, über `prefers-color-scheme`; kein Umschalter, keine Persistenz |
| Git-Submodule | ausgeschlossen — nicht automatisch aktualisierbar |

Die Submodul-Variante scheidet aus, weil Dependabot und Renovate Submodule
nicht verlässlich aktualisieren, `go.mod` dagegen schon. Ortus, Tempus und
Hostus haben bereits eine `dependabot.yml`; Situs und Expertus bekommen eine.

Damit ein einziger Verteilmechanismus für alle fünf Dienste reicht, wird
Expertus von nginx auf einen eigenen Go-Server umgestellt (siehe unten).

## Aufbau des Repositories

```
fieldworksdiary-designsystem/
├── go.mod                  module github.com/jobrunner/fieldworksdiary-designsystem
├── designsystem.go         //go:embed css/*.css → TokensCSS(), BaseCSS(), CSS()
├── css/
│   ├── tokens.css          Farben (hell und dunkel), Typografie, Abstände,
│   │                       Radien, Schatten
│   └── base.css            Reset, Fokus, Formularelemente, Buttons, Card,
│                           Tabs, Badge, Tabelle, Spinner, Fehlerbox,
│                           sr-only, reduced-motion, Breakpoints
├── contrast_test.go        prüft jeden Farbwert programmatisch
├── demo/index.html         Referenzseite, alle Komponenten in hell und dunkel
└── docs/                   Anwendungsregeln und Migrationsleitfaden
```

Die Go-Schnittstelle bleibt bewusst klein: drei Funktionen, die Bytes liefern.
Kein Templating, keine Handler, keine Meinung darüber, wie ein Dienst sein
HTML baut.

```go
package designsystem

//go:embed css/tokens.css
//go:embed css/base.css

func TokensCSS() []byte  // nur die Variablen
func BaseCSS() []byte    // Reset und Komponenten, setzt Tokens voraus
func CSS() []byte        // beides, in dieser Reihenfolge
```

Ein Dienst liefert das Ergebnis unter einer eigenen Route aus, üblicherweise
`/assets/designsystem.css`, und verlinkt es aus seinem HTML. Wer sein CSS
heute inline in einem `<style>`-Block hält, kann `CSS()` auch dort einsetzen.

## Farb-Tokens

Jeder Wert ist gegen **beide** Flächen geprüft — Seitenhintergrund und Karte.
Welche der beiden die strengere ist, hängt vom Thema ab: im hellen Thema ist
es `--bg`, im dunklen `--card`. `contrast_test.go` liest die Werte aus
`css/tokens.css` und prüft beide Flächen; es ist die verbindliche Quelle, die
folgenden Tabellen geben sie wieder.

### Herkunft der Palette

Die Akzentfarbe ist seit dem ursprünglichen Entwurf blau (`#1e40af`). Die
zugehörige iOS-App (`fieldworksdiary-ng`) führt seit jeher Grün als
Markenfarbe; Web und App sollen zusammenpassen, deshalb wechselt das Modul
auf Grün.

Web und App weisen der Markenfarbe unterschiedliche Rollen zu: die App setzt
sie als **Fläche** (Knopf, Kopfbereich) mit weißer Schrift darauf, das Web
nutzt sie zusätzlich als **Textfarbe** (Link, ausgewählter Reiter, Erfolg).
Ein einziger Wert kann beide Rollen nicht tragen — ein Grün, das im dunklen
Thema als Text lesbar wäre, trägt weiße Schrift nur mit rund 2:1. Die Palette
teilt die Rolle deshalb auf zwei Token auf: `--accent` (Fläche, dunkel genug
für weiße Schrift, in beiden Themen derselbe Wert) und `--accent-text`
(Textfarbe, themenabhängig wie jede andere Textfarbe).

Die iOS-Werte selbst erfüllen die 7:1-Zusage dieses Moduls nicht; `--accent`
und `--accent-hover` sind deshalb minimal von den App-Werten nachgezogene
Werte — die iOS-App übernimmt sie im Gegenzug, damit beide Plattformen
exakt dieselbe Markenfarbe zeigen.

Mit dem Wechsel auf Grün entfällt `--success` als eigenes Token: das
Markengrün (135,5° im HSL-Farbkreis) und das frühere Erfolgsgrün (141,7°)
liegen nur 6° auseinander und wären nebeneinander nicht unterscheidbar. In
einem grünen System ist Erfolg die Markenfarbe — `--accent-text` steht
überall dort, wo zuvor `--success` stand. Das frei gewordene vormalige
Markenblau lebt als `--info` weiter, in derselben Rolle, die `systemBlue` in
der iOS-App für wertungsfreie Hinweise trägt.

### Hell

| Token | Wert | gegen `--bg` | gegen `--card` | Anforderung |
|---|---|---|---|---|
| `--bg` | `#f8fafc` | — | — | Seitenhintergrund |
| `--card` | `#ffffff` | — | — | Kartenfläche |
| `--text` | `#1e293b` | 13.98:1 | 14.63:1 | ≥ 7:1 |
| `--text-muted` | `#475569` | 7.24:1 | 7.58:1 | ≥ 7:1 |
| `--accent` (Fläche) | `#186029` | 7.31:1 (als Text gelesen) | 7.65:1 (als Text gelesen) | eigene Prüfung, siehe unten |
| `--accent-on` (Text auf `--accent`) | `#ffffff` | — | — | 7.65:1 auf `--accent`, ≥ 7:1 |
| `--accent-text` | `#186029` | 7.31:1 | 7.65:1 | ≥ 7:1 |
| `--accent-hover` (Fläche) | `#124a1f` | — | — | 10.37:1 unter `--accent-on`, ≥ 7:1 |
| `--info` | `#1e40af` | 8.34:1 | 8.72:1 | ≥ 7:1 |
| `--error` | `#991b1b` | 7.94:1 | 8.31:1 | ≥ 7:1 |
| `--warning` | `#93390e` | 7.09:1 | 7.41:1 | ≥ 7:1 |
| `--text-disabled` | `#475569` | 7.24:1 | 7.58:1 | ≥ 7:1 |
| `--control-line` | `#767676` | 4.34:1 | 4.54:1 | ≥ 3:1 |
| `--border` | `#e2e8f0` | dekorativ | dekorativ | ausgenommen |

`--accent` und `--accent-hover` sind Flächen, keine Textfarben — als Text
gegen `--bg`/`--card` gelesen bestünden sie im hellen Thema zwar knapp
(7.31:1 / 7.65:1 für `--accent`), im dunklen Thema aber nicht (siehe unten).
Sie werden deshalb nicht gegen die Seitenflächen geprüft, sondern dagegen,
was tatsächlich auf ihnen steht: `--accent-on` erreicht 7.65:1 auf `--accent`
und 10.37:1 auf `--accent-hover`, in beiden Themen gleich, weil beide Werte
themenunabhängig sind.

`--text-disabled` stammt aus Tempus und fehlt Ortus — ein Beispiel für die
Drift, die das gemeinsame Modul beendet. Es erhält denselben Wert wie
`--text-muted`: ein deaktiviertes Bedienelement darf nicht durch schwachen
Kontrast erkennbar gemacht werden, sondern durch `aria-disabled` und
Cursor.

### Dunkel

| Token | Wert | gegen `--card` | gegen `--bg` | Anforderung |
|---|---|---|---|---|
| `--bg` | `#0f172a` | — | — | Seitenhintergrund |
| `--card` | `#1e293b` | — | — | Kartenfläche |
| `--text` | `#f1f5f9` | 13.35:1 | 16.30:1 | ≥ 7:1 |
| `--text-muted` | `#b4c0ce` | 7.92:1 | 9.67:1 | ≥ 7:1 |
| `--accent` (Fläche, unverändert) | `#186029` | 1.91:1 (als Text gelesen) | 2.33:1 (als Text gelesen) | eigene Prüfung, siehe oben |
| `--accent-on` (Text auf `--accent`, unverändert) | `#ffffff` | — | — | 7.65:1 auf `--accent`, ≥ 7:1 |
| `--accent-text` | `#6ecb86` | 7.35:1 | 8.97:1 | ≥ 7:1 |
| `--accent-hover` (Fläche, unverändert) | `#124a1f` | — | — | 10.37:1 unter `--accent-on`, ≥ 7:1 |
| `--info` | `#93c5fd` | 8.11:1 | 9.90:1 | ≥ 7:1 |
| `--error` | `#fca5a5` | 7.71:1 | 9.41:1 | ≥ 7:1 |
| `--warning` | `#fcd34d` | 10.15:1 | 12.38:1 | ≥ 7:1 |
| `--text-disabled` | `#b4c0ce` | 7.92:1 | 9.67:1 | ≥ 7:1 |
| `--control-line` | `#8695a8` | 4.79:1 | 5.85:1 | ≥ 3:1 |
| `--border` | `#334155` | dekorativ | dekorativ | ausgenommen |

`--accent` bleibt im dunklen Thema derselbe Wert wie im hellen — als Text
gegen dunklen Grund gelesen erreicht es hier nur 1.91:1 bzw. 2.33:1, weit
unter 7:1. Das ist beabsichtigt: `--accent` ist eine Fläche, keine Textfarbe
(siehe Herkunft der Palette unten), geprüft wird sie gegen `--accent-on`,
nicht gegen `--bg`/`--card`.

Zwei Konsequenzen aus dem vorigen (blauen) System sind weiterhin
erwähnenswert:

**Sekundärtext wird dunkler.** Ortus' und Tempus' `#64748b` (4.55:1) weicht
`#475569` (7.6:1). Das ist die sichtbarste Änderung im hellen Thema.

**Grün und Orange mussten nachgezogen werden.** Ein erster Entwurf des
ursprünglichen Systems setzte das damalige Erfolgsgrün auf `#166534` und
`--warning` auf `#92400e`. Beide erreichten 7:1 gegen die weiße Karte,
verfehlten es aber gegen den Seitenhintergrund (6.81:1 und 6.78:1) — ein
Fall, den nur die Prüfung gegen beide Flächen findet. Die jetzigen Werte für
`--warning` erfüllen beide.

**Karten brauchen im Dunkeln einen Rand.** `--card: #1e293b` steht gegen
`--bg: #0f172a` nur bei 1.22:1. Der `box-shadow`, der die Karte im hellen
Thema vom Hintergrund abhebt, leistet das auf dunklem Grund nicht. Im dunklen
Thema tritt deshalb `border: 1px solid var(--border)` an die Stelle des
Schattens. Dieselbe Regel gilt für jede andere Fläche, die sich allein durch
Helligkeit abheben soll.

`--border` und `--control-line` bleiben getrennt. Die Begründung übernimmt
`tokens.css` im Wortlaut aus Hostus: WCAG 1.4.11 fordert Kontrast für die
Begrenzung, die ein **Bedienelement** erkennbar macht; dekorative Trennlinien
sind ausgenommen und dürfen zart bleiben.

## Weitere Tokens

Typografie, Abstände, Radien und Schatten folgen den Werten, die Ortus und
Tempus bereits gemeinsam verwenden — sie haben sich bewährt und ihre
Übernahme hält die Migration dieser beiden Dienste klein.

- Schriftfamilie: `-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif`
- Monospace: `'SF Mono', Monaco, ui-monospace, monospace`
- Zeilenhöhe `1.5`, Basisgröße `1rem`
- `--radius: 8px`, kleinerer Radius `4px` für eingebettete Flächen
- `--shadow: 0 1px 3px rgba(0,0,0,0.1)` (nur im hellen Thema wirksam)
- Ein Breakpoint bei `640px`, mobile-first; ein zweiter wird heute nicht
  gebraucht und wird erst ergänzt, wenn ein Dienst ihn tatsächlich braucht

## Wozu die Dienste da sind — und was daraus folgt

Die fünf Frontends sind keine fertigen Erzeugnisse, sondern Bruchstücke: Jedes
dient einer anderen Anwendergruppe als MVP, und alle sollen später in einer
einheitlichen Anwendung wiederverwendet werden.

Das kehrt die übliche Abwägung um. Ein Gestaltungsbestandteil, den heute nur
ein Dienst braucht — ein Symbol, eine Komponente —, ist keine Sonderlocke,
sondern ein Baustein der künftigen Anwendung. Wiedererkennbarkeit wiegt hier
schwerer als Sparsamkeit. Im Zweifel gehört ein gestalterischer Bestandteil
deshalb ins gemeinsame Modul, auch wenn er zunächst nur an einer Stelle
auftaucht.

Die Grenze verläuft nicht zwischen „einer braucht es" und „mehrere brauchen
es", sondern zwischen **Gestaltung und Fachlogik**:

| gehört ins Modul | bleibt beim Dienst |
|---|---|
| die Zeichnung eines Symbols | welcher Zustand welches Symbol bekommt |
| Kopf- und Fußzeile als Bauform | die Beschriftungen darin |
| Farben, Raster, Strichstärke | die Auswertung, die zu einem Wert führt |

Ein Beispiel: Das Modul liefert das Wettersymbol für Nebel. Die Aussage
„WMO-Code 45 bedeutet Nebel" ist Wetterlogik, ändert sich mit der Schnittstelle
des Datenanbieters und bleibt in Tempus.

## Seitengerüst

Ortus und Tempus führen Kopf- und Fußzeile zeichengleich — sie unterscheiden
sich in drei Wörtern:

```html
<header><h1>Ortus</h1><p>Point-in-Polygon Abfrage über Datenquellen</p></header>
<header><h1>Tempus</h1><p>Ortsdaten mit Zeitbezug</p></header>

<footer>API Dokumentation · OpenAPI Spec · Health Status
        <div class="footer-version">ortus __ORTUS_VERSION__</div></footer>
```

Dieselbe Fußzeile, dieselben drei Verweise, dieselbe Versionszeile. Alle vier
Go-Dienste bieten `/docs`, `/openapi.json` und `/health` an; Expertus hat
stattdessen andere Ziele, braucht aber dieselbe Bauform.

Das Modul liefert deshalb Kopf- und Fußzeile als Bauform: Name, Untertitel,
Verweisliste und Versionsangabe werden übergeben, das Markup kommt aus dem
Modul. Damit rückt zugleich ein Mangel aus Expertus gerade — dort beginnt der
Kopf am Fensterrand, während der Inhalt zentriert ist; Kopf, Inhalt und Fuß
teilen sich künftig dieselbe Spur.

## Symbole

Heute zeichnet jeder Dienst für sich: Ortus führt denselben Chevron dreimal in
einer Datei, Situs und Hostus haben gar keine Symbole, und Tempus benutzt
**Emoji** — `&#x2600;` für Sonne, 🌑 bis 🌘 für die Mondphasen, dazu 📊, 🌍, 💧.

Emoji sind für diesen Zweck ungeeignet:

- Sie werden vom Betriebssystem gezeichnet. Auf macOS, Windows und Android
  sehen sie verschieden aus; die Darstellung ist nicht bestimmbar.
- Sie sind mehrfarbig und gehorchen keiner Palette. Im dunklen Thema bleiben
  sie, wie sie sind.
- Sie teilen weder Strichstärke noch Raster — nebeneinander wirken sie wie
  eine Sammlung, nicht wie ein Satz.
- Bei den Mondphasen wiegt das am schwersten: 🌒 und 🌓 sind auf vielen
  Systemen kaum zu unterscheiden, und die Phase **ist** die Information.

Das Modul liefert die Symbole deshalb als Go-Funktionen in einem Unterpaket
`icons/`, je eine Funktion pro Symbol, die fertiges SVG zurückgibt.

**Bauregel.** Sie ist keine Erfindung, sondern übernimmt, was Ortus bereits
tut:

- `viewBox="0 0 24 24"`, keine festen `width`/`height` — die Größe bestimmt
  der Einsatzort über CSS
- `fill="none"`, `stroke="currentColor"`, `stroke-width="2"`,
  `stroke-linecap="round"`, `stroke-linejoin="round"`
- `currentColor` ist der Kern: ein Symbol nimmt damit die Farbe seines
  Umfelds an und folgt Thema und Zustand, ohne eine eigene Farbe zu kennen.
  Ein Symbol mit eigener Farbe stünde außerhalb der Kontrastprüfung.
- `aria-hidden="true"` und `focusable="false"` als Vorgabe: ein Symbol neben
  einer Beschriftung ist Schmuck und darf nicht doppelt vorgelesen werden.
  Steht ein Symbol allein — etwa in einem Knopf ohne Text —, trägt das
  umgebende Bedienelement die Beschriftung, nicht das Symbol.

**Eine begründete Ausnahme:** Die Mondphasen dürfen Flächen verwenden
(`fill="currentColor"`), weil bei ihnen die Aufteilung zwischen beleuchtetem
und unbeleuchtetem Teil die Aussage trägt. Eine reine Strichzeichnung könnte
zunehmenden und abnehmenden Halbmond nicht unterscheiden.

## Komponenten in `base.css`

Aufgenommen wird, was mindestens zwei Dienste brauchen:

- Reset (`box-sizing`, Margin-Reset), `body`, `.container`
- Sichtbarer Fokusring auf allen fokussierbaren Elementen
- `label`, `input`, `select`, `textarea` samt Fokuszustand und Platzhalter
- `.btn`, `.btn-secondary`, `.btn-row`, Zustände `hover`/`disabled`/`focus-visible`
- `.card`, `.card-title`
- `.tabs`, `.tab` mit `aria-selected`-Zustand
- `.badge` samt Erfolgs-, Fehler-, Warnungs- und Hinweisvariante
- Tabellen-Grundlayout und `.table-wrap` für waagerechtes Rollen
- `.spinner`, `.loading`
- `.error`, `.warning` und `.info` als Meldungsboxen
- `.sr-only`, `.skip-link`
- `@media (prefers-reduced-motion: reduce)`

Nicht aufgenommen wird Domänenspezifisches: `.weather-panel` und
`.provider-row` bleiben in Tempus, `.source-card` und `.gazetteer-block` in
Ortus, `.hits` in Situs, die Statusflächen in Hostus. Diese Komponenten
greifen auf die Tokens zu, leben aber weiter im jeweiligen Projekt.

Mindestgrößen für Bedienelemente: 44 × 44 CSS-Pixel wo möglich, 24 × 24 als
harte Untergrenze (WCAG 2.2, SC 2.5.8). Expertus' heutige Regel wird damit
zum Standard.

## Prüfung im Design-System-Repo

`contrast.go` berechnet die relative Leuchtdichte nach WCAG und das
Kontrastverhältnis zweier Farben. `tokens_test.go` liest `css/tokens.css` ein
und prüft:

- jedes Token, das keiner ausdrücklichen Ausnahme unterliegt (Flächen,
  Linien, Radien, Schrift, `--accent`, `--accent-hover`, `--accent-on`),
  gegen `--card` und `--bg`, in beiden Themen, ≥ 7:1 — die Ausnahmeliste
  steht im Testcode, alles andere ergibt sich aus der Datei selbst, ein neues
  Textfarb-Token wird also ohne weiteres Zutun erfasst. Jeder Eintrag der
  Ausnahmeliste trägt im Testcode einen Verweis auf die Prüfung, die ihn
  stattdessen abdeckt — eine Ausnahme ohne eigene gezielte Prüfung wäre
  schlicht ungeprüft.
- `--control-line` gegen `--card` und `--bg`, beide Themen, ≥ 3:1
- die Knopfbeschriftung (`--accent-on`) auf `--accent` und `--accent-hover`,
  beide Themen, ≥ 7:1 — das ist die gezielte Prüfung für die drei oben
  genannten Ausnahmen: `--accent` und `--accent-hover` sind Flächen, keine
  Textfarben, `--accent-on` ist Text auf dieser Fläche, nicht auf
  `--bg`/`--card`
- dass `tokens.css` genau einen `@media`-Block enthält und dieser
  `prefers-color-scheme: dark` lautet — die Zerlegung in helles und dunkles
  Thema setzt das voraus

`css_test.go` sucht Farbliterale außerhalb von `tokens.css` (Hex, `rgb()`/
`hsl()`, moderne Funktionen wie `oklch()`/`color-mix()`, die vollständige
CSS-Farbwortliste, prozentkodierte Hex-Werte in data-URIs) und stellt sicher,
dass `base.css` nur Token nutzt, die `tokens.css` auch definiert.

Was diese Tests **nicht** abdecken: dass eine tatsächliche Regel in
`base.css` zwei einzeln geprüfte Token auch sinnvoll kombiniert. Zwei für
sich genommen unauffällige Token können in einer neuen Regel zu einem
unlesbaren Paar werden — nachgewiesen etwa `color: var(--text-muted)` auf
`background: var(--accent)`, 1.15:1, während alle bis dahin genannten Tests
grün bleiben. Ein eigener Test in `farbpaarung_test.go` deckt den
naheliegendsten Fall davon ab: jede Regel, die `color: var(--…)` UND
`background`/`background-color: var(--…)` **im selben Block** setzt, wird
gegeneinander auf 7:1 geprüft. Die bewusste Grenze dieser Prüfung: sie
kombiniert keine Farben über mehrere Regeln hinweg — eine Farbe, die (wie
`.btn:hover`) aus einer anderen Regel geerbt wird, weil sie in der eigenen
Regel nicht neu gesetzt ist, entgeht ihr. Das verlangte den vollen
CSS-Kaskadenalgorithmus, den eine einfache Textzerlegung nicht nachbildet.
Geprüft ist also: dass jedes Token für sich gegen die Flächen besteht, und
dass jede Regel, die Text- und Flächenfarbe gemeinsam benennt, zueinander
passt — nicht, dass jede im CSS durch Kaskade und Vererbung tatsächlich
zustande kommende Kombination zueinander passt.

Die Demo-Seite unter `demo/` zeigt jede Komponente in beiden Themen und dient
der visuellen Abnahme. Zwei Tests (`demo/demo_test.go`) stellen sicher, dass
sie das echte Stylesheet über `<link rel="stylesheet">` einbindet und keinen
eigenen `<style>`-Block mitbringt — sonst ließe sich jede der oben genannten
Prüfungen durch eine lokal „reparierte" Demo-Seite unterlaufen.

## Expertus: Umstellung auf Go

Expertus wird heute von nginx ausgeliefert, mit Konfiguration aus
Umgebungsvariablen. Damit es dasselbe Go-Modul nutzen kann wie die anderen
vier, tritt ein eigener Go-Server an die Stelle von nginx. Der Server muss
genau das leisten, was `docker/nginx.conf` und `docker/entrypoint.sh` heute
tun:

| heute | künftig |
|---|---|
| `entrypoint.sh` schreibt `config.json` aus `ORTUS_BASE_URL`, `HABITATUS_BASE_URL`, `HOSTUS_BASE_URL` | Handler `/config.json` rendert dieselben drei Werte aus der Umgebung |
| CSP mit `connect-src` aus denselben Variablen gebaut | Middleware baut denselben CSP-String |
| `X-Content-Type-Options`, `Referrer-Policy`, `Cache-Control`, `etag` je Location | dieselben Header in einem Middleware-Block |
| `try_files $uri $uri/ /index.html` | `http.FileServer` über `embed.FS`, Fallback auf `index.html` |

Mit berührt:

- `Dockerfile` — mehrstufiger Go-Build statt nginx-Image
- `docker/nginx.conf` und `docker/entrypoint.sh` — entfallen
- `playwright.config.js` — `webServer.command` startet das Binary statt
  `python3 -m http.server`
- `Makefile` — `serve` und `smoke` entsprechend
- `.github/workflows/` — Go-Toolchain, `go test`

Der Grundsatz, den der Dockerfile heute festhält — das Ausgelieferte ist
identisch mit dem Repository-Inhalt — bleibt gewahrt: `go:embed` legt die
Dateien unverändert ins Binary, es findet weiterhin keine Transformation
statt. Der Kommentar wird entsprechend neu formuliert.

## Migrationsreihenfolge

Jeder Schritt ist für sich abgeschlossen und überprüfbar.

1. **Design-System-Repo** anlegen: Tokens, `base.css`, `contrast_test.go`,
   Demo-Seite. Ohne Konsumenten prüfbar, visuelle Abnahme an der Demo-Seite.
2. **Expertus** — zuerst der Go-Server, dann die Tokens. Expertus geht
   voran, weil es noch nicht bei Endanwendern ist: Fehler im frischen
   System treffen dort niemanden, und die vier produktiven Frontends
   bleiben unangetastet, bis sich das Modul bewährt hat.
3. **Tempus** — kürzester Weg unter den produktiven Diensten, weil das CSS
   bereits in einer eigenen `index.html` per `go:embed` liegt.
4. **Ortus** — dabei die Go-String-Konstante in `frontend.go` auf
   `go:embed index.html` umstellen. Die Datei schrumpft um rund 700 Zeilen.
5. **Hostus** — behält seine Statusflächen als lokale Ergänzung.
6. **Situs** — der aufwendigste Schritt. Nicht nur die Palette wechselt: der
   Explorer benutzt mit `.panel`, `.globals`, `.row` ein eigenes
   Klassenvokabular, das sich nicht eins zu eins auf `.card`, `.form-group`
   und `.btn-row` abbilden lässt. Hier wird Markup umgeschrieben, nicht nur
   CSS ersetzt. Sein Dark Mode kommt künftig aus den Tokens statt aus
   eigener Media Query.

Diese Reihenfolge hat einen Preis, der die Demo-Seite aufwertet: Expertus
bringt heute 75 Zeilen CSS mit und kennt weder Karten noch Tabs, Badges oder
Tabellen. Als erster Konsument übt es das Modul also nur zu einem kleinen
Teil aus. Die Demo-Seite in `demo/` ist damit für die Mehrzahl der
Komponenten die einzige Prüfstelle, bevor sie produktiven Code berühren —
sie muss jede Komponente in beiden Themen zeigen und wird mit derselben
Sorgfalt abgenommen wie ein Dienst.

## Absicherung

Die vorhandenen Tests der Projekte sind das Sicherheitsnetz der Migration:
Hostus' Mutationstest auf seine Assets, Expertus' axe-core-Lauf über
Playwright, Ortus' `frontend_test.go`. Sie bleiben unverändert bestehen und
müssen nach jedem Migrationsschritt grün sein.

Ergänzt wird je Dienst ein Test, der prüft, dass die ausgelieferte Seite das
Design-System-CSS tatsächlich einbindet — sonst fällt ein Dienst
unbemerkt auf lokale Reste zurück.

## Bewusst nicht enthalten

- Kein Umschalter für das Thema und keine gespeicherte Wahl; allein die
  Systemeinstellung entscheidet.
- Keine dienstspezifischen Akzentfarben.
- Kein JavaScript für Verhalten, das HTML selbst mitbringt. Das Akkordeon
  entsteht aus `<details>` und `<summary>`: auf- und zuklappbar,
  tastaturbedienbar, von Screenreadern richtig angesagt — ohne eine Zeile
  Skript. Das Modul liefert dafür nur die Gestaltung. Dasselbe gilt für die
  Reiter, deren Zustand am `aria-selected` des Dienstes hängt.
- **Eine Ausnahme: die Combobox mit Vorschlagsliste.** Sie lässt sich nicht
  ohne Skript bauen, und ihr Verhalten ist der schwierige Teil — Pfeiltasten,
  `aria-activedescendant`, Entprellung, Abbruch überholter Anfragen. Genau
  dort entstehen die Fehler, und genau dort lohnt es sich, sie einmal richtig
  zu lösen statt in jedem Dienst neu. Das Modul liefert deshalb neben `CSS()`
  auch `JS()`.

  Die Grenze bleibt dieselbe wie bei den Symbolen: Das Modul liefert das
  Verhalten der Bedienform, nicht die Fachlogik. Woher die Vorschläge kommen
  und was ein Treffer bedeutet, entscheidet der Dienst.
- Kein Bundler, kein npm, kein Git-Submodul.
