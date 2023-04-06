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
	"testing"
)

const URL = "http://localhost:" + constants.DEFAULT_PORT

func BeforeEach(query string) (*http.Response, error) {
	// Changes to project
	// TODO: Make it work with every pc.
	changeErr := os.Chdir("C:\\Users\\sande\\GolandProjects\\assignment-2")
	if changeErr != nil {
		return nil, changeErr
	}
	// Creates a test server on handler history.
	server := httptest.NewServer(http.HandlerFunc(HandlerHistory))
	server.URL = server.URL + query
	fmt.Println(server.URL)
	resp, getReqErr := http.Get(server.URL + query)
	if getReqErr != nil {
		return nil, getReqErr
	}
	return resp, nil
}

// TestHandlerHistory_NoParams Testing the base return from history endpoint.
func TestHandlerHistory_NoParams(t *testing.T) {
	resp, error := BeforeEach("")
	if error != nil {
		t.Fatal(error.Error())
	}
	// Waits for the body to close.
	defer resp.Body.Close()
	// Checks if the request is of status ok.
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Handler returned wrong status code.")
}

func TestCountryCodeLimiter(t *testing.T) {
	resp, error := BeforeEach(constants.HISTORY_PATH + "/nor")
	if error != nil {
		t.Fatal(error.Error())
	}
	body, ioReadErr := io.ReadAll(resp.Body)
	if ioReadErr != nil {
		t.Fatal("Body read error: " + ioReadErr.Error())
	}
	var list []structs.HistoricalRSE
	json.Unmarshal(body, &list)

	// Waits for the body to close.
	defer resp.Body.Close()

	for i := 1; i < len(list); i++ {
		assert.Equal(t, list[i-1].IsoCode, list[i].IsoCode, "Country codes does not match.")
	}
}

func TestBeginEndQuery(t *testing.T) {

}

func TestBeginQuery(t *testing.T) {

}
func TestEndQuery(t *testing.T) {

}

func TestSortedByPercentage(t *testing.T) {

}

func TestMeanCalculated(t *testing.T) {

}
