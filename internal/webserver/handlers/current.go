package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/structs"
	"net/http"
	"strings"
)

// HandlerCurrent is a handler for the /current endpoint.
func HandlerCurrent(w http.ResponseWriter, r *http.Request) {
	// Query for printing information about endpoint.
	if r.URL.Query().Has("information") && strings.Contains(strings.ToLower(r.URL.Query().Get("information")), "true") {
		_, writeErr := w.Write([]byte("To use API, remove ?information=true, from the URL.\n"))
		if writeErr != nil {
			return
		}
		utility.Encoder(w, constants.CURRENT_QUERIES)
		return
	}

	// Runs initialise method for handler.
	originalList, initError := InitHandler(w, r)
	if initError != nil {
		return
	}

	// Retrieves the list of current year records.
	currentList := getCurrentList(originalList)

	// Collects parameter from url path. Returns empty string if none exists.
	countryIdentifier := utility.GetParams(r.URL.Path, constants.CURRENT_PATH)
	if countryIdentifier != "" {
		filteredList, filterErr := CountryFilterer(w, currentList, countryIdentifier)
		if filterErr != nil {
			return
		}

		// Checks if query neighbours is presented.
		if r.URL.Query().Has("neighbours") && len(filteredList) > 0 {
			if strings.ToLower(r.URL.Query().Get("neighbours")) == "true" {
				// Collects iso code from filtered list. Filtered list is never nil or empty.
				countryIdentifier = filteredList[0].IsoCode
				// Retrieves the neighbour countries using country code.
				neighbourList, neighbourErr := retrieveNeighbours(currentList, countryIdentifier)
				if neighbourErr != nil {
					http.Error(w, "Error:"+neighbourErr.Error(), http.StatusInternalServerError)
					return
				}
				// Sets the filtered list to currentList, which is the one to be shown.
				currentList = filteredList
				// Appends neighbours into the list to be shown.
				for _, v := range neighbourList {
					currentList = append(currentList, v)
				}
			} else { // Neighbours query is not correctly prompted.
				http.Error(w, "Neighbour query must be =true to work.", http.StatusBadRequest)
				return
			}
		} else {
			// If neighbours is not passed, the filtered list is the one to be shown.
			currentList = filteredList
		}
	}

	// Handles sorting queries.
	var sortErr error
	currentList, sortErr = SortQueryHandler(r, currentList)
	if sortErr != nil {
		http.Error(w, sortErr.Error(), http.StatusBadRequest)
		return
	}

	// If list is empty, error is passed.
	if len(currentList) == 0 {
		http.Error(w, "No search results matching your parameters.", http.StatusBadRequest)
		return
	}
	// Resets country identifier.
	countryIdentifier = ""
	// Encodes currentList to the client.
	utility.Encoder(w, currentList)
}

// getCurrentYear Retrieves the latest year. Which in turn is the largest number.
func getCurrentYear(list []structs.RenewableShareEnergyElement) int {
	currentYear := 0
	for _, v := range list {
		if currentYear < v.Year { // If year is lower than current year it will be replaced.
			currentYear = v.Year
		}
	}
	return currentYear // Returns the current year.
}

// getCurrentList Retrieves the list of element corresponding to the current year.
func getCurrentList(originalList []structs.RenewableShareEnergyElement) []structs.RenewableShareEnergyElement {
	// Current year.
	currentYear := getCurrentYear(originalList)
	var currentList []structs.RenewableShareEnergyElement

	// Iterates through the original list to collect current year elements.
	for _, v := range originalList {
		if v.Year == currentYear {
			currentList = append(currentList, v)
		}
	}
	return currentList
}

// retrieveNeighbours Checks the neighbouring countries and includes them in output list.
func retrieveNeighbours(list []structs.RenewableShareEnergyElement, countryCode string) ([]structs.RenewableShareEnergyElement, error) {
	// Collects country from API.
	country, countryGetErr := utility.GetCountry(countryCode, true)
	if countryGetErr != nil {
		return nil, countryGetErr
	}
	var neighbourList []structs.RenewableShareEnergyElement

	// Iterates through borders and list to append neighbours.
	for _, v := range country.Borders {
		for _, w := range list {
			if strings.ToLower(v) == strings.ToLower(w.IsoCode) {
				neighbourList = append(neighbourList, w)
			}
		}
	}
	return neighbourList, nil
}
