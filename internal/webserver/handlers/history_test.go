package handlers

import (
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/structs"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"strconv"
	"testing"
)

const ORIGINAL_LENGTH = 1
const SHORTER = 2

// dirChanger Changes the directory to project root.
func dirChanger() error {
	// Gets the filepath of history_test.go.
	_, filename, _, _ := runtime.Caller(0)
	// Jumps back 3 folders.
	dir := path.Join(path.Dir(filename), "..", "..", "..")
	// Changes to the new dir structure.
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	return nil
}

// getBody a function which decodes body into a template.
func getBody(response *http.Response, template interface{}) error {
	body, ioReadErr := io.ReadAll(response.Body)
	if ioReadErr != nil {
		return ioReadErr
	}
	json.Unmarshal(body, template)
	return nil
}

// TestHandlerHistory_NoParams Testing the base return from history endpoint.
func TestHandlerHistory_NoParams(t *testing.T) {
	// Changes the working directory to the project directory.
	changeErr := dirChanger()
	if changeErr != nil {
		t.Fatal(changeErr.Error())
	}
	// Creates a test server on handler history.
	server := httptest.NewServer(http.HandlerFunc(HandlerHistory))
	resp, getReqErr := http.Get(server.URL)
	if getReqErr != nil {
		t.Fatal("Error when requesting: " + getReqErr.Error())
	}
	var testList []structs.RenewableShareEnergyElement
	err := getBody(resp, &testList)
	if err != nil {
		t.Fatal("Error when getting body: " + err.Error())
	}
	// Waits for the body to close.
	defer resp.Body.Close()
	// Checks if the request is of status ok.
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Handler returned wrong status code.")
	// Checks if json from body contains anything.
	assert.NotEmpty(t, testList, "JSON list from body is empty.")
}

// prepareList Function which prepares the lists for testing.
func prepareList(setting int) ([]structs.RenewableShareEnergyElement, error) {
	// Changes the working directory to the project directory.
	changeErr := dirChanger()
	if changeErr != nil {
		return nil, changeErr
	}
	// Collects CSV file into a list of JSON structs.
	list, err := utility.RSEToJSON()
	if err != nil {
		return nil, err
	}
	// Switch case for different settings of the list. For now only original and shorter setting.
	switch setting {
	case SHORTER: // Returns the first 100 elements from original list.
		return list[0:100], nil
	default:
		return list, nil
	}
}

// TestCountryCodeLimiter Test function which tests the country code limiter from history endpoint.
func TestCountryCodeLimiter(t *testing.T) {
	// Constant for testing.
	testCountryCode := "SWE"
	// Retrieves the original list from CSV.
	originalList, err := prepareList(ORIGINAL_LENGTH)
	if err != nil {
		t.Fatal("Error when getting list, " + err.Error())
	}
	// Calls the countryCode limiter method.
	filteredList := countryCodeLimiter(originalList, testCountryCode)

	// Checks if all country codes is the same in the return list.
	for i := 0; i < len(filteredList); i++ {
		assert.Equal(t, filteredList[i].IsoCode, "SWE", "Country codes does not match.")
	}
	// Makes sure the list is not empty.
	assert.NotEmpty(t, filteredList, "List is empty.")
}

// TestBeginEndQuery Test function which tests the queries for getting data between certain years.
func TestBeginEndQuery(t *testing.T) {
	// Constants for testing.
	begin := 2010
	end := 2011

	originalList, err := prepareList(ORIGINAL_LENGTH)
	filteredList, err := beginEndLimiter(strconv.Itoa(begin), strconv.Itoa(end), originalList)
	if err != nil {
		t.Fatal("Error when filtering list: " + err.Error())
	}

	// Checks if year is between the specified in query.
	for i := 0; i < len(filteredList); i++ {
		assert.GreaterOrEqual(t, end, filteredList[i].Year, "Year is lower than begin query.")
		assert.LessOrEqual(t, begin, filteredList[i].Year, "Year is greater than end query.")
	}
}

func TestBeginQuery(t *testing.T) {
	// Constants for testing.
	begin := 2020
	originalList, err := prepareList(ORIGINAL_LENGTH)
	filteredList, err := beginEndLimiter(strconv.Itoa(begin), "", originalList)
	if err != nil {
		t.Fatal("Error when filtering list: " + err.Error())
	}

	// Checks if year is between the specified in query.
	for i := 0; i < len(filteredList); i++ {
		assert.LessOrEqual(t, begin, filteredList[i].Year, "Year is lower than begin query.")
	}
}

func TestEndQuery(t *testing.T) {
	// Constants for testing.
	end := 1980
	originalList, err := prepareList(ORIGINAL_LENGTH)
	filteredList, err := beginEndLimiter("", strconv.Itoa(end), originalList)
	if err != nil {
		t.Fatal("Error when filtering list: " + err.Error())
	}
	// Checks if year is between the specified in query.
	for i := 0; i < len(filteredList); i++ {
		assert.GreaterOrEqual(t, end, filteredList[i].Year, "Year is lower than begin query.")
	}
}

// TestSortedByPercentage Testing the sorting function from history.
func TestSortedByPercentage(t *testing.T) {
	// Prepares a shorter list for testing. This is due to the sorting method being slow.
	list, err := prepareList(ORIGINAL_LENGTH)
	if err != nil {
		t.Fatal("Error when getting list: " + err.Error())
	}
	// Sorts the list by percentage.
	sortedList := sliceSortingByValue(list, 1) // Ascending sorting.

	// Checks if list is sorted by percentage.
	for i := 1; i < len(sortedList); i++ {
		assert.GreaterOrEqualf(t, sortedList[i-1].Percentage, sortedList[i].Percentage, "List is not sorted.")
	}

	sortedList = sliceSortingByValue(list, 2) // Descending value.
	for i := 1; i < len(sortedList); i++ {
		assert.LessOrEqualf(t, sortedList[i-1].Percentage, sortedList[i].Percentage, "List is not sorted correctly.")
	}
}

// TestMeanCalculated Tests the mean list.
func TestMeanCalculated(t *testing.T) {
	// Gets the shorter list to calculate the first country's mean.
	shortList, err := prepareList(SHORTER)
	if err != nil {
		t.Fatal("Error getting list: " + err.Error())
	}
	// Calculated mean.
	meanAfrica := 7.436156868
	meanList := meanCalculation(shortList)

	// Checks if the average of first country is correct.
	assert.Equal(t, math.Round(meanAfrica), math.Round(meanList[0].Percentage), "The average is wrong.")
}
