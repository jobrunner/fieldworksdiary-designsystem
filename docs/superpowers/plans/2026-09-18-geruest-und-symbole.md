# Seitengerüst und Symbole — Implementierungsplan

> **Für agentische Bearbeiter:** ERFORDERLICHE UNTER-SKILL: `superpowers:subagent-driven-development` (empfohlen) oder `superpowers:executing-plans`, um diesen Plan Aufgabe für Aufgabe umzusetzen. Die Schritte nutzen Kontrollkästchen (`- [ ]`) zur Nachverfolgung.

**Ziel:** Das Design-System um Kopf- und Fußzeile sowie rund 38 Symbole erweitern, damit die fünf Dienste — und die daraus später entstehende einheitliche Anwendung — dieselbe Bildsprache teilen.

**Architektur:** Ein Unterpaket `icons/` mit je einer Go-Funktion pro Symbol, die fertiges SVG liefert. Alle Symbole auf 24er-Raster mit `currentColor`, sodass sie Thema und Zustand ihres Umfelds übernehmen, ohne eine eigene Farbe zu kennen. Kopf- und Fußzeile kommen als Bauform aus dem Wurzelpaket: Beschriftungen werden übergeben, das Markup stammt aus dem Modul.

**Tech-Stack:** Go 1.26, `html/template` für das Gerüst, reines SVG von Hand. Keine externen Abhängigkeiten.

**Spec:** `docs/design.md`, Abschnitte „Wozu die Dienste da sind", „Seitengerüst" und „Symbole"

## Globale Randbedingungen

- Modulpfad `github.com/jobrunner/fieldworksdiary-designsystem`, Go `1.26.0`, **keine** externen Abhängigkeiten.
- **Bauregel für jedes Symbol:** `viewBox="0 0 24 24"`, keine festen `width`/`height`, `fill="none"`, `stroke="currentColor"`, `stroke-width="2"`, `stroke-linecap="round"`, `stroke-linejoin="round"`, `aria-hidden="true"`, `focusable="false"`.
- **Einzige Ausnahme:** Mondphasen dürfen `fill="currentColor"` verwenden — bei ihnen trägt die Flächenaufteilung die Aussage. Jede Ausnahme ist im Code zu begründen.
- Kein Symbol enthält eine Farbe im Klartext. `currentColor` ist verbindlich; ein Symbol mit eigener Farbe stünde außerhalb der Kontrastprüfung.
- Das Modul liefert **Zeichnungen**, keine Fachlogik. Die Zuordnung „WMO-Code 45 → Nebel" bleibt in Tempus.
- Alle Kommentare auf Deutsch und begründend. Commits auf Deutsch mit Conventional-Commits-Präfix.
- Nach jeder Aufgabe: `go build ./... && go vet ./... && go test ./... && gofmt -l .` sauber.
- **Die Palette ist seit der Erstellung dieses Plans von Blau auf Grün gewechselt.** `--accent` ist jetzt eine **Fläche** (`#186029`, in beiden Themen gleich, trägt weiße Schrift über `--accent-on`); als Textfarbe dient `--accent-text` (`#186029` hell / `#6ECB86` dunkel). Ein eigenes `--success` gibt es nicht mehr — Erfolg ist die Markenfarbe. Neu hinzugekommen ist `--info`. Wo dieser Plan eine Akzentfarbe für **Text, Linien oder Umrisse** vorsieht, ist `--accent-text` gemeint.

---

## Dateistruktur

| Datei | Verantwortung |
|---|---|
| `icons/icons.go` | Bauform eines Symbols und die gemeinsame Hülle — sonst nichts |
| `icons/bedienung.go` | die 12 Bedien-Symbole |
| `icons/wetter.go` | die 12 Wettersymbole |
| `icons/himmel.go` | 8 Mondphasen, Sonnenauf- und -untergang |
| `icons/messwerte.go` | Tropfen, Globus, Diagramm, Thermometer |
| `icons/icons_test.go` | prüft die Bauregel für **jedes** Symbol |
| `layout.go` | Kopf- und Fußzeile als Bauform |
| `layout_test.go` | dazugehörige Prüfungen |
| `demo/index.html` | zeigt Gerüst und alle Symbole |

Die Symbole liegen nach Thema in getrennten Dateien, damit keine davon
unübersichtlich wird und ein Beitrag erkennbar bleibt. `icons.go` trägt die
gemeinsame Hülle, damit die Bauregel an genau einer Stelle steht.

---

### Aufgabe 1: Bauform eines Symbols

**Dateien:**
- Anlegen: `icons/icons.go`, `icons/icons_test.go`

**Schnittstellen:**
- Liefert: `type Icon string`, `func (i Icon) String() string`, die interne Hülle `huelle(inhalt string, gefuellt bool) Icon` — alle folgenden Aufgaben bauen darauf.

- [ ] **Schritt 1: Den fehlschlagenden Test schreiben**

`icons/icons_test.go`:

```go
package icons

import "testing"

func TestHuelleFolgtDerBauregel(t *testing.T) {
	got := string(huelle(`<path d="M6 9l6 6 6-6"/>`, false))

	for _, pflicht := range []string{
		`viewBox="0 0 24 24"`,
		`fill="none"`,
		`stroke="currentColor"`,
		`stroke-width="2"`,
		`stroke-linecap="round"`,
		`stroke-linejoin="round"`,
		`aria-hidden="true"`,
		`focusable="false"`,
		`<path d="M6 9l6 6 6-6"/>`,
	} {
		if !contains(got, pflicht) {
			t.Errorf("die Hülle enthält %q nicht:\n%s", pflicht, got)
		}
	}
	// Feste Maße würden die Größe am Einsatzort festnageln; sie soll über
	// CSS bestimmt werden.
	for _, verboten := range []string{`width="`, `height="`} {
		if contains(got, verboten) {
			t.Errorf("die Hülle enthält %q — die Größe gehört ins CSS", verboten)
		}
	}
}

func TestGefuellteHuelleFuerFlaechensymbole(t *testing.T) {
	// Die Mondphasen sind die einzige Ausnahme von der Strichzeichnung:
	// bei ihnen trägt die Flächenaufteilung die Aussage.
	got := string(huelle(`<circle cx="12" cy="12" r="9"/>`, true))
	if !contains(got, `fill="currentColor"`) {
		t.Errorf("gefüllte Hülle nutzt kein fill=currentColor:\n%s", got)
	}
	if contains(got, `fill="none"`) {
		t.Errorf("gefüllte Hülle enthält noch fill=none:\n%s", got)
	}
}

func contains(h, n string) bool {
	return len(n) <= len(h) && indexOf(h, n) >= 0
}

func indexOf(h, n string) int {
	for i := 0; i+len(n) <= len(h); i++ {
		if h[i:i+len(n)] == n {
			return i
		}
	}
	return -1
}
```

- [ ] **Schritt 2: Test laufen lassen, Fehlschlag bestätigen**

Ausführen: `go test ./icons/ -v`
Erwartet: Übersetzungsfehler — `undefined: huelle`.

- [ ] **Schritt 3: Die Umsetzung schreiben**

`icons/icons.go`:

```go
// Package icons liefert die Symbole der fieldworksdiary-Dienste als fertiges
// SVG. Je Symbol eine Funktion; die gemeinsame Hülle steht an genau einer
// Stelle, damit die Bauregel nicht 38-mal wiederholt und dabei verletzt wird.
//
// Das Modul liefert Zeichnungen, keine Fachlogik: welcher Zustand welches
// Symbol bekommt — etwa welcher WMO-Wettercode Nebel bedeutet — entscheidet
// der Dienst.
package icons

import "strings"

// Icon ist fertiges SVG-Markup. Eigener Typ statt string, damit an der
// Aufrufstelle erkennbar bleibt, dass hier Markup eingesetzt wird und keine
// Beschriftung.
type Icon string

func (i Icon) String() string { return string(i) }

// huelle setzt ein Symbol aus seinem Inhalt zusammen.
//
// currentColor ist der Kern der Bauregel: ein Symbol nimmt damit die Farbe
// seines Umfelds an und folgt Thema und Zustand, ohne eine eigene Farbe zu
// kennen. Ein Symbol mit eigener Farbe stünde außerhalb der Kontrastprüfung
// des Moduls.
//
// Ohne width und height: die Größe bestimmt der Einsatzort über CSS. Feste
// Maße würden ein Symbol im Knopf und in der Tabelle gleich groß machen.
//
// aria-hidden, weil ein Symbol neben einer Beschriftung Schmuck ist und
// sonst doppelt vorgelesen würde. Steht es allein — etwa in einem Knopf ohne
// Text —, trägt das umgebende Bedienelement die Beschriftung.
//
// gefuellt kehrt die Zeichenart um: Flächen statt Striche. Das ist allein
// für die Mondphasen vorgesehen, bei denen die Aufteilung zwischen
// beleuchtetem und unbeleuchtetem Teil die Aussage trägt — eine reine
// Strichzeichnung könnte zunehmenden und abnehmenden Halbmond nicht
// unterscheiden.
func huelle(inhalt string, gefuellt bool) Icon {
	var b strings.Builder
	b.WriteString(`<svg viewBox="0 0 24 24" `)
	if gefuellt {
		b.WriteString(`fill="currentColor" stroke="none" `)
	} else {
		b.WriteString(`fill="none" stroke="currentColor" stroke-width="2" `)
		b.WriteString(`stroke-linecap="round" stroke-linejoin="round" `)
	}
	b.WriteString(`aria-hidden="true" focusable="false">`)
	b.WriteString(inhalt)
	b.WriteString(`</svg>`)
	return Icon(b.String())
}

// strich baut ein Symbol in der Regelform: Strichzeichnung.
func strich(inhalt string) Icon { return huelle(inhalt, false) }

// flaeche baut ein Symbol als Flächenzeichnung. Nur für die Mondphasen —
// siehe die Begründung an huelle.
func flaeche(inhalt string) Icon { return huelle(inhalt, true) }
```

- [ ] **Schritt 4: Test laufen lassen, Erfolg bestätigen**

Ausführen: `go test ./icons/ -v`
Erwartet: PASS.

- [ ] **Schritt 5: Einchecken**

```bash
git add icons/
git commit -m "feat: Bauform für Symbole

Die Bauregel steht an einer Stelle statt 38-mal wiederholt zu werden.
currentColor ist verbindlich: ein Symbol mit eigener Farbe stünde
außerhalb der Kontrastprüfung des Moduls."
```

---

### Aufgabe 2: Bedien-Symbole

**Dateien:**
- Anlegen: `icons/bedienung.go`
- Ändern: `icons/icons_test.go` (Regelprüfung über alle Symbole)

**Schnittstellen:**
- Nutzt: `strich` aus Aufgabe 1.
- Liefert: `Standort()`, `ChevronUnten()`, `Kalender()`, `Schliessen()`, `Suche()`, `Herunterladen()`, `Kopieren()`, `Haken()`, `Warnung()`, `Information()`, `Fehler()`, `Menue()` — Expertus und die vier Go-Dienste rufen sie auf.
- Liefert außerdem: `Alle() map[string]Icon` — die Sammlung aller Symbole, damit Prüfung und Demo-Seite nicht händisch gepflegt werden müssen.

- [ ] **Schritt 1: Den fehlschlagenden Test schreiben**

Die entscheidende Prüfung: **jedes** Symbol folgt der Bauregel. Sie leitet sich aus `Alle()` ab, nicht aus einer Liste im Testcode — sonst fehlt das nächste Symbol still.

An `icons/icons_test.go` anhängen:

```go
func TestJedesSymbolFolgtDerBauregel(t *testing.T) {
	if len(Alle()) == 0 {
		t.Fatal("Alle() ist leer — dann prüft dieser Test nichts")
	}
	for name, icon := range Alle() {
		s := string(icon)
		t.Run(name, func(t *testing.T) {
			for _, pflicht := range []string{
				`viewBox="0 0 24 24"`,
				`aria-hidden="true"`,
				`focusable="false"`,
				`currentColor`,
			} {
				if !contains(s, pflicht) {
					t.Errorf("%s enthält %q nicht", name, pflicht)
				}
			}
			for _, verboten := range []string{`width="`, `height="`, `#`, `rgb(`, `hsl(`} {
				if contains(s, verboten) {
					t.Errorf("%s enthält %q — Größe gehört ins CSS, Farbe kommt aus currentColor", name, verboten)
				}
			}
			if !contains(s, `<svg`) || !contains(s, `</svg>`) {
				t.Errorf("%s ist kein vollständiges SVG", name)
			}
		})
	}
}

func TestAlleEnthaeltDieBediensymbole(t *testing.T) {
	// Diese Namen sind eine Zusage an die Dienste: sie rufen sie auf.
	for _, name := range []string{
		"standort", "chevron-unten", "kalender", "schliessen", "suche",
		"herunterladen", "kopieren", "haken", "warnung", "information",
		"fehler", "menue",
	} {
		if _, da := Alle()[name]; !da {
			t.Errorf("Alle() kennt %q nicht", name)
		}
	}
}
```

- [ ] **Schritt 2: Test laufen lassen, Fehlschlag bestätigen**

Ausführen: `go test ./icons/ -run TestJedesSymbol -v`
Erwartet: Übersetzungsfehler — `undefined: Alle`.

- [ ] **Schritt 3: Die Symbole zeichnen**

`icons/bedienung.go`. Die Pfade stammen, wo vorhanden, aus dem Bestand — der
Chevron und das Standort-Symbol sind Ortus entnommen, damit dessen Aussehen
erhalten bleibt:

```go
package icons

// Bedien-Symbole: was in jedem Dienst vorkommt, unabhängig vom Fach.

// Standort — für den Knopf, der die aktuelle Position bestimmt.
// Übernommen aus Ortus, wo dieses Symbol bereits in Gebrauch ist.
func Standort() Icon {
	return strich(`<circle cx="12" cy="12" r="3"/><path d="M12 2v4m0 12v4M2 12h4m12 0h4"/>`)
}

// ChevronUnten — Aufklappen. In Ortus steht derselbe Pfad dreimal in einer
// Datei; das ist der Grund, warum es dieses Paket gibt.
func ChevronUnten() Icon {
	return strich(`<path d="M6 9l6 6 6-6"/>`)
}

// Kalender — Datums- und Zeitauswahl. Übernommen aus Tempus.
func Kalender() Icon {
	return strich(`<rect x="3" y="4" width="18" height="18" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/>`)
}

func Schliessen() Icon {
	return strich(`<path d="M18 6L6 18M6 6l12 12"/>`)
}

func Suche() Icon {
	return strich(`<circle cx="11" cy="11" r="8"/><path d="M21 21l-4.35-4.35"/>`)
}

func Herunterladen() Icon {
	return strich(`<path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3"/>`)
}

func Kopieren() Icon {
	return strich(`<rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>`)
}

func Haken() Icon {
	return strich(`<path d="M20 6L9 17l-5-5"/>`)
}

func Warnung() Icon {
	return strich(`<path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><path d="M12 9v4M12 17h.01"/>`)
}

func Information() Icon {
	return strich(`<circle cx="12" cy="12" r="10"/><path d="M12 16v-4M12 8h.01"/>`)
}

func Fehler() Icon {
	return strich(`<circle cx="12" cy="12" r="10"/><path d="M15 9l-6 6M9 9l6 6"/>`)
}

func Menue() Icon {
	return strich(`<path d="M3 12h18M3 6h18M3 18h18"/>`)
}

// bediensymbole ist die Sammlung dieser Datei. Alle() setzt daraus und aus
// den übrigen Dateien die Gesamtmenge zusammen; so muss niemand beim
// Hinzufügen eines Symbols an eine zweite Stelle denken.
func bediensymbole() map[string]Icon {
	return map[string]Icon{
		"standort":      Standort(),
		"chevron-unten": ChevronUnten(),
		"kalender":      Kalender(),
		"schliessen":    Schliessen(),
		"suche":         Suche(),
		"herunterladen": Herunterladen(),
		"kopieren":      Kopieren(),
		"haken":         Haken(),
		"warnung":       Warnung(),
		"information":   Information(),
		"fehler":        Fehler(),
		"menue":         Menue(),
	}
}

// Alle liefert jedes Symbol des Pakets unter seinem Namen. Prüfung und
// Referenzseite leiten sich daraus ab, damit ein neues Symbol nicht still
// ungeprüft und ungezeigt bleibt.
func Alle() map[string]Icon {
	out := map[string]Icon{}
	for _, teil := range []map[string]Icon{bediensymbole()} {
		for k, v := range teil {
			out[k] = v
		}
	}
	return out
}
```

- [ ] **Schritt 4: Test laufen lassen, Erfolg bestätigen**

Ausführen: `go test ./icons/ -v`
Erwartet: PASS, 12 Unterfälle in `TestJedesSymbolFolgtDerBauregel`.

- [ ] **Schritt 5: Einchecken**

```bash
git add icons/
git commit -m "feat: zwölf Bedien-Symbole

Chevron und Standort aus Ortus übernommen, wo der Chevron bisher dreimal
in derselben Datei stand. Die Regelprüfung leitet sich aus Alle() ab, damit
ein neues Symbol nicht still ungeprüft bleibt."
```

---

### Aufgabe 3: Wettersymbole

**Dateien:**
- Anlegen: `icons/wetter.go`
- Ändern: `icons/bedienung.go` (nur die Zeile in `Alle()`)
- Test: `icons/wetter_test.go`

**Schnittstellen:**
- Nutzt: `strich` aus Aufgabe 1, `Alle()` aus Aufgabe 2.
- Liefert: `WetterKlarTag()`, `WetterKlarNacht()`, `WetterLeichtBewoelktTag()`, `WetterLeichtBewoelktNacht()`, `WetterBewoelkt()`, `WetterNebel()`, `WetterRegen()`, `WetterSchauer()`, `WetterSchnee()`, `WetterSchneeschauer()`, `WetterGewitter()`, `WetterHagel()`.

- [ ] **Schritt 1: Den fehlschlagenden Test schreiben**

`icons/wetter_test.go`:

```go
package icons

import "testing"

func TestWettersymboleSindVollstaendig(t *testing.T) {
	// Die Namen sind eine Zusage: Tempus bildet seine WMO-Codes darauf ab.
	// Die Zuordnung selbst bleibt in Tempus — hier liegt nur die Zeichnung.
	for _, name := range []string{
		"wetter-klar-tag", "wetter-klar-nacht",
		"wetter-leicht-bewoelkt-tag", "wetter-leicht-bewoelkt-nacht",
		"wetter-bewoelkt", "wetter-nebel", "wetter-regen", "wetter-schauer",
		"wetter-schnee", "wetter-schneeschauer", "wetter-gewitter", "wetter-hagel",
	} {
		if _, da := Alle()[name]; !da {
			t.Errorf("Alle() kennt %q nicht", name)
		}
	}
}

func TestTagUndNachtUnterscheidenSich(t *testing.T) {
	// Sonne und Mond müssen verschieden gezeichnet sein — sonst trägt die
	// Unterscheidung zwischen Tag und Nacht keine Information.
	paare := [][2]string{
		{"wetter-klar-tag", "wetter-klar-nacht"},
		{"wetter-leicht-bewoelkt-tag", "wetter-leicht-bewoelkt-nacht"},
	}
	for _, p := range paare {
		if Alle()[p[0]] == Alle()[p[1]] {
			t.Errorf("%s und %s sind identisch gezeichnet", p[0], p[1])
		}
	}
}
```

- [ ] **Schritt 2: Test laufen lassen, Fehlschlag bestätigen**

Ausführen: `go test ./icons/ -run TestWetter -v`
Erwartet: FAIL — `Alle() kennt "wetter-klar-tag" nicht`.

- [ ] **Schritt 3: Die Symbole zeichnen**

`icons/wetter.go`:

```go
package icons

// Wettersymbole. Das Paket liefert die Zeichnungen; welcher WMO-Code welches
// Symbol bekommt, entscheidet der Dienst — diese Zuordnung ist Wetterlogik
// und ändert sich mit der Schnittstelle des Datenanbieters.

func WetterKlarTag() Icon {
	return strich(`<circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>`)
}

func WetterKlarNacht() Icon {
	return strich(`<path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>`)
}

func WetterLeichtBewoelktTag() Icon {
	return strich(`<path d="M12 4V2M5.6 5.6L4.2 4.2M4 12H2M18.4 5.6l1.4-1.4"/><path d="M8.5 12a3.5 3.5 0 015.9-2.5"/><path d="M17 20H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 20z"/>`)
}

func WetterLeichtBewoelktNacht() Icon {
	return strich(`<path d="M16 7a5 5 0 01-4.5-7 7 7 0 106.9 8.6"/><path d="M17 20H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 20z"/>`)
}

func WetterBewoelkt() Icon {
	return strich(`<path d="M18 17H7a5 5 0 010-10 6 6 0 0111.5-1.5A4.5 4.5 0 0118 17z"/>`)
}

func WetterNebel() Icon {
	return strich(`<path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M4 17h16M7 21h13"/>`)
}

func WetterRegen() Icon {
	return strich(`<path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M8 17v3M12 17v4M16 17v3"/>`)
}

func WetterSchauer() Icon {
	return strich(`<path d="M12 5V3M6.5 6.5L5 5M20 13h-1"/><path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M9 17v3M15 17v3"/>`)
}

func WetterSchnee() Icon {
	return strich(`<path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M8 18h.01M12 20h.01M16 18h.01M10 21h.01M14 17h.01"/>`)
}

func WetterSchneeschauer() Icon {
	return strich(`<path d="M12 5V3M6.5 6.5L5 5"/><path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M9 18h.01M13 20h.01M16 18h.01"/>`)
}

func WetterGewitter() Icon {
	return strich(`<path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M13 16l-3 5h4l-3 4"/>`)
}

func WetterHagel() Icon {
	return strich(`<path d="M17 13H7a4 4 0 010-8 5 5 0 019.5-1A3.5 3.5 0 0117 13z"/><path d="M8 17v1M12 17v1M16 17v1M10 20v1M14 20v1"/>`)
}

func wettersymbole() map[string]Icon {
	return map[string]Icon{
		"wetter-klar-tag":              WetterKlarTag(),
		"wetter-klar-nacht":            WetterKlarNacht(),
		"wetter-leicht-bewoelkt-tag":   WetterLeichtBewoelktTag(),
		"wetter-leicht-bewoelkt-nacht": WetterLeichtBewoelktNacht(),
		"wetter-bewoelkt":              WetterBewoelkt(),
		"wetter-nebel":                 WetterNebel(),
		"wetter-regen":                 WetterRegen(),
		"wetter-schauer":               WetterSchauer(),
		"wetter-schnee":                WetterSchnee(),
		"wetter-schneeschauer":         WetterSchneeschauer(),
		"wetter-gewitter":              WetterGewitter(),
		"wetter-hagel":                 WetterHagel(),
	}
}
```

In `icons/bedienung.go` die Sammlung in `Alle()` erweitern:

```go
	for _, teil := range []map[string]Icon{bediensymbole(), wettersymbole()} {
```

- [ ] **Schritt 4: Test laufen lassen, Erfolg bestätigen**

Ausführen: `go test ./icons/ -v`
Erwartet: PASS, 24 Unterfälle in der Regelprüfung.

- [ ] **Schritt 5: Einchecken**

```bash
git add icons/
git commit -m "feat: zwölf Wettersymbole

Ersetzen die Emoji in Tempus, die das Betriebssystem zeichnete und die
weder Palette noch dunkles Thema kannten. Die Zuordnung der WMO-Codes
bleibt in Tempus — hier liegt nur die Zeichnung."
```

---

### Aufgabe 4: Mondphasen und Himmelssymbole

**Dateien:**
- Anlegen: `icons/himmel.go`, `icons/himmel_test.go`
- Ändern: `icons/bedienung.go` (Sammlung in `Alle()`)

**Schnittstellen:**
- Nutzt: `strich` und `flaeche` aus Aufgabe 1.
- Liefert: `MondNeu()`, `MondZunehmendeSichel()`, `MondErstesViertel()`, `MondZunehmendGibbous()`, `MondVoll()`, `MondAbnehmendGibbous()`, `MondLetztesViertel()`, `MondAbnehmendeSichel()`, `SonnenAufgang()`, `SonnenUntergang()`.

Dies ist die Aufgabe, für die der Anstoß kam: Die Mondphasen-Emoji sind auf
vielen Systemen kaum zu unterscheiden — und die Phase **ist** die Information.

- [ ] **Schritt 1: Den fehlschlagenden Test schreiben**

`icons/himmel_test.go`:

```go
package icons

import "testing"

func TestAlleAchtMondphasen(t *testing.T) {
	for _, name := range []string{
		"mond-neu", "mond-zunehmende-sichel", "mond-erstes-viertel",
		"mond-zunehmend-gibbous", "mond-voll", "mond-abnehmend-gibbous",
		"mond-letztes-viertel", "mond-abnehmende-sichel",
	} {
		if _, da := Alle()[name]; !da {
			t.Errorf("Alle() kennt %q nicht", name)
		}
	}
}

func TestJedeMondphaseIstEigenstaendigGezeichnet(t *testing.T) {
	// Der Anlass für dieses Paket: als Emoji waren die Phasen auf vielen
	// Systemen nicht auseinanderzuhalten. Zwei gleich gezeichnete Phasen
	// wären derselbe Mangel in neuer Form.
	phasen := []string{
		"mond-neu", "mond-zunehmende-sichel", "mond-erstes-viertel",
		"mond-zunehmend-gibbous", "mond-voll", "mond-abnehmend-gibbous",
		"mond-letztes-viertel", "mond-abnehmende-sichel",
	}
	gesehen := map[string]string{}
	for _, name := range phasen {
		s := string(Alle()[name])
		if vorher, da := gesehen[s]; da {
			t.Errorf("%s ist identisch gezeichnet wie %s", name, vorher)
		}
		gesehen[s] = name
	}
}

func TestMondphasenDuerfenFlaechenNutzen(t *testing.T) {
	// Die einzige zugelassene Ausnahme von der Strichzeichnung. Geprüft
	// wird, dass sie auch wirklich nur hier auftritt.
	for name, icon := range Alle() {
		gefuellt := contains(string(icon), `fill="currentColor"`)
		istMond := len(name) >= 5 && name[:5] == "mond-"
		if gefuellt && !istMond {
			t.Errorf("%s nutzt Flächen, ist aber keine Mondphase — die Regelform ist die Strichzeichnung", name)
		}
	}
}

func TestSonnenAufUndUntergangUnterscheidenSich(t *testing.T) {
	if Alle()["sonnenaufgang"] == Alle()["sonnenuntergang"] {
		t.Error("Sonnenauf- und -untergang sind identisch gezeichnet")
	}
}
```

- [ ] **Schritt 2: Test laufen lassen, Fehlschlag bestätigen**

Ausführen: `go test ./icons/ -run TestMond -v`
Erwartet: FAIL — `Alle() kennt "mond-neu" nicht`.

- [ ] **Schritt 3: Die Symbole zeichnen**

`icons/himmel.go`. Die Phasen entstehen aus einem Kreis und einer
zweiten Form, die den unbeleuchteten Teil abdeckt:

```go
package icons

// Mondphasen und Sonnenstände.
//
// Die Mondphasen sind die einzige Ausnahme von der Strichzeichnung: bei ihnen
// trägt die Aufteilung zwischen beleuchtetem und unbeleuchtetem Teil die
// Aussage. Eine reine Strichzeichnung könnte zunehmenden und abnehmenden
// Halbmond nicht unterscheiden — und genau diese Unterscheidbarkeit war der
// Anlass, die vorherigen Emoji zu ersetzen.

// MondNeu — unbeleuchtet: nur der Umriss.
func MondNeu() Icon {
	return strich(`<circle cx="12" cy="12" r="9"/>`)
}

// MondZunehmendeSichel — schmale Sichel rechts.
func MondZunehmendeSichel() Icon {
	return flaeche(`<path d="M12 3a9 9 0 000 18 9 9 0 010-18z" opacity="0"/><circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 3a9 9 0 010 18 6 6 0 000-18z"/>`)
}

// MondErstesViertel — rechte Hälfte beleuchtet.
func MondErstesViertel() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 3a9 9 0 010 18z"/>`)
}

// MondZunehmendGibbous — mehr als halb, rechts.
func MondZunehmendGibbous() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 3a9 9 0 010 18 6 6 0 010-18z"/>`)
}

// MondVoll — vollständig beleuchtet.
func MondVoll() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9"/>`)
}

// MondAbnehmendGibbous — mehr als halb, links.
func MondAbnehmendGibbous() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 3a9 9 0 000 18 6 6 0 000-18z"/>`)
}

// MondLetztesViertel — linke Hälfte beleuchtet.
func MondLetztesViertel() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 3a9 9 0 000 18z"/>`)
}

// MondAbnehmendeSichel — schmale Sichel links.
func MondAbnehmendeSichel() Icon {
	return flaeche(`<circle cx="12" cy="12" r="9" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 3a9 9 0 000 18 6 6 0 010-18z"/>`)
}

// SonnenAufgang — Sonne über dem Horizont, Pfeil nach oben.
func SonnenAufgang() Icon {
	return strich(`<path d="M17 18a5 5 0 00-10 0"/><path d="M12 2v6M4.22 10.22l1.42 1.42M1 18h2M21 18h2M18.36 11.64l1.42-1.42M23 22H1"/><path d="M8 6l4-4 4 4"/>`)
}

// SonnenUntergang — Sonne über dem Horizont, Pfeil nach unten.
func SonnenUntergang() Icon {
	return strich(`<path d="M17 18a5 5 0 00-10 0"/><path d="M12 8V2M4.22 10.22l1.42 1.42M1 18h2M21 18h2M18.36 11.64l1.42-1.42M23 22H1"/><path d="M16 4l-4 4-4-4"/>`)
}

func himmelssymbole() map[string]Icon {
	return map[string]Icon{
		"mond-neu":                MondNeu(),
		"mond-zunehmende-sichel":  MondZunehmendeSichel(),
		"mond-erstes-viertel":     MondErstesViertel(),
		"mond-zunehmend-gibbous":  MondZunehmendGibbous(),
		"mond-voll":               MondVoll(),
		"mond-abnehmend-gibbous":  MondAbnehmendGibbous(),
		"mond-letztes-viertel":    MondLetztesViertel(),
		"mond-abnehmende-sichel":  MondAbnehmendeSichel(),
		"sonnenaufgang":           SonnenAufgang(),
		"sonnenuntergang":         SonnenUntergang(),
	}
}
```

In `Alle()` ergänzen: `bediensymbole(), wettersymbole(), himmelssymbole()`.

**Hinweis zur Umsetzung:** Die oben angegebenen Pfade sind ein Ausgangspunkt,
kein Dogma. Entscheidend ist das Ergebnis: Die acht Phasen müssen bei 24 px
Kantenlänge **auf einen Blick** unterscheidbar sein, insbesondere zunehmend
gegen abnehmend. Prüfe das, indem du die Symbole in die Referenzseite
einsetzt (Aufgabe 6) und ansiehst. Weicht deine Zeichnung ab, begründe sie im
Kommentar.

- [ ] **Schritt 4: Test laufen lassen, Erfolg bestätigen**

Ausführen: `go test ./icons/ -v`
Erwartet: PASS, 34 Unterfälle in der Regelprüfung.

Beachte: `TestJedesSymbolFolgtDerBauregel` verbietet `#`, `rgb(` und `hsl(`.
Die Mondphasen enthalten `fill="none" stroke="currentColor"` innerhalb der
gefüllten Hülle — das ist zulässig, weil keine Farbe im Klartext auftritt.
Schlägt die Regelprüfung dennoch an, liegt es an einem Attribut, nicht an der
Ausnahme; korrigiere dann das Symbol, nicht die Prüfung.

- [ ] **Schritt 5: Einchecken**

```bash
git add icons/
git commit -m "feat: acht Mondphasen und die Sonnenstände

Die Phasen waren zuvor Emoji und auf vielen Systemen nicht auseinander-
zuhalten — dabei ist die Phase die eigentliche Information. Sie sind die
einzige zugelassene Ausnahme von der Strichzeichnung; ein Test hält fest,
dass die Ausnahme nur hier auftritt."
```

---

### Aufgabe 5: Messwert-Symbole und Kopf-/Fußzeile

**Dateien:**
- Anlegen: `icons/messwerte.go`, `layout.go`, `layout_test.go`
- Ändern: `icons/bedienung.go` (Sammlung), `css/base.css` (Regeln fürs Gerüst)

**Schnittstellen:**
- Liefert: `Tropfen()`, `Globus()`, `Diagramm()`, `Thermometer()`, `Wind()`, `Druck()`.
- Liefert außerdem im Wurzelpaket: `func Kopfzeile(k KopfDaten) template.HTML`, `func Fusszeile(f FussDaten) template.HTML` samt der beiden Datentypen.

- [ ] **Schritt 1: Den fehlschlagenden Test schreiben**

`layout_test.go`:

```go
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
```

- [ ] **Schritt 2: Test laufen lassen, Fehlschlag bestätigen**

Ausführen: `go test . -run TestKopfzeile -v`
Erwartet: Übersetzungsfehler — `undefined: Kopfzeile`.

- [ ] **Schritt 3: Die Messwert-Symbole zeichnen**

`icons/messwerte.go`:

```go
package icons

// Symbole für Messgrößen.

// Tropfen — Taupunkt und Feuchte. Ersetzt das Emoji in Tempus.
func Tropfen() Icon {
	return strich(`<path d="M12 2.7l5.3 5.3a7.5 7.5 0 11-10.6 0z"/>`)
}

// Globus — großräumige Angaben wie Bioklima.
func Globus() Icon {
	return strich(`<circle cx="12" cy="12" r="10"/><path d="M2 12h20M12 2a15 15 0 010 20 15 15 0 010-20z"/>`)
}

// Diagramm — zusammengefasste Werte über einen Zeitraum.
func Diagramm() Icon {
	return strich(`<path d="M3 3v18h18"/><path d="M7 15l4-4 3 3 5-6"/>`)
}

func Thermometer() Icon {
	return strich(`<path d="M14 14.76V3.5a2.5 2.5 0 00-5 0v11.26a4.5 4.5 0 105 0z"/>`)
}

func Wind() Icon {
	return strich(`<path d="M9.6 4.6A2 2 0 1111 8H2M12.6 19.4A2 2 0 1014 16H2M17.7 7.7A2.5 2.5 0 1119.5 12H2"/>`)
}

func Druck() Icon {
	return strich(`<circle cx="12" cy="12" r="9"/><path d="M12 12l4-3M12 7v1"/>`)
}

func messwertsymbole() map[string]Icon {
	return map[string]Icon{
		"tropfen":     Tropfen(),
		"globus":      Globus(),
		"diagramm":    Diagramm(),
		"thermometer": Thermometer(),
		"wind":        Wind(),
		"druck":       Druck(),
	}
}
```

In `Alle()` ergänzen: `messwertsymbole()`.

- [ ] **Schritt 4: Das Gerüst schreiben**

`layout.go`:

```go
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
```

- [ ] **Schritt 5: Die Gerüstregeln in `base.css` ergänzen**

An `css/base.css` anhängen. Beachte: keine Farbe im Klartext, sonst schlägt
`TestBaseCSSEnthaeltKeineFarbliterale` an.

```css
/* --- Seitengerüst --- */

/* Kopf, Inhalt und Fuß teilen sich .container und stehen damit in derselben
   Spur. In Expertus begann der Kopf zuvor am Fensterrand, während der Inhalt
   zentriert war — die Seite wirkte zerrissen. */

.ds-kopf {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  justify-content: space-between;
  gap: 0.5rem 1rem;
  padding-bottom: 1rem;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid var(--border);
}

.ds-kopf h1 {
  margin: 0;
  font-size: 1.5rem;
  /* --accent-text, nicht --accent: seit der Umstellung auf Grün ist
     --accent eine Fläche (dunkel in beiden Themen, für weiße Schrift) und
     als Textfarbe auf dunklem Grund mit 1,91:1 unlesbar. */
  color: var(--accent-text);
}

.ds-kopf-titel p {
  margin: 0.125rem 0 0;
  font-size: 0.875rem;
}

.ds-fuss {
  padding-top: 1rem;
  margin-top: 2.5rem;
  border-top: 1px solid var(--border);
  font-size: 0.8125rem;
}

.ds-fuss p {
  margin: 0;
}

.ds-fuss-fassung {
  margin-top: 0.4rem;
  font-variant-numeric: tabular-nums;
}

/* --- Symbole --- */

/* Ein Symbol bemisst sich an der Schrift daneben, nicht an einer festen
   Pixelzahl: bei 200 % Textgröße wächst es mit (WCAG 1.4.4). */
.icon {
  width: 1.25em;
  height: 1.25em;
  flex: none;
  vertical-align: -0.2em;
}

/* Im Knopf steht das Symbol neben der Beschriftung. */
.btn .icon {
  margin-right: 0.4em;
}

/* Ein Knopf, der nur ein Symbol trägt, bleibt quadratisch und behält die
   Mindestgröße für Zielflächen. */
.btn-icon {
  width: 2.75rem;
  padding: 0;
}

.btn-icon .icon {
  margin: 0;
}
```

- [ ] **Schritt 6: Test laufen lassen, Erfolg bestätigen**

```bash
go build ./... && go vet ./... && go test ./... -v && gofmt -l .
```
Erwartet: alles grün, `gofmt -l` leer. Beachte, dass jetzt auch
`TestFarbpaarungenInBaseCSSErreichenAAA` über die neuen Regeln läuft.

- [ ] **Schritt 7: Einchecken**

```bash
git add icons/ layout.go layout_test.go css/base.css
git commit -m "feat: Seitengerüst und Messwert-Symbole

Ortus und Tempus führten Kopf- und Fußzeile zeichengleich, bis auf drei
Wörter. Name, Untertitel und Verweise werden jetzt übergeben, das Markup
kommt aus dem Modul — samt Maskierung, weil die Werte aus der
Konfiguration eines Dienstes stammen."
```

---

### Aufgabe 6: Referenzseite und Dokumentation

**Dateien:**
- Ändern: `demo/index.html`, `demo/demo_test.go`, `README.md`

**Schnittstellen:**
- Nutzt: alles Vorherige.

Die Referenzseite ist laut Entwurf für die Mehrzahl der Bestandteile die
einzige Prüfstelle vor dem produktiven Einsatz. Ein Symbol, das dort nicht
auftaucht, hat niemand angesehen.

- [ ] **Schritt 1: Den fehlschlagenden Test schreiben**

An `demo/demo_test.go` anhängen:

```go
func TestDemoZeigtJedesSymbol(t *testing.T) {
	// Die Referenzseite ist die einzige Stelle, an der die Symbole vor dem
	// Einsatz angesehen werden. Ein Symbol, das hier fehlt, wurde nie
	// beurteilt — und bei den Mondphasen war genau die Unterscheidbarkeit
	// der Anlass für dieses Paket.
	seite := string(demoHTML)
	for name := range icons.Alle() {
		if !strings.Contains(seite, `data-icon="`+name+`"`) {
			t.Errorf("die Referenzseite zeigt das Symbol %q nicht", name)
		}
	}
}

func TestDemoZeigtDasSeitengeruest(t *testing.T) {
	seite := string(demoHTML)
	for _, teil := range []string{`class="ds-kopf"`, `class="ds-fuss"`} {
		if !strings.Contains(seite, teil) {
			t.Errorf("die Referenzseite zeigt %q nicht", teil)
		}
	}
}
```

Der Import `icons "github.com/jobrunner/fieldworksdiary-designsystem/icons"` kommt dazu.

- [ ] **Schritt 2: Test laufen lassen, Fehlschlag bestätigen**

Ausführen: `go test ./demo/ -run TestDemoZeigtJedesSymbol -v`
Erwartet: FAIL für jedes der Symbole.

- [ ] **Schritt 3: Die Referenzseite erweitern**

Die Seite ist statisches HTML, die Symbole kommen aber aus Go. Damit beides
zusammenfindet, wird die Seite beim Ausliefern ergänzt: `demo.Handler()`
ersetzt einen Platzhalter durch die Symbolübersicht. Jedes Symbol bekommt
`data-icon="<name>"`, damit der Test es findet und beim Ansehen erkennbar
ist, welches Symbol man vor sich hat.

Ergänze in `demo/index.html` einen Abschnitt:

```html
      <section class="card">
        <h2 class="card-title">Seitengerüst</h2>
        <p class="muted">
          Kopf- und Fußzeile kommen als Bauform aus dem Modul; die
          Beschriftungen übergibt der Dienst.
        </p>
        <header class="ds-kopf">
          <div class="ds-kopf-titel">
            <h3>Ortus</h3>
            <p class="muted">Point-in-Polygon Abfrage über Datenquellen</p>
          </div>
        </header>
        <footer class="ds-fuss">
          <p class="ds-fuss-verweise">
            <a href="/designsystem.css">API Dokumentation</a> &middot;
            <a href="/designsystem.css">OpenAPI Spec</a> &middot;
            <a href="/designsystem.css">Health Status</a>
          </p>
          <p class="ds-fuss-fassung muted">ortus 1.4.2</p>
        </footer>
      </section>

      <section class="card">
        <h2 class="card-title">Symbole</h2>
        <p class="muted">
          Alle Symbole nehmen über <code>currentColor</code> die Farbe ihres
          Umfelds an. Die acht Mondphasen müssen auf einen Blick
          unterscheidbar sein — insbesondere zunehmend gegen abnehmend.
        </p>
        <div class="icon-galerie">__SYMBOLE__</div>
      </section>
```

Im Handler (`demo/demo.go`) den Platzhalter ersetzen: für jedes Symbol aus
`icons.Alle()` eine Kachel mit dem SVG, dem Namen als Beschriftung und
`data-icon="<name>"` am umschließenden Element. Sortiere die Namen, damit die
Reihenfolge zwischen zwei Aufrufen gleich bleibt.

Ergänze die Galerie-Regeln in `css/base.css`:

```css
/* Übersicht der Symbole auf der Referenzseite. */
.icon-galerie {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(7rem, 1fr));
  gap: 0.75rem;
}

.icon-galerie figure {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.35rem;
  margin: 0;
  padding: 0.75rem 0.25rem;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
}

.icon-galerie figcaption {
  font-size: 0.6875rem;
  color: var(--text-muted);
  text-align: center;
  word-break: break-word;
}

.icon-galerie .icon {
  width: 1.75em;
  height: 1.75em;
}
```

- [ ] **Schritt 4: Test laufen lassen, Erfolg bestätigen**

```bash
go build ./... && go vet ./... && go test ./... -v && gofmt -l .
```

- [ ] **Schritt 5: Die Mondphasen mit eigenen Augen prüfen**

```bash
go run ./cmd/demo
```

`http://127.0.0.1:5180` öffnen, zur Symbolübersicht gehen. **Die Prüfung, auf
die es ankommt:** Sind die acht Mondphasen bei ihrer tatsächlichen Größe
auseinanderzuhalten — besonders zunehmende gegen abnehmende Sichel und
erstes gegen letztes Viertel? Wenn nicht, sind die Pfade zu überarbeiten; der
Test kann nur feststellen, dass sie verschieden *sind*, nicht dass man sie
unterscheiden *kann*. Einmal im hellen, einmal im dunklen Thema ansehen.

- [ ] **Schritt 6: Das README ergänzen**

Um einen Abschnitt zu den Symbolen und zum Gerüst:

```markdown
## Symbole

    import "github.com/jobrunner/fieldworksdiary-designsystem/icons"

    icons.Standort()      // <svg …> für den Standort-Knopf
    icons.WetterNebel()   // Zeichnung; welcher WMO-Code das ist, weiß der Dienst
    icons.Alle()          // alle Symbole unter ihrem Namen

Jedes Symbol folgt derselben Bauregel: 24er-Raster, keine festen Maße,
`currentColor`, `aria-hidden`. Ein Symbol nimmt damit die Farbe seines
Umfelds an — es folgt Thema und Zustand, ohne eine eigene Farbe zu kennen.
Die Größe bestimmt das CSS über die Klasse `.icon`.

Das Modul liefert Zeichnungen, keine Fachlogik: welcher Zustand welches
Symbol bekommt, entscheidet der Dienst.

## Seitengerüst

    designsystem.Kopfzeile(designsystem.KopfDaten{Name: "Ortus", Untertitel: "…"})
    designsystem.Fusszeile(designsystem.FussDaten{Verweise: …, Name: "ortus", Fassung: "1.4.2"})

Name, Untertitel und Verweise übergibt der Dienst; Auszeichnung und Klassen
kommen aus dem Modul. Alle Werte werden maskiert.
```

- [ ] **Schritt 7: Einchecken**

```bash
git add demo/ css/base.css README.md
git commit -m "feat: Referenzseite zeigt Gerüst und alle Symbole

Ein Symbol, das die Seite nicht zeigt, hat niemand vor dem Einsatz
angesehen — der Test leitet die erwartete Liste deshalb aus icons.Alle()
ab statt sie zu wiederholen."
```

---

## Abschluss

Nach Aufgabe 6 steht:

- rund 40 Symbole, jedes nach derselben Bauregel, jedes geprüft
- Kopf- und Fußzeile als Bauform, mit Maskierung
- die Referenzseite zeigt beides in hellem und dunklem Thema

Erst danach beginnt die Umstellung von Expertus — es nutzt Gerüst und Symbole
dann von Anfang an, statt beides kurz darauf zu ersetzen.

## Bewusst nicht enthalten

- Die Zuordnung von WMO-Codes zu Wettersymbolen. Sie ist Wetterlogik und
  bleibt in Tempus.
- Ein SVG-Sprite als zweiter Ausspielweg. Die Go-Funktionen genügen für alle
  fünf Dienste; ein zweiter Weg für dieselben Symbole liefe auseinander.
- Symbole in mehreren Größen oder Strichstärken. Eine Form, über CSS skaliert.
