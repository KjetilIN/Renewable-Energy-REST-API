package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitHandler(t *testing.T) {
	// Create a new request.
	req, err := http.NewRequest("GET", constants.HISTORY_PATH, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response.
	rr := httptest.NewRecorder()

	// Call the handler function.
	list, err := InitHandler(rr, req)
	if err != nil {
		return
	}

	assert.NotEmpty(t, list, "List is empty.")
}

// TestCountryFilterer Tests country filtering method used in current and history endpoints.
func TestCountryFilterer(t *testing.T) {
	rr := httptest.NewRecorder()
	err := utility.DirChanger(2)
	if err != nil {
		t.Fatal("Directory change error.")
	}
	list, jsonErr := prepareList(1)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}
	countryIdentifier := "NOR"

	filteredList, err := CountryFilterer(rr, list, countryIdentifier)
	if err != nil {
		t.Fatal("Filter error: " + err.Error())
	}
	for _, v := range filteredList {
		if v.IsoCode != countryIdentifier {
			t.Fatal("List is not sorted.")
		}
	}
}

// TestCountryNameLimiter Tests name limiter.
func TestCountryNameLimiter(t *testing.T) {
	// Changes directory to root.
	utility.DirChanger(2)
	// Retrieves the csv into a json file.
	list, err := utility.RSEToJSON()
	if err != nil {
		t.Fatal("Error when getting JSON.")
	}
	// Country name for testing.
	countryName := "Norway"
	list = countryNameLimiter(list, "Norway")
	// Iterates through the filtered list.
	for _, v := range list {
		if v.Name != countryName {
			t.Fatal("List is not sorted.")
		}
	}
}

// TestSortQueryHandler Tests sorting used in current and history endpoint.
func TestSortQueryHandler(t *testing.T) {
	// Prepare a list for testing.
	list, jsonErr := prepareList(1)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}
	// Create a mock request with the desired query parameters.
	req, getErr := http.NewRequest("GET", constants.HISTORY_PATH+"?descending=true&sortbyvalue=true", nil)
	if getErr != nil {
		t.Fatal(getErr)
	}

	// Call the handler function.
	sortedList, err := SortQueryHandler(req, list)

	// Check for any errors returned by the handler.
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	// Checks if list is sorted by percentage.
	for i := 1; i < len(sortedList); i++ {
		assert.GreaterOrEqualf(t, sortedList[i-1].Percentage, sortedList[i].Percentage, "List is not sorted.")
		return
	}

	req, getErr = http.NewRequest("GET", constants.HISTORY_PATH+"?sortbyvalue=true", nil)
	if getErr != nil {
		t.Fatal(err)
	}
	sortedList, sortErr := SortQueryHandler(req, list)
	if sortErr != nil {
		t.Fatal(sortErr)
	}
	for i := 1; i < len(sortedList); i++ {
		assert.LessOrEqual(t, sortedList[i-1].Percentage, sortedList[i].Percentage, "List is not sorted correctly.")

		req, getErr = http.NewRequest("GET", constants.HISTORY_PATH+"?sortalphabetically=true", nil)
		if getErr != nil {
			t.Fatal(err)
		}

		// Checks if list is sorted alphabetically.
		sortedList, sortErr = SortQueryHandler(req, list) // Ascending sorting.
		if sortErr != nil {
			t.Fatal(sortErr)
		}
		for i := 1; i < len(sortedList); i++ {
			if sortedList[i-1].Name < sortedList[i].Name {
				t.Fatal("List is not sorted correctly.")
			}
		}
		// Checks if list is sorted alphabetically descending.
		req, getErr = http.NewRequest("GET", constants.HISTORY_PATH+"?sortalphabetically=true&descending=true", nil)
		if getErr != nil {
			t.Fatal(err)
		}

		// Checks if list is sorted alphabetically.
		sortedList, sortErr = SortQueryHandler(req, list) // Ascending sorting.
		if sortErr != nil {
			t.Fatal(sortErr)
		}
		for i := 1; i < len(sortedList); i++ {
			if sortedList[i-1].Name > sortedList[i].Name {
				t.Fatal("List is not sorted correctly.")
			}
		}
	}
}
