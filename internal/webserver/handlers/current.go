package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/structs"
	"encoding/json"
	"net/http"
	"strings"
)

// HandlerCurrent is a handler for the /current endpoint.
func HandlerCurrent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	// Collects the CSV list into JSON.
	originalList, err := utility.RSEToJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Current year.
	currentYear := getCurrentYear(originalList)
	var currentList []structs.RenewableShareEnergyElement

	// Iterates through the original list to collect current year elements.
	for _, v := range originalList {
		if v.Year == currentYear {
			currentList = append(currentList, v)
		}
	}

	// Collects parameters, separated by /
	params := strings.Split(r.URL.Path, "/") //Used to split the / in path to collect search parameters.
	// Checks if an optional parameter is passed.
	if len(params) == 6 {
		countryCode := params[5]
		filteredList := countryCodeLimiter(currentList, countryCode)

		// Checks if query neighbours is presented.
		if len(filteredList) > 0 && strings.ToLower(r.URL.Query().Get("neighbours")) == "true" {
			// Retrieves the neighbour countries using country code.
			neighbourList, neighbourErr := retrieveNeighbours(currentList, countryCode)
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

// retrieveNeighbours Checks the neighbouring countries and includes them in output list.
func retrieveNeighbours(list []structs.RenewableShareEnergyElement, countryCode string) ([]structs.RenewableShareEnergyElement, error) {
	// Collects country from API.
	country, countryGetErr := getCountryFromAPI(countryCode)
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

// getCountryFromAPI Function which gets data as byte slice based on country code search parameter.
func getCountryFromAPI(countryCode string) (structs.Country, error) {
	// Declare variables used.
	var client http.Client
	var countryFromAPI []structs.Country

	// Performs a get request to country api using country code search parameter.
	resp, getError := client.Get(constants.COUNTRY_API_ADDRESS + countryCode)
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
