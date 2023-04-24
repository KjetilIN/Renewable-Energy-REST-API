package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/mock"
	"assignment-2/internal/webserver/structs"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandlerCurrent_NoParams Tests the base handler without any params.
func TestHandlerCurrent_NoParams(t *testing.T) {
	// Changes working directory to root directory.
	dirChangeErr := utility.DirChanger() // Function in history test.
	if dirChangeErr != nil {
		t.Fatal("Error switching working directory: " + dirChangeErr.Error())
	}
	server := httptest.NewServer(http.HandlerFunc(mock.StubHandlerCurrent))
	resp, getReqErr := http.Get(server.URL)
	if getReqErr != nil {
		t.Fatal("Error when requesting: " + getReqErr.Error())
	}
	var testList []structs.RenewableShareEnergyElement
	err := getBody(resp, &testList) // Function in history test.
	if err != nil {
		t.Fatal("Error when getting body: " + err.Error())
	}
	// Waits for the body to close.
	defer resp.Body.Close()

	// Checks if the request is of status ok.
	assert.Equal(t, 200, resp.StatusCode, "Handler returned wrong status code.")
	// Checks if json from body contains anything.
	assert.NotEmpty(t, testList, "JSON list from body is empty.")

	currentYear := getCurrentYear(testList)

	for _, v := range testList {
		if v.Year != currentYear {
			t.Fatal("Value is not of the current year.")
		}
	}
}

// TestCurrentMockHandler tests GET and POST requests on the Current mock handler
func TestCurrentMockHandler(t *testing.T) {
	// Changes the working directory to the project directory.
	err := utility.DirChanger()
	if err != nil {
		return
	}

	// Testing a get request on local host
	getRequest, _ := http.NewRequest("GET", constants.MOCK_CURRENT_API_URL, nil)
	response := httptest.NewRecorder()
	//Executing the handler
	mock.StubHandlerCurrent(response, getRequest)
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
	postRequest, _ := http.NewRequest("POST", constants.MOCK_CURRENT_API_URL, nil)
	postResponse := httptest.NewRecorder()
	mock.StubHandlerCurrent(postResponse, postRequest)
	resultPost := postResponse.Result()
	defer resultPost.Body.Close()
	if resultPost.StatusCode != http.StatusNotImplemented {
		t.Errorf("Test case POST failed: Not marked as not implemented")
	}
}

// TestNeighbourRetrieval Tests if neighbour retrieval works.
// Tests API retrieval at the same time.
func TestNeighbourRetrieval(t *testing.T) {
	// Changes directory.
	dirChangeErr := utility.DirChanger() // Function in history test.
	if dirChangeErr != nil {
		t.Fatal("Error changing directory: " + dirChangeErr.Error())
	}
	countryCodeTest := "NOR"
	origList, err := utility.RSEToJSON()
	if err != nil {
		t.Fatal("Error when getting list: " + err.Error())
	}
	// Retrieves the list of current year elements.
	currentList := getCurrentList(origList)
	// Checks if list is empty.
	assert.NotEmpty(t, currentList, "Current list is empty, and is not supposed to be.")

	// Retrieves the neighbouring countries of the test country.
	neighbours, err := retrieveNeighbours(currentList, countryCodeTest)
	if err != nil {
		t.Fatal("Error when retrieving neighbouring countries.")
	}
	assert.NotEmpty(t, neighbours, "Neighbour list is empty, it should not be empty.")

	// Expected borders presented.
	expectedBorders := []string{"FIN", "SWE", "RUS"}

	// Checks if the iso codes matched expected iso codes.
	for i := 0; i < len(neighbours); i++ {
		assert.Contains(t, neighbours[i].IsoCode, expectedBorders[i], "Borders does not correspond.")
	}

	// Negative test.
	_, expectedError := retrieveNeighbours(currentList, "COUNTRY_THAT_DOES_NOT_EXIST")
	assert.Error(t, expectedError, "Error was not returned.")
}
