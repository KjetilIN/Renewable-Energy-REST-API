package handlers

import (
	"assignment-2/internal/constants"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const URL = "http://localhost:" + constants.DEFAULT_PORT

// TestHandlerHistory_NoParams Testing the base return from history endpoint.
func TestHandlerHistory_NoParams(t *testing.T) {
	// Changes to project
	// TODO: Make it work with every pc.
	changeErr := os.Chdir("C:\\Users\\sande\\GolandProjects\\assignment-2")
	fmt.Println(os.Getwd())
	if changeErr != nil {
		t.Error(changeErr)
	}
	// Creates a test server on handler history.
	server := httptest.NewServer(http.HandlerFunc(HandlerHistory))
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal("Something went wrong: " + err.Error())
	}
	// Checks if the request is of status ok.
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Handler returned wrong status code.")
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
