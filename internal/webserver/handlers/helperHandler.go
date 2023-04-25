package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/structs"
	"errors"
	"net/http"
	"strings"
)

// InitHandler Duplicate initializer code from handlers.
func InitHandler(w http.ResponseWriter, r *http.Request) ([]structs.RenewableShareEnergyElement, error) {
	// Checks the request type.
	if !utility.CheckRequest(r, http.MethodGet) {
		http.Error(w, "Request not supported.", http.StatusNotImplemented)
		return nil, errors.New("faulty request method")
	}
	// Sets the content type of client to be json format.
	w.Header().Set("content-type", "application/json")

	// Reads from csv and returns json list.
	listOfRSE, jsonError := utility.RSEToJSON()
	if jsonError != nil {
		http.Error(w, jsonError.Error(), http.StatusInternalServerError)
		return nil, jsonError
	}
	return listOfRSE, nil
}

// SortQueryHandler Handler duplicate queries from history and current endpoint.
func SortQueryHandler(r *http.Request, list []structs.RenewableShareEnergyElement) []structs.RenewableShareEnergyElement {

	// Checks if sortByValue query is passed. If so it sorts it by percentage.
	if r.URL.Query().Has("sortbyvalue") && strings.Contains(strings.ToLower(r.URL.Query().Get("sortbyvalue")), "true") {
		// Sorts percentage descending if descending query is true.
		if strings.Contains(strings.ToLower(r.URL.Query().Get("descending")), "true") {
			list = utility.SortRSEList(list, false, constants.DESCENDING)
		} else { // Sorting standard is ascending if nothing else is passed.
			list = utility.SortRSEList(list, false, constants.ASCENDING)
		}
	}

	// Checks if sortAlphabetically query is passed.
	if r.URL.Query().Has("sortalphabetically") && strings.Contains(strings.ToLower(r.URL.Query().Get("sortalphabetically")), "true") {
		// Sorts list descending if descending query is true.
		if strings.Contains(strings.ToLower(r.URL.Query().Get("descending")), "true") {
			list = utility.SortRSEList(list, true, constants.DESCENDING)
		} else { // Sorting standard is ascending if nothing else is passed.
			list = utility.SortRSEList(list, true, constants.ASCENDING)
		}
	}
	return list
}
