package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/structs"
	"net/http"
	"strings"
)

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
