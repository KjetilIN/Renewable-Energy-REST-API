package handlers

import (
	"assignment-2/internal/constants"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const URL = "http://localhost:" + constants.DEFAULT_PORT

// TestHandlerHistory_NoParams Testing the base return from history endpoint.
func TestHandlerHistory_NoParams(t *testing.T) {
	// Creates a new request for history endpoint.
	req, err := http.NewRequest("GET", URL+constants.HISTORY_PATH, nil)
	if err != nil {
		t.Error("Request error: " + err.Error())
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerHistory)

	handler.ServeHTTP(response, req)
	assert.Equal(t, nil, response.Result(), "Handler returned wrong status code.")
}

func TestCountryCodeLimiter(t *testing.T) {

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
