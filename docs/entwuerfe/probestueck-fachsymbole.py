# Probestueck: dieselben Motive in zwei Registern.
# Register A (Bedienung): 24er-Raster, stroke-width 2, reduziert.
# Register B (fachlich):  48er-Raster, stroke-width 1.5, mehr Detail.

A = {  # 24er-Raster, stroke 2
 'art': '<path d="M12 21v-7"/><circle cx="12" cy="9" r="2.5"/><path d="M12 6.5V4M9.5 9H7m7 0h2.5M10 7l-1.5-1.5M14 7l1.5-1.5"/>',
 'vegetationstyp': '<path d="M5 21v-6M12 21v-9M19 21v-5"/><circle cx="5" cy="13" r="1.6"/><path d="M12 12l-2.5-2.5M12 12l2.5-2.5"/><path d="M19 14.5c-1.5-1.5-1.5-3 0-4.5 1.5 1.5 1.5 3 0 4.5z"/>',
 'habitat': '<path d="M2 20h20"/><path d="M2 16l5-6 4 4 3-4 4 5 4-2"/><path d="M8 20v-3M16 20v-4"/>',
 'wald': '<path d="M12 21v-4"/><path d="M12 3L6 11h12z"/><path d="M12 9l-4.5 6h9z"/>',
 'grasland': '<path d="M3 21c1.5-4 2-7 1.5-10M8 21c2-5 2.5-9 1.5-13M13 21c-1.5-4-2-7-1.5-10M18 21c2-5 2.5-8 2-11M21 21H3"/>',
}
B = {  # 48er-Raster, stroke 1.5
 'art': '<path d="M24 43V26"/><circle cx="24" cy="18" r="4.5"/><path d="M24 13.5V7M19.5 18H11m17.5 0H37M20 13l-5-5M28 13l5-5M20 23l-5 5M28 23l5 5"/><path d="M24 33c-4-1-6.5-3.5-7.5-7M24 37c4-1 6.5-3.5 7.5-7"/>',
 'vegetationstyp': '<path d="M8 43V29M18 43V22M28 43V26M38 43V31"/><circle cx="8" cy="26" r="2.6"/><path d="M18 22l-4-4M18 22l4-4M18 28l-3.5-3.5M18 28l3.5-3.5"/><path d="M28 26c-2.5-2.5-2.5-5.5 0-8 2.5 2.5 2.5 5.5 0 8z"/><path d="M38 31c-2-2-2-4.5 0-6.5 2 2 2 4.5 0 6.5z"/><path d="M4 43h40"/>',
 'habitat': '<path d="M4 41h40"/><path d="M4 33l9-12 7 8 6-9 8 11 6-5"/><path d="M13 41v-6M27 41v-8M38 41v-5"/><path d="M4 37c4 1.5 7-1.5 11 0s7-1.5 11 0 7-1.5 11 0"/>',
 'wald': '<path d="M24 43v-7"/><path d="M24 5L14 20h20z"/><path d="M24 15l-8 12h16z"/><path d="M24 25l-10 13h20z"/>',
 'grasland': '<path d="M6 43c3-8 4-14 3-20M14 43c4-10 5-18 3-26M22 43c-3-8-4-14-3-20M30 43c4-10 5-16 4-22M38 43c3-8 3.5-13 3-18M44 43H4"/><circle cx="19" cy="20" r="1.4"/><circle cx="33" cy="24" r="1.4"/>',
}

def svg(inhalt, raster, strich, groesse):
    return (f'<svg viewBox="0 0 {raster} {raster}" fill="none" stroke="currentColor" '
            f'stroke-width="{strich}" stroke-linecap="round" stroke-linejoin="round" '
            f'aria-hidden="true" style="width:{groesse};height:{groesse}">{inhalt}</svg>')

import subprocess, urllib.request
css = urllib.request.urlopen('http://127.0.0.1:5180/designsystem.css').read().decode()

zeilen = []
for name in A:
    zeilen.append(f'''<tr><th scope="row">{name}</th>
<td>{svg(A[name],24,2,"24px")}</td><td>{svg(A[name],24,2,"3em")}</td>
<td>{svg(B[name],48,1.5,"24px")}</td><td>{svg(B[name],48,1.5,"3em")}</td></tr>''')

html = f'''<!doctype html><html lang="de"><head><meta charset="utf-8"><title>Probestück</title><style>{css}
table{{width:auto}} th,td{{padding:.7rem 1.1rem;vertical-align:middle}}
thead th{{font-size:.75rem;text-transform:uppercase;letter-spacing:.04em}}
tbody th{{font-weight:500;font-size:.85rem;text-align:left}}
td{{text-align:center}}
.trenn{{border-left:2px solid var(--border)}}
</style></head><body><div class="container">
<h1>Probestück — zwei Register</h1>
<p class="muted">Dieselben fünf Motive. Links im Bedien-Register (24er-Raster, Strichstärke 2),
rechts im fachlichen Register (48er-Raster, Strichstärke 1,5). Je einmal bei 24 Pixeln
Einsatzgröße und einmal groß.</p>
<table><thead><tr><th></th>
<th>A · 24 px</th><th>A · groß</th>
<th class="trenn">B · 24 px</th><th>B · groß</th></tr></thead>
<tbody>{"".join(zeilen)}</tbody></table>
</div></body></html>'''
open('/Users/jbrunner/work/projects/expertus/.playwright-mcp/probe.html','w').write(html)
print('erzeugt')
