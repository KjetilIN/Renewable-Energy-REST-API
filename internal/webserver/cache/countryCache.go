package cache

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Creates a new cache.
var cachedCountry = make(map[string]structs.Country)

// AddCountryToCache Adds country to the cache.
func AddCountryToCache(country structs.Country) error {
	// Stores the country name.
	cachedCountryName := strings.ToUpper(country.Name["common"].(string))
	// Checks if it exists in the cache.
	if _, exists := cachedCountry[cachedCountryName]; exists {
		// If it exists in cache an error is returned.
		return errors.New(fmt.Sprintf("%s is already cached.", cachedCountryName))
	}
	// Inserts the time of entry to the cache.
	country.Cache = time.Now()
	cachedCountry[cachedCountryName] = country
	return nil
}

// GetCountryByIsoCodeFromCache Retrieves a cached country by its isoCode.
func GetCountryByIsoCodeFromCache(isoCode string) (structs.Country, error) {
	isoCode = strings.ToUpper(isoCode)
	for _, country := range cachedCountry {
		if strings.ToUpper(country.CountryCode) == isoCode || strings.ToUpper(country.CountryCode) == isoCode {
			return GetCachedCountryByName(strings.ToUpper(country.Name["common"].(string)))
		}
	}
	return structs.Country{}, errors.New(fmt.Sprintf("%s is not cached.", isoCode))
}

// GetCountryFromCache Get a country by name.
func GetCountryFromCache(cachedCountryName string) (structs.Country, error) {
	cachedCountryName = strings.ToUpper(cachedCountryName)
	return GetCachedCountryByName(cachedCountryName)
}

// GetCachedCountryByName  Retrieves a cached country from the cache by its common name.
func GetCachedCountryByName(cachedCountryName string) (structs.Country, error) {
	// Checks if country exists in cache.
	if country, exists := cachedCountry[cachedCountryName]; exists {
		// Checks if cache limit is reached.
		if time.Since(country.Cache).Hours() > constants.LIMIT_CACHE_TIME {
			// Delete the cached country if it has surpassed the limit
			delete(cachedCountry, country.Name["common"].(string))
			return structs.Country{}, errors.New(fmt.Sprintf("%s was cached, but it has gone over %v hours since it was cached.", cachedCountryName, constants.LIMIT_CACHE_TIME))
		}
		// Returns country if found.
		return country, nil
	} else {
		// Returns error if entry is not found.
		return structs.Country{}, errors.New(fmt.Sprintf("%s is not cached.", cachedCountryName))
	}
}
