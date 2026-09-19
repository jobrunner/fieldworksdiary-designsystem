# Fachsymbole — Übergabe an die nächste Sitzung

Stand: 2026-09-19

Dieses Dokument ist die Grundlage für die Fortsetzung in einer neuen Sitzung.
Es hält fest, was entschieden ist, was offen blieb, und was beim Bau der
ersten 40 Symbole gelernt wurde — damit die Fehler nicht ein zweites Mal
gemacht werden.

## Worum es geht

Das Design-System hat 40 Symbole für Bedienung, Wetter, Himmelserscheinungen
und Messwerte. Es fehlen die **fachlichen** Symbole: Pflanzen, Habitate,
Syntaxa — das, was die Dienste inhaltlich behandeln.

Als Anregung dient `floraveg.eu` (Datenbank der europäischen Vegetation,
Masaryk-Universität Brno). Ausdrücklich als **Anregung**: Die dortigen
Zeichnungen sind urheberrechtlich geschützt (© Vegetation Science Group) und
werden nicht übernommen. Übernommen wird eine Idee.

## Die Idee, die von FloraVeg.EU stammt

Dort entspricht die Abstraktionsebene des Bildes der Abstraktionsebene der
Daten:

| Bereich | Bild | Aussage |
|---|---|---|
| Species | **eine** Blüte, detailliert | ein Einzelding |
| Vegetation | **mehrere** Pflanzen verschiedener Wuchsform nebeneinander | eine Gesellschaft |
| Habitats | eine **Landschaft** mit Relief, Wasser, Zonen | ein Ort |

Das trägt die fachliche Hierarchie ins Bild: Ein Syntaxon *ist* keine
Pflanze, sondern eine wiederkehrende Kombination von Pflanzen — und genau das
zeigt das Bild. Diese Übertragung ist der Kern dessen, was übernommen werden
soll.

## Entschieden

- **Vier Motivgruppen**, alle gewünscht:
  1. **EUNIS-Habitatgruppen** — Wald, Grasland, Moor und Sumpf, Heide und
     Gestrüpp, Fels und Geröll, Binnengewässer, Küste, Acker und Siedlung
     (8). Für Expertus unmittelbar nützlich: dort wird das Ergebnis einer
     Habitatbestimmung angezeigt.
  2. **Klassifikationsebenen** — Art, Vegetationstyp (Syntaxon), Habitat,
     dazu die syntaxonomischen Ränge Klasse, Ordnung, Verband (6).
  3. **Wuchsformen** — Baum, Strauch, Zwergstrauch, Kraut, Gras, Moos,
     Flechte, Farn (8).
  4. **Standortfaktoren** — Feuchte, Nährstoffe, Licht, Bodenreaktion, Salz
     (5). Die Achsen der ökologischen Zeigerwerte.

  Zusammen rund **27 Zeichnungen**.

- **Erst ein Probestück, dann entscheiden** — das Probestück liegt vor
  (siehe unten).

## Offen — hier beginnt die neue Sitzung

Nach dem Probestück stehen drei Fragen offen, die der Mensch beantworten muss:

1. **Trifft die Bildsprache seinen Geschmack?** Besonders `vegetationstyp`
   (mehrere Wuchsformen nebeneinander als „Gesellschaft") ist der Kern der
   übernommenen Idee.

2. **Ein Register oder zwei?** Das Probestück zeigt dieselben fünf Motive im
   Bedien-Register (24er-Raster, Strichstärke 2) und in einem fachlichen
   Register (48er-Raster, Strichstärke 1,5).

   **Befund aus dem Probestück:** Das feine Register versagt bei 24 Pixeln —
   die Striche verschwinden, `art` und `grasland` zerfallen. In großer
   Darstellung gewinnt es dagegen deutlich (Blütenform mit Kelch,
   Samenstände, Wasserlinie unter dem Relief). Die Register unterscheiden
   sich also weniger im Stil als in der **Einsatzgröße**.

   **Empfehlung:** ein Register im 24er-Raster, dafür sorgfältig gezeichnet.
   Ein zweiter Satz verdoppelt die Pflegearbeit für einen Gewinn, den man nur
   groß sieht.

3. **Umfang und Reihenfolge:** 27 Zeichnungen sind ein eigenes Vorhaben. Die
   Mondphasen haben gezeigt, wie lang der Weg von „geometrisch korrekt" zu
   „auf einen Blick erkennbar" sein kann — dort brauchte es eine
   Nacharbeitsrunde nach visueller Prüfung, obwohl alle Tests grün waren.

## Das Probestück

- `docs/entwuerfe/probestueck-fachsymbole.png` — das Bild
- `docs/entwuerfe/probestueck-fachsymbole.py` — erzeugt es neu; enthält die
  Pfade beider Register für `art`, `vegetationstyp`, `habitat`, `wald`,
  `grasland`

Das Skript holt sich das CSS vom laufenden Demo-Server
(`http://127.0.0.1:5180/designsystem.css`, Start mit `go run ./cmd/demo`).

## Was beim Bau der ersten 40 Symbole gelernt wurde

Diese Punkte haben jeweils mindestens eine Nacharbeitsrunde gekostet. Sie
gehören in den Auftrag jeder Person, die Symbole zeichnet.

### Geometrie

1. **SVG-Bogenflags sind je ein Zeichen und dürfen ohne Trennzeichen
   angrenzen.** In `a4 4 0 010-8` steht `010-8` für `0`, `1`, `0`, `-8`. Ein
   Zerleger, der Zahlen gierig liest, rechnet ab da falsch. Mir ist das
   passiert; ich habe daraufhin einen intakten Pfad für fehlerhaft gehalten.

2. **Für Ellipsenbögen gilt nicht „Abstand ≤ 2·Radius".** Maßgeblich ist die
   Bedingung aus der SVG-Spezifikation, Abschnitt F.6.6: mit den
   Halbdifferenzen dx, dy muss `dx²/rx² + dy²/ry² ≤ 1` gelten. Ist der Wert
   größer, vergrößert der Renderer beide Radien **still** — gezeichnet wird
   eine andere Form als angegeben. Genau das hat ein Wettersymbol
   unleserlich gemacht.

3. **Die halbe Strichbreite zählt zur Ausdehnung.** Ein Pfad, dessen Punkte
   bis 24 reichen, ragt mit `stroke-width="2"` bis 25 und wird beschnitten.
   Das hat eine Prüfung übersehen, die nur die Punkte betrachtete.

Diese drei Prüfungen liegen inzwischen als `geometrie_test.go` vor und laufen
über `icons.Alle()`. Neue Symbole werden automatisch erfasst.

### Tests, die grün sind, ohne etwas zu beweisen

Die Schlussprüfung des Symbolpakets hat zwölf Mutationen durchgespielt — elf
blieben grün. Die Muster, die sich wiederholt haben:

1. **Eine Prüfung, die sich aus einer Liste ableitet, prüft nur, was auf der
   Liste steht.** Ein Symbol, das niemand einträgt, wird nirgends geprüft und
   nirgends gezeigt. Gegenmittel: die Anzahl der Symbolfunktionen im Quelltext
   mit der Anzahl der Einträge vergleichen.

2. **Textsuche findet auch Kommentare.** Die Combobox-Prüfung suchte
   `mousedown` im Dateitext — und fand es im Kommentar, der erklärt, *warum*
   `mousedown` nötig ist. Der Code durfte auf `click` wechseln, ohne dass
   etwas anschlug. Gegenmittel: Kommentare vor der Suche entfernen.

3. **Gültiges SVG sagt nichts über Lesbarkeit.** Ein Pfad kann syntaktisch
   einwandfrei, im Raster, mit `currentColor` versehen sein — und trotzdem
   nichts Erkennbares zeichnen. Bei den Mondphasen war die erste Fassung
   geometrisch fehlerfrei (alle Bögen exakt an der Grenze) und trotzdem
   unbrauchbar, weil eine astronomisch korrekte Sichel bei 24 Pixeln
   verschwindet.

**Die Konsequenz:** Jede Symbolgruppe wird nach dem Zeichnen **gerendert und
angesehen**, nicht nur getestet. Ein Test, der nie rot war, beweist nichts —
also jede Zusage einmal absichtlich brechen und den Test rot sehen.

### Was sich bewährt hat

- Eine messbare Größe statt einer Geschmacksfrage: Bei den Mondphasen wurde
  der **Flächenanteil** je Phase gerastert (0 / 25 / 50 / 83 / 100 %). Das
  machte „unterscheidbar" überprüfbar. Für Wuchsformen wäre etwa die
  Silhouettenhöhe eine solche Größe.
- Die Gruppenzugehörigkeit gehört ins Paket (`icons.Gruppen()`), nicht in die
  Demo — sonst erfindet die künftige Anwendung die Reihenfolge neu.
- Wiederkehrende Formen zeichengleich halten und das prüfen. Bei den Wolken
  gab es zwischenzeitlich drei verschiedene Radien für dieselbe Form.

## Verwandte Stellen

- `docs/design.md` — der Entwurf, Abschnitt „Symbole" mit der Bauregel
- `icons/icons.go` — `huelle()`, `strich()`, `flaeche()`,
  `darfFlaechenNutzen()`
- `icons/gruppen.go` — `Gruppen()`, die Reihenfolge der Symbole
- `geometrie_test.go` — die drei Geometrieprüfungen
- `demo/` — die Referenzseite; jedes neue Symbol muss dort erscheinen, ein
  Test erzwingt es
