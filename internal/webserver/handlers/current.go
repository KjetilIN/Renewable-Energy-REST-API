package handlers

import (
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/structs"
	"net/http"
)

// HandlerCurrent is a handler for the /current endpoint.
func HandlerCurrent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	// Current year.
	currentYear := 2021

	// Collects the CSV list into JSON.
	originalList, err := utility.RSEToJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var currentList []structs.HistoricalRSE

	// Iterates through the original list to collect current year elements.
	for _, v := range originalList {
		if v.Year == currentYear {
			currentList = append(currentList, v)
		}
	}
	// Encodes currentList to the client.
	utility.Encoder(w, currentList)
}
