package utility

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"encoding/json"
	"net/http"
)

// GetCountryFromAPI Function which gets data as byte slice based on country code search parameter.
func GetCountryFromAPI(countryIdentifier string, countryCode bool) (structs.Country, error) {
	// Declare variables used.
	var client http.Client
	var countryFromAPI []structs.Country
	var resp *http.Response
	var getError error

	// One method to retrieve based on country name and code.
	if !countryCode {
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
