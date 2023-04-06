package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

// BeforeEach Function to run before all test functions.
// Creates a new server to be tested on, using default handler.
func BeforeEach(query string) (*http.Response, error) {
	// Changes to project
	// TODO: Make it work with every pc.
	changeErr := os.Chdir("C:\\Users\\sande\\GolandProjects\\assignment-2")
	if changeErr != nil {
		return nil, changeErr
	}
	// Creates a test server on handler history.
	server := httptest.NewServer(http.HandlerFunc(HandlerDefault))
	server.URL = server.URL + query
	fmt.Println(server.URL)
	resp, getReqErr := http.Get(server.URL + query)
	if getReqErr != nil {
		return nil, getReqErr
	}
	return resp, nil
}

// getBody a function which decodes body into a template.
func getBody(response *http.Response, template interface{}) error {
	body, ioReadErr := io.ReadAll(response.Body)
	if ioReadErr != nil {
		return ioReadErr
	}
	unmarshallErr := json.Unmarshal(body, &template)
	if unmarshallErr != nil {
		return unmarshallErr
	}
	return nil
}

// TestHandlerHistory_NoParams Testing the base return from history endpoint.
func TestHandlerHistory_NoParams(t *testing.T) {
	resp, error := BeforeEach(constants.HISTORY_PATH)
	if error != nil {
		t.Fatal(error.Error())
	}
	// Waits for the body to close.
	defer resp.Body.Close()
	// Checks if the request is of status ok.
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Handler returned wrong status code.")
}

// TestCountryCodeLimiter Test function which tests the country code limiter from history endpoint.
func TestCountryCodeLimiter(t *testing.T) {
	// Retrieves a response from a newly created server.
	resp, err := BeforeEach(constants.HISTORY_PATH + "/nor")
	if err != nil {
		t.Fatal(err.Error())
	}

	var list []structs.HistoricalRSE
	ioReadErr := getBody(resp, list)
	if ioReadErr != nil {
		t.Fatal("Body read error: " + ioReadErr.Error())
	}
	// Waits for the body to close.
	defer resp.Body.Close()

	// Checks if all country codes is the same in the return list.
	for i := 1; i < len(list); i++ {
		assert.Equal(t, list[i-1].IsoCode, list[i].IsoCode, "Country codes does not match.")
	}
}

// TestBeginEndQuery Test function which tests the queries for getting data between certain years.
func TestBeginEndQuery(t *testing.T) {
	// Constants for testing.
	begin := 2010
	end := 2011
	// Retrieves a response from a newly created server.
	resp, err := BeforeEach(constants.HISTORY_PATH + "?begin=" + strconv.Itoa(begin) + "&end=" + strconv.Itoa(end))
	if err != nil {
		t.Fatal(err.Error())
	}
	var list []structs.HistoricalRSE
	// Unmarshalls the response body into the list using referencing.
	ioReadErr := getBody(resp, list)
	if ioReadErr != nil {
		t.Fatal("Body read error: " + ioReadErr.Error())
	}
	// Waits for the body to close.
	defer resp.Body.Close()

	// Checks if year is between the specified in query.
	for i := 0; i < len(list); i++ {
		assert.GreaterOrEqualf(t, begin, list[i].Year, "Year is lower than begin query.")
		assert.Less(t, end, list[i].Year, "Year is greater than end query.")
	}
}

func TestBeginQuery(t *testing.T) {

}
func TestEndQuery(t *testing.T) {

}

func TestSortedByPercentage(t *testing.T) {

}

func TestMeanCalculated(t *testing.T) {

}
