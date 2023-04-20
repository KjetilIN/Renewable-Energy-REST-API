package cache

import (
	"assignment-2/internal/webserver/structs"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestAddAndGetCache Tests adding and retrieving from cache.
func TestAddAndGetCache(t *testing.T) {
	addToCacheErr := AddCountryToCache(structs.Country{
		Name:        map[string]interface{}{"common": "Norway"},
		CountryCode: "NOR",
		Borders:     []string{"SWE", "FIN", "RUS"},
		Cache:       time.Time{},
	})
	if addToCacheErr != nil {
		t.Fatal("Error adding to cache.")
	}

	country, getFromCacheErr := GetCountryFromCache("norway")
	if getFromCacheErr != nil {
		t.Fatal("Error when retrieving cached element.")
	}
	assert.Equal(t, country.CountryCode, "NOR", "Country code is not correct.")

	country1, getFromCacheErr := GetCountryByIsoCodeFromCache("NOR")
	assert.Equal(t, country1, country, "Countries is not the same.")

	addToCacheErr = AddCountryToCache(structs.Country{
		Name:        map[string]interface{}{"common": "Norway"},
		CountryCode: "NOR",
		Borders:     []string{"SWE", "FIN", "RUS"},
		Cache:       time.Time{},
	})
	assert.Error(t, addToCacheErr, "Error was not sent.")

}
