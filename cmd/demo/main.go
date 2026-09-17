// Startet die Referenzseite des Design-Systems zur visuellen Abnahme.
package main

import (
	"log"
	"net/http"

	"github.com/jobrunner/fieldworksdiary-designsystem/demo"
)

func main() {
	log.Println("Referenzseite auf http://127.0.0.1:5180")
	log.Fatal(http.ListenAndServe("127.0.0.1:5180", demo.Handler()))
}
