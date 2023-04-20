package utility

import (
	"assignment-2/internal/webserver/cache"
	"assignment-2/internal/webserver/structs"
	"fmt"
)

// GetCountry Function which finds the country code based on name of country. Uses countryCache and API if not found
// in cache.
func GetCountry(countryIdentifier string, countryCode bool) (structs.Country, error) {
	// Checks cache if countryFromCache exists.
	countryFromCache, cacheErr := cache.GetCachedCountryByName(countryIdentifier)
	if cacheErr != nil {
		// Retrieve country information from API.
		country, retrievalErr := GetCountryFromAPI(countryIdentifier, countryCode)
		if retrievalErr != nil {
			return structs.Country{}, retrievalErr
		}
		// Adds information to cache if it does not exist.
		if cache.ExistInCache(countryIdentifier, countryCode) {
			cacheAddErr := cache.AddCountryToCache(country)
			fmt.Println(cacheAddErr)
		}
		return country, nil
	} else { // If cache does not return an error, country is found.
		return countryFromCache, nil
	}
}
