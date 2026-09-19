// Verhalten der Koordinateneingabe, verallgemeinert aus der ausgereiftesten
// der drei Fassungen bei Ortus (internal/adapters/http/frontend.go, etwa
// Zeile 1020–1140). Ortus bietet eine Auswahl aus sieben Systemen mit
// tauschender Achsenreihenfolge und einem MGRS-Sonderfall (ein Textfeld
// statt zweier Zahlenfelder); Tempus kennt nur WGS 84 und keine Auswahl.
// Dieses Modul muss beides können. Ohne Bündler und ohne Abhängigkeiten
// lauffähig, genau wie js/combobox.js — deshalb ein ES-Modul ohne import
// von außerhalb.

// DEFAULT_ID_PREFIX ist die Vorgabe für idPrefix — aus demselben Grund wie
// DEFAULT_OPTION_ID_PREFIX in combobox.js: zwei Koordinateneingaben auf
// einer Seite (z. B. "von" und "nach") dürften sich sonst die id-Namen
// ihrer Felder teilen.
const DEFAULT_ID_PREFIX = 'koord'

// --- Zerlegung eines eingefügten Koordinatenpaars, ohne DOM -------------
//
// Aus demselben Grund wie erzeugeZustand() in combobox.js getrennt vom
// DOM-Code unten: die Fallunterscheidung (Semikolon, Komma als Trenner
// oder als Dezimaltrenner, Leerzeichen, ein einzelner Wert) ist genug
// Logik, um still falsch zu sein, und exportiert, damit sie ohne Browser
// und ohne DOM geprüft werden kann.
//
// Der wichtigste der fünf Fälle ist der letzte: ein EINZELNER Wert mit
// deutschem Dezimalkomma ("35,016132") darf nicht als Paar gelesen werden.
// Eine naive Zerlegung an jedem Komma macht daraus die zwei Zahlen 35 und
// 016132 — beide für sich gültig, das Ergebnis trotzdem falsch, und zwar
// unauffällig falsch.
export function parseKoordinatenPaar(text) {
  const t = (text || '').trim()
  if (!t) return null
  const hatSemikolon = t.indexOf(';') >= 0
  const hatKomma = t.indexOf(',') >= 0
  const hatPunkt = t.indexOf('.') >= 0
  const hatLeerzeichen = /\s/.test(t)
  let teile, kommaIstDezimal
  if (hatSemikolon) {
    teile = t.split(';'); kommaIstDezimal = true // "35,01;32,67" -> Komma ist Dezimaltrenner
  } else if (hatKomma && hatPunkt) {
    teile = t.split(','); kommaIstDezimal = false // "35.01, 32.67" -> Punkt Dezimaltrenner, Komma trennt
  } else if (hatKomma && hatLeerzeichen) {
    teile = t.split(/\s+/); kommaIstDezimal = true // "35,01 32,67" -> Leerzeichen trennt, Komma ist Dezimaltrenner
  } else if (!hatKomma && hatLeerzeichen) {
    teile = t.split(/\s+/); kommaIstDezimal = false // "35.01 32.67"
  } else {
    return null // einzelner Wert (auch das deutsche "35,016132") -> normales Einfügen, der Browser übernimmt
  }
  if (!teile || teile.length < 2) return null
  const a = normZahl(teile[0], kommaIstDezimal)
  const b = normZahl(teile[1], kommaIstDezimal) // weitere Teile (>2) werden verworfen
  if (a === null || b === null) return null
  return [a, b]
}

function normZahl(s, kommaIstDezimal) {
  s = (s || '').trim()
  if (kommaIstDezimal) s = s.replace(',', '.')
  if (!/^[+-]?(\d+\.?\d*|\.\d+)$/.test(s)) return null
  return s
}

// --- DOM-Anbindung -------------------------------------------------------

// mountKoordinaten verdrahtet die Koordinateneingabe in `container` und
// liefert { destroy }. Anders als bei der Combobox erwartet dieses Modul
// KEINE frisch erzeugten Elemente, sondern bindet an ein bereits im
// Markup vorhandenes Gerüst — genau wie Ortus' eigenes Skript an
// coordGrid/groupX/groupY/labelX/labelY/coordX/coordY/coordMgrs/sridSelect
// bindet, die alle schon im HTML stehen, bevor das Skript läuft. Das
// Gerüst innerhalb von `container` (Beschriftung/Feld-Paare finden sich
// über idPrefix):
//   - optional  <select id="{idPrefix}-system">          (mehrere Systeme)
//   - <div class="koord-gitter"> mit zwei .form-group, deren Felder
//     <input id="{idPrefix}-x"> und <input id="{idPrefix}-y"> heißen
//   - optional  <input id="{idPrefix}-einzel">            (MGRS o. Ä.)
//
// systeme: [{ id, name, xLabel, yLabel, xPlaceholder, yPlaceholder,
//             yZuerst?, einzelfeld? }]
// onChange({ system, x, y, text }) läuft bei jeder Änderung; bei einem
// Einzelfeld-System steht der Wert in text, x und y bleiben leer — und
// umgekehrt.
export function mountKoordinaten({
  container,
  systeme,
  onChange,
  idPrefix = DEFAULT_ID_PREFIX,
}) {
  const idSystem = `${idPrefix}-system`
  const idX = `${idPrefix}-x`
  const idY = `${idPrefix}-y`
  const idEinzel = `${idPrefix}-einzel`

  const gitter = container.querySelector('.koord-gitter')
  const feldX = container.querySelector(`#${idX}`)
  const feldY = container.querySelector(`#${idY}`)
  const labelX = container.querySelector(`label[for="${idX}"]`)
  const labelY = container.querySelector(`label[for="${idY}"]`)
  const gruppeX = feldX.closest('.form-group')
  const gruppeY = feldY.closest('.form-group')

  const auswahl = container.querySelector(`#${idSystem}`)
  const gruppeAuswahl = auswahl ? auswahl.closest('.form-group') : null

  const feldEinzel = container.querySelector(`#${idEinzel}`)
  const labelEinzel = feldEinzel ? container.querySelector(`label[for="${idEinzel}"]`) : null
  const gruppeEinzel = feldEinzel ? feldEinzel.closest('.form-group') : null

  // Hat systeme nur EINEN Eintrag (der Tempus-Fall: keine
  // Koordinatentransformation, also keine Wahl zwischen Systemen), wird
  // keine Auswahl dargestellt — ein Auswahlfeld mit einer einzigen
  // Möglichkeit ist eine Bedienlast ohne Nutzen. Das gilt unabhängig
  // davon, ob der Dienst überhaupt ein <select>-Gerüst mitgebracht hat.
  if (gruppeAuswahl) {
    const mitAuswahl = systeme.length > 1
    gruppeAuswahl.style.display = mitAuswahl ? '' : 'none'
    if (mitAuswahl) {
      auswahl.replaceChildren(
        ...systeme.map((sys) => {
          const option = document.createElement('option')
          option.value = sys.id
          option.textContent = sys.name
          return option
        })
      )
    }
  }

  let aktuelles = systeme[0]

  function findeSystem(id) {
    return systeme.find((s) => s.id === id) || systeme[0]
  }

  // anwenden setzt Beschriftungen, Platzhalter und sichtbare Reihenfolge
  // für ein System und leert danach die Werte.
  //
  // yZuerst tauscht NUR DIE POSITION der beiden Felder im Baum — welches
  // Feld welchen Wert trägt (die Bindung feldX↔idX↔labelX bzw.
  // feldY↔idY↔labelY), bleibt unverändert. Bei WGS 84 steht die Breite
  // oben, weil man Koordinaten so SPRICHT ("52,52 Nord, 13,40 Ost"); bei
  // projizierten Systemen steht der Rechtswert oben, weil Vermesser sie so
  // SCHREIBEN. Wer Position und Bedeutung verwechselt, baut eine Eingabe,
  // die Koordinaten stillschweigend vertauscht — und das fällt erst auf,
  // wenn ein Fundort in der Nordsee liegt.
  function anwenden(sys) {
    aktuelles = sys
    const einzelfeld = !!sys.einzelfeld && !!gruppeEinzel

    if (gitter) gitter.style.display = einzelfeld ? 'none' : ''
    if (gruppeEinzel) gruppeEinzel.style.display = einzelfeld ? '' : 'none'

    if (einzelfeld) {
      if (labelEinzel) labelEinzel.textContent = sys.xLabel
      feldEinzel.placeholder = sys.xPlaceholder || ''
    } else {
      labelX.textContent = sys.xLabel
      labelY.textContent = sys.yLabel
      feldX.placeholder = sys.xPlaceholder || ''
      feldY.placeholder = sys.yPlaceholder || ''
      if (gitter) {
        if (sys.yZuerst) gitter.insertBefore(gruppeY, gruppeX)
        else gitter.insertBefore(gruppeX, gruppeY)
      }
    }

    // Beim Wechsel des Systems werden die Werte geleert: ein Rechtswert,
    // der als Breitengrad stehen bleibt, ist schlimmer als ein leeres Feld.
    feldX.value = ''
    feldY.value = ''
    if (feldEinzel) feldEinzel.value = ''

    melden()
  }

  function melden() {
    if (aktuelles.einzelfeld && gruppeEinzel) {
      onChange?.({ system: aktuelles.id, x: '', y: '', text: feldEinzel.value })
    } else {
      onChange?.({ system: aktuelles.id, x: feldX.value, y: feldY.value, text: '' })
    }
  }

  // Smart-Einfügen: ein volles Paar wie "35.016132, 32.670024" (auch
  // ";"-getrennt oder mit deutschem Dezimalkomma
  // "35,016132;32,670024") in EINES der beiden Felder verteilt es auf
  // beide — der erste Teil in das SICHTBAR erste Feld, der zweite in das
  // zweite. Ein einzelner Wert wird normal eingefügt, der Browser
  // übernimmt das selbst.
  function aufEinfuegen(e) {
    const clip = e.clipboardData || window.clipboardData
    if (!clip) return
    // 'text/plain' ist der Standardtyp; 'text' ein Fallback für ältere
    // IE/Edge-Fassungen.
    const eingefuegt = clip.getData('text/plain') || clip.getData('text')
    const paar = parseKoordinatenPaar(eingefuegt)
    if (!paar) return // einzelner Wert -> normales Einfügen zulassen
    e.preventDefault()
    const yZuerst = !!aktuelles.yZuerst
    ;(yZuerst ? feldY : feldX).value = paar[0]
    ;(yZuerst ? feldX : feldY).value = paar[1]
    melden()
  }

  feldX.addEventListener('paste', aufEinfuegen)
  feldY.addEventListener('paste', aufEinfuegen)
  feldX.addEventListener('input', melden)
  feldY.addEventListener('input', melden)
  if (feldEinzel) feldEinzel.addEventListener('input', melden)

  function aufSystemwechsel() {
    anwenden(findeSystem(auswahl.value))
  }
  if (auswahl) auswahl.addEventListener('change', aufSystemwechsel)

  anwenden(systeme[0])

  // destroy() meldet Ereignishörer ab, damit ein zweites mountKoordinaten()
  // auf demselben Gerüst nicht doppelt reagiert — dieselbe Sorge wie bei
  // destroy() in combobox.js.
  return {
    destroy() {
      feldX.removeEventListener('paste', aufEinfuegen)
      feldY.removeEventListener('paste', aufEinfuegen)
      feldX.removeEventListener('input', melden)
      feldY.removeEventListener('input', melden)
      if (feldEinzel) feldEinzel.removeEventListener('input', melden)
      if (auswahl) auswahl.removeEventListener('change', aufSystemwechsel)
    },
  }
}
