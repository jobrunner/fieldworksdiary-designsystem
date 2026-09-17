# fieldworksdiary-designsystem

Gemeinsame Gestaltungsgrundlage der fieldworksdiary-Dienste — Ortus, Tempus,
Situs, Hostus und Expertus.

Das Modul liefert zwei Stylesheets: `css/tokens.css` mit den Farb-, Typo- und
Abstandsvariablen für helles und dunkles Thema, und `css/base.css` mit Reset
und den Komponenten, die mehr als ein Dienst braucht. Verhalten bleibt Sache
der Dienste; hier liegt ausschließlich CSS.

Alle Farbwerte sind auf 7:1 für Text und 3:1 für Bedienelement-Ränder
ausgelegt. `contrast_test.go` rechnet das im Build nach.

**Status: in Planung.** Der Entwurf steht in [docs/design.md](docs/design.md);
Code gibt es noch nicht.
