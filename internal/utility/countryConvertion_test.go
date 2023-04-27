package utility

import (
	"assignment-2/internal/webserver/cache"
	"assignment-2/internal/webserver/structs"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestGetCountry Tests the get country.
func TestGetCountry(t *testing.T) {
	// Creates a country for testing.
	testCountry := structs.Country{
		Name:        map[string]interface{}{"common": "Norway"},
		CountryCode: "NOR",
		Borders:     []string{"SWE", "FIN", "RUS"},
		Cache:       time.Time{},
	}
	// Adds the test country to cache.
	addToCacheErr := cache.AddCountryToCache(testCountry)
	if addToCacheErr != nil {
		t.Fatal("Add to cache error.")
	}
	// Gets the country information, will access the cache.
	country, getErr := GetCountry("norway", false)
	if getErr != nil {
		t.Fatal("Error retrieving country.")
	}
	// Checks if the information is the same.
	assert.Equal(t, country.Borders, testCountry.Borders, "Countries is not the same.")
	country, getErr = GetCountry("nor", true)
	if getErr != nil {
		t.Fatal("Error retrieving country.")
	}
	assert.Equal(t, country.Borders, testCountry.Borders, "Countries is not the same.")
}
