package utility

import (
	"assignment-2/internal/webserver/cache"
	"assignment-2/internal/webserver/structs"
)

// GetCountry Function which finds the country code based on name of country. Uses countryCache and API if not found
// in cache.
func GetCountry(countryIdentifier string, countryCode bool) (structs.Country, error) {
	// Checks cache if countryFromCache exists.
	var countryFromCache structs.Country
	var cacheErr error

	// Checks if iso code.
	if countryCode {
		countryFromCache, cacheErr = cache.GetCountryByIsoCodeFromCache(countryIdentifier)
	} else {
		countryFromCache, cacheErr = cache.GetCountryFromCache(countryIdentifier)
	}
	// Retrieves from api and also adds to cache if an error is returned.
	if cacheErr != nil {
		// Retrieve country information from API.
		country, retrievalErr := GetCountryFromAPI(countryIdentifier, countryCode)
		if retrievalErr != nil {
			return structs.Country{}, retrievalErr
		}
		// Adds information to cache if it does not exist.
		_ = cache.AddCountryToCache(country)
		return country, nil
	} else { // If cache does not return an error, country is found.
		return countryFromCache, nil
	}
}
