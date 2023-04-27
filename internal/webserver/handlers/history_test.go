package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/mock"
	"assignment-2/internal/webserver/structs"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const ORIGINAL_LENGTH = 1
const SHORTER = 2

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
	err := utility.DirChanger(3)
	if err != nil {
		return
	}

	// Creates a new HTTP request
	req, err := http.NewRequest("GET", constants.MOCK_HISTORY_API_URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Creates a new response recorder to record the response from the server
	rr := httptest.NewRecorder()

	// Creates a stub HTTP handler
	handler := http.HandlerFunc(mock.StubHandlerHistory)

	// Sends the request to the server and records the response
	handler.ServeHTTP(rr, req)

	// Checks if the request is of status ok.
	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code.")

	// Checks if json from body contains anything.
	var testList []structs.RenewableShareEnergyElement
	err = json.Unmarshal(rr.Body.Bytes(), &testList)
	if err != nil {
		t.Fatal("Error when unmarshalling body: " + err.Error())
	}
	assert.NotEmpty(t, testList, "JSON list from body is empty.")
}

// TestHistoryMockHandler tests GET and POST requests on the History mock handler
func TestHistoryMockHandler(t *testing.T) {
	// Changes the working directory to the project directory.
	err := utility.DirChanger(3)
	if err != nil {
		return
	}

	// Testing a get request on local host
	getRequest, _ := http.NewRequest("GET", constants.MOCK_HISTORY_API_URL, nil)
	response := httptest.NewRecorder()
	//Executing the handler
	mock.StubHandlerHistory(response, getRequest)
	resultGet := response.Result()
	defer resultGet.Body.Close()

	//Error if not implemented or not correct
	if resultGet.StatusCode != http.StatusOK {
		t.Error("Test case on GET failed, should be 200")
	}
	expected1 := "application/json"
	resultGetHeader := resultGet.Header.Get("content-type")
	if resultGetHeader != expected1 {
		t.Errorf("Test case failed on GET: wrong header information")
	}

	// Test case 2: POST request
	postRequest, _ := http.NewRequest("POST", constants.MOCK_HISTORY_API_URL, nil)
	postResponse := httptest.NewRecorder()
	mock.StubHandlerHistory(postResponse, postRequest)
	resultPost := postResponse.Result()
	defer resultPost.Body.Close()
	if resultPost.StatusCode != http.StatusNotImplemented {
		t.Errorf("Test case POST failed: Not marked as not implemented")
	}
}

// prepareList Function which prepares the lists for testing.
func prepareList(setting int) ([]structs.RenewableShareEnergyElement, error) {
	// Changes the working directory to the project directory.
	changeErr := utility.DirChanger(2)
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
		return list[0:57], nil
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

	// Creates list and filters it.
	originalList, err := prepareList(SHORTER)
	filteredList, err := beginEndLimiter(strconv.Itoa(begin), strconv.Itoa(end), originalList)
	if err != nil {
		t.Fatal("Error when filtering list: " + err.Error())
	}
	// Calculated by hand.
	calculatedMeanAfrica := 7.1637459

	// Checks if mean is correct for the test years.
	for i := 0; i < len(filteredList); i++ {
		assert.Equal(t, calculatedMeanAfrica, filteredList[0].Percentage)
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

// TestSorted Testing the sorting function from history.
func TestSorted(t *testing.T) {
	// Prepares a list for testing.
	list, err := prepareList(ORIGINAL_LENGTH)
	if err != nil {
		t.Fatal("Error when getting list: " + err.Error())
	}
	// Sorts the list by percentage.
	sortedList := utility.SortRSEList(list, false, 1) // Ascending sorting.

	// Checks if list is sorted by percentage.
	for i := 1; i < len(sortedList); i++ {
		assert.LessOrEqualf(t, sortedList[i-1].Percentage, sortedList[i].Percentage, "List is not sorted.")
	}
	sortedList = utility.SortRSEList(list, false, 2) // Descending value.
	for i := 1; i < len(sortedList); i++ {
		assert.GreaterOrEqual(t, sortedList[i-1].Percentage, sortedList[i].Percentage, "List is not sorted correctly.")
	}

	// Checks if list is sorted alphabetically.
	sortedList = utility.SortRSEList(list, true, 1) // Ascending sorting.
	for i := 1; i < len(sortedList); i++ {
		if sortedList[i-1].Name < sortedList[i].Name {
			t.Fatal("List is not sorted correctly.")
		}
	}

	// Checks if list is sorted descending alphabetically.
	sortedList = utility.SortRSEList(list, true, 2) // Descending sorting.
	for i := 1; i < len(sortedList); i++ {
		if sortedList[i-1].Name > sortedList[i].Name {
			t.Fatal("List is not sorted correctly.")
		}
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
	meanAfrica := 7.436156868421055
	meanList := meanCalculation(shortList)
	// As mean uses maps, the return is unsorted.
	meanList = utility.SortRSEList(meanList, false, 1)

	// Checks if the average of first country is correct.
	assert.Equal(t, meanAfrica, meanList[0].Percentage, "The average is wrong.")
}
