package handlers

import (
	"assignment-2/db"
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
		return nil, errors.New("method not supported") // Error is returned, but without name. This is as http error is present.
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

// CountryFilterer Filters list based on country code or name.
func CountryFilterer(w http.ResponseWriter, list []structs.RenewableShareEnergyElement, countryIdentifier string) ([]structs.RenewableShareEnergyElement, error) {
	// Checks if country identifier exists.
	if countryIdentifier != "" {
		var filteredList []structs.RenewableShareEnergyElement
		// Tries to filter list by country code.
		filteredList = countryCodeLimiter(list, countryIdentifier)

		// Checks if filtered list is empty, if so the identifier might not be a country code. Checks for country names.
		if len(filteredList) == 0 {
			// Parses country name to country code.
			country, getCountryError := utility.GetCountry(countryIdentifier, false)
			if getCountryError != nil {
				http.Error(w, "Did not find country based on search parameters.", http.StatusBadRequest)
				return nil, errors.New("no matching country")
			}
			if country.CountryCode != "" {
				// Assigns the country identifier to be the country code from api.
				countryIdentifier = country.CountryCode
				// Using country code from api it filters list.
				filteredList = countryCodeLimiter(list, countryIdentifier)
			} else { // If country code does not exist, it is handled here.
				http.Error(w, "No country code corresponding to country.", http.StatusNotFound)
				return nil, errors.New("no country code")
			}
		}
		// Increment the invocations for the given country code
		dbErr := db.IncrementInvocations(strings.ToUpper(countryIdentifier), constants.FIRESTORE_COLLECTION)
		if dbErr != nil {
			http.Error(w, "Error: "+dbErr.Error(), http.StatusBadRequest)
			return nil, dbErr
		}
		return filteredList, nil
	} else {
		return list, nil
	}
}

// countryCodeLimiter Method to limit a list based on country code.
func countryCodeLimiter(listToIterate []structs.RenewableShareEnergyElement, countryCode string) []structs.RenewableShareEnergyElement {
	var limitedList []structs.RenewableShareEnergyElement
	for i, v := range listToIterate { // Iterates through input list.
		if strings.Contains(strings.ToLower(listToIterate[i].IsoCode), strings.ToLower(countryCode)) { // If country code match it is
			// appended to new list.
			limitedList = append(limitedList, v)
		}
	}
	return limitedList // Returns list containing all matching countries.
}
