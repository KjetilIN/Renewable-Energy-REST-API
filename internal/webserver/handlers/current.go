package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/structs"
	"encoding/json"
	"net/http"
	"strings"
)

// COUNTRY_CODE_RETRIVAL Constant used in api retrieval function.
const COUNTRY_CODE_RETRIVAL = 0

// COUNTRY_NAME_RETRIVAL Constant used in api retrieval function.
const COUNTRY_NAME_RETRIVAL = 1

// HandlerCurrent is a handler for the /current endpoint.
func HandlerCurrent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	// Collects the CSV list into JSON.
	originalList, err := utility.RSEToJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	currentList := getCurrentList(originalList)

	// Collects parameters, separated by /
	params := strings.Split(r.URL.Path, "/") //Used to split the / in path to collect search parameters.
	// Checks if an optional parameter is passed.
	if len(params) == 6 {
		countryIdentifier := params[5]
		var filteredList []structs.RenewableShareEnergyElement

		// Checks if countryIdentifier is not empty, and then if it is less or more than 3 characters,
		// if so it is not a country code.
		if len(countryIdentifier) > 0 && len([]byte(countryIdentifier)) != 3 {
			// Parses country name to country code.
			countryCode, getError := parseCCToCountryName(countryIdentifier)
			if getError != nil {
				http.Error(w, "Error when parsing country name to country code: "+getError.Error(), http.StatusBadRequest)
				return
			}
			countryIdentifier = countryCode
		}
		// Gets the countries based on country code, uses api.
		filteredList = countryCodeLimiter(currentList, countryIdentifier)

		// Checks if query neighbours is presented.
		if len(filteredList) > 0 && strings.ToLower(r.URL.Query().Get("neighbours")) == "true" {
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
		} else {
			// If neighbours is not passed, the filtered list is the one to be shown.
			currentList = filteredList
		}
	}
	// If list is empty, error is passed.
	if len(currentList) == 0 {
		http.Error(w, "No search results matching your parameters.", http.StatusNotFound)
		return
	}
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
	country, countryGetErr := getCountryFromAPI(countryCode, COUNTRY_CODE_RETRIVAL)
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

// parseCCToCountryName Function which finds the country code based on name of country.
func parseCCToCountryName(countryName string) (string, error) {
	// Retrieve country code based on struct.
	country, retrievalErr := getCountryFromAPI(countryName, COUNTRY_NAME_RETRIVAL)
	if retrievalErr != nil {
		return "", retrievalErr
	}
	return country.CountryCode, nil
}

// getCountryFromAPI Function which gets data as byte slice based on country code search parameter.
func getCountryFromAPI(countryIdentifier string, retrievalSetting int) (structs.Country, error) {
	// Declare variables used.
	var client http.Client
	var countryFromAPI []structs.Country
	var resp *http.Response
	var getError error

	// One method to retrieve based on country name and code.
	if retrievalSetting == COUNTRY_NAME_RETRIVAL {
		resp, getError = client.Get(constants.COUNTRYNAME_API_ADDRESS + countryIdentifier)
	} else {
		resp, getError = client.Get(constants.COUNTRYCODE_API_ADDRESS + countryIdentifier)
	}
	// Performs a get request to country api using country code search parameter.
	if getError != nil {
		return structs.Country{}, getError
	}
	defer resp.Body.Close() //Waits for the body to return, then closes the request.
	// Decodes body into countryFromAPI struct.
	err := json.NewDecoder(resp.Body).Decode(&countryFromAPI)
	if err != nil {
		return structs.Country{}, err
	}
	// Only one country returned, therefore first index is the correct country.
	return countryFromAPI[0], nil
}
