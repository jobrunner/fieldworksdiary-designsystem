// Startet die Referenzseite des Design-Systems zur visuellen Abnahme.
package main

import (
	"log"
	"net/http"

	designsystem "github.com/jobrunner/fieldworksdiary-designsystem"
)

func main() {
	log.Println("Referenzseite auf http://127.0.0.1:5180")
	log.Fatal(http.ListenAndServe("127.0.0.1:5180", designsystem.DemoHandler()))
}
