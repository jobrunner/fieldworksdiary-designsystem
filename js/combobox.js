// Verhalten der Combobox mit Vorschlagsliste (ARIA-Combobox-Pattern),
// verallgemeinert aus Expertus' src/views/combobox.js und
// src/combobox-state.js. Jeder Dienst bindet dieses Modul ein und liefert
// nur, was fachlich ist: die Herkunft der Vorschläge (suggest) und was ein
// Treffer bedeutet (onPick). Ohne Bündler und ohne Abhängigkeiten lauffähig
// — deshalb ein ES-Modul ohne import von außerhalb.

// DEFAULT_OPTION_ID_PREFIX ist die Vorgabe für idPrefix (siehe
// mountCombobox). Zwei Comboboxen auf derselben Seite mit der Vorgabe
// erzeugten beide id="combobox-option-0" — aria-activedescendant zeigte
// dann möglicherweise in die falsche Liste. Ein Dienst mit mehr als einer
// Combobox auf einer Seite übergibt deshalb je Instanz ein eigenes
// idPrefix; das Modul kann die Kollision nicht selbst lösen, weil es keine
// der beiden Instanzen kennt, solange nicht beide verdrahtet sind.
const DEFAULT_OPTION_ID_PREFIX = 'combobox-option-'

// --- Zustandslogik, ohne DOM -------------------------------------------
//
// Getrennt von der DOM-Anbindung unten, aus demselben Grund wie in Expertus'
// combobox-state.js: das ARIA-Combobox-Pattern hat genug Zustand (offen,
// aktiver Index, activedescendant), um still falsch zu sein, und ein Test,
// der nur die Attribute im DOM prüft, sieht davon nichts. Eine eigene Datei
// wäre hier ein zweiter Baustein mehr, als die Aufgabe vorsieht (nur
// js/combobox.js ist vorgesehen) — die Trennung bleibt deshalb innerhalb
// dieser Datei bestehen, nicht als eigene Datei.
function erzeugeZustand(idPrefix) {
  let query = ''
  let options = []
  let open = false
  let activeIndex = -1

  function momentaufnahme() {
    return {
      query,
      options,
      open,
      activeIndex,
      activeId: activeIndex >= 0 ? `${idPrefix}${activeIndex}` : null,
    }
  }

  return {
    getState: momentaufnahme,

    setOptions(liste) {
      options = liste ?? []
      open = options.length > 0
      activeIndex = -1
    },

    setQuery(q) {
      query = q
      activeIndex = -1
    },

    close() {
      open = false
      activeIndex = -1
    },

    move(delta) {
      if (!open || options.length === 0) return
      const n = options.length
      activeIndex = activeIndex < 0 ? (delta > 0 ? 0 : n - 1) : (activeIndex + delta + n) % n
    },

    select(index) {
      if (!Number.isInteger(index) || index < 0 || index >= options.length) return
      activeIndex = index
    },

    pick() {
      if (activeIndex < 0) return null
      const gewaehlt = options[activeIndex]
      open = false
      activeIndex = -1
      return gewaehlt
    },
  }
}

// --- DOM-Anbindung -------------------------------------------------------

// mountCombobox verdrahtet ein Eingabefeld (role="combobox") mit einer Liste
// (role="listbox") und liefert { destroy }. suggest ist der einzige
// Anknüpfungspunkt an die Fachlichkeit — was ein Treffer bedeutet und woher
// er kommt, entscheidet der Dienst; das Modul kennt nur id, text und einen
// optionalen hinweis, der gedämpft hinter dem Text steht.
//
// idPrefix: siehe DEFAULT_OPTION_ID_PREFIX — nur bei mehr als einer
// Combobox auf derselben Seite nötig, dann je Instanz ein eigener Wert.
//
// clearOnPick (Vorgabe: true): nach einer Auswahl wird das Eingabefeld
// geleert. Das ist Expertus' Semantik — dort wird wiederholt eine Art
// HINZUGEFÜGT, das Feld muss für die nächste Eingabe leer sein, und dieser
// Fall ist über die fünf Dienste hinweg der häufigere (mehrfaches
// Hinzufügen in einer Sitzung ist die Regel, eine einmalige Ortsauswahl
// die Ausnahme). Für eine ORTSAUSWAHL in Ortus oder Tempus ist Leeren
// dagegen falsch: der gewählte Ort soll im Feld stehen bleiben, sichtbar
// als das, was gerade gilt. Deshalb konfigurierbar: clearOnPick: false
// lässt den Text des gewählten Eintrags im Feld stehen (oder den
// eingegebenen Freitext, wenn kein Eintrag getroffen wurde).
export function mountCombobox({
  input,
  listbox,
  suggest,
  onPick,
  onError,
  debounceMs = 200,
  idPrefix = DEFAULT_OPTION_ID_PREFIX,
  clearOnPick = true,
}) {
  const zustand = erzeugeZustand(idPrefix)
  let timer = null
  let laufend = null

  function paint() {
    const s = zustand.getState()
    input.setAttribute('aria-expanded', String(s.open))
    if (s.activeId) input.setAttribute('aria-activedescendant', s.activeId)
    else input.removeAttribute('aria-activedescendant')

    listbox.replaceChildren()
    listbox.hidden = !s.open
    if (!s.open) return

    s.options.forEach((eintrag, i) => {
      const li = document.createElement('li')
      li.id = `${idPrefix}${i}`
      li.setAttribute('role', 'option')
      li.className = 'option'
      li.setAttribute('aria-selected', String(i === s.activeIndex))
      li.textContent = eintrag.text
      if (eintrag.hinweis) {
        const hinweis = document.createElement('span')
        hinweis.className = 'muted'
        hinweis.textContent = ` — ${eintrag.hinweis}`
        li.append(hinweis)
      }
      // mousedown statt click: click käme erst nach blur, und blur
      // schließt die Liste bereits (siehe unten) — die Auswahl ginge
      // verloren, bevor click überhaupt ausgelöst wird.
      li.addEventListener('mousedown', (e) => {
        e.preventDefault()
        zustand.select(i)
        uebernehmen()
      })
      listbox.append(li)
    })
  }

  function uebernehmen() {
    const gewaehlt = zustand.pick()
    // Freitext bleibt zulässig: ohne Hervorhebung gilt beim Übernehmen der
    // eingegebene Text, mit id: null. Manche Fachbegriffe stehen in keinem
    // Verzeichnis, aus dem suggest schöpfen könnte.
    const ergebnis = gewaehlt ?? { id: null, text: input.value.trim() }
    // clearOnPick steuert, was im Feld stehen bleibt — siehe die
    // Begründung an mountCombobox.
    input.value = clearOnPick ? '' : ergebnis.text
    zustand.setQuery(input.value)
    zustand.setOptions([])
    paint()
    onPick(ergebnis)
  }

  // Die Entprellung ist nicht verhandelbar: ohne sie löst jeder einzelne
  // Tastenanschlag sofort eine eigene Netzanfrage aus — bei einem
  // dreizehnstelligen Namen dreizehn statt zwei Anfragen. Die Dienste laufen
  // im Gelände über Mobilfunk, wo Datenvolumen, Akku und Wartezeit zählen.
  // Der AbortController weiter unten ersetzt das nicht: er verhindert nur,
  // dass eine überholte ANTWORT eine neuere überschreibt, nicht, dass die
  // ANFRAGE gestellt wird.
  function aufEingabe() {
    zustand.setQuery(input.value)
    clearTimeout(timer)
    timer = setTimeout(async () => {
      laufend?.abort()
      laufend = new AbortController()
      try {
        const treffer = await suggest(input.value, { signal: laufend.signal })
        zustand.setOptions(treffer)
        paint()
      } catch (err) {
        if (err?.name === 'AbortError') return
        zustand.setOptions([])
        paint()
        onError?.(err)
      }
    }, debounceMs)
  }

  function aufTastendruck(e) {
    if (e.key === 'ArrowDown') { e.preventDefault(); zustand.move(1); paint() }
    else if (e.key === 'ArrowUp') { e.preventDefault(); zustand.move(-1); paint() }
    else if (e.key === 'Enter') { e.preventDefault(); if (input.value.trim()) uebernehmen() }
    else if (e.key === 'Escape') {
      // Escape bricht auch eine laufende Anfrage ab — sonst füllt eine
      // spät eintreffende Antwort die gerade geschlossene Liste erneut.
      clearTimeout(timer)
      laufend?.abort()
      zustand.close()
      paint()
    }
  }

  function aufBlur() {
    zustand.close()
    paint()
  }

  input.addEventListener('input', aufEingabe)
  input.addEventListener('keydown', aufTastendruck)
  input.addEventListener('blur', aufBlur)

  paint()

  // destroy() räumt Zeitgeber, laufende Anfrage UND die drei Ereignishörer
  // auf dem Eingabefeld ab. Ohne die Hörer bleiben sie auch nach dem
  // Entfernen des Feldes aus dem Baum aktiv und halten diese Instanz am
  // Leben — ein zweiter mountCombobox() auf dasselbe Feld ließe dann beide
  // Instanzen gleichzeitig reagieren.
  return {
    destroy() {
      clearTimeout(timer)
      laufend?.abort()
      input.removeEventListener('input', aufEingabe)
      input.removeEventListener('keydown', aufTastendruck)
      input.removeEventListener('blur', aufBlur)
    },
  }
}
