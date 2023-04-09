package handlers

import (
	"assignment-2/internal/webserver/structs"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"
)

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

// TestHandlerCurrent_NoParams Tests the base handler without any params.
func TestHandlerCurrent_NoParams(t *testing.T) {
	// Changes working directory to root directory.
	dirChangeErr := dirChanger()
	if dirChangeErr != nil {
		t.Fatal("Error switching working directory: " + dirChangeErr.Error())
	}
	server := httptest.NewServer(http.HandlerFunc(HandlerCurrent))
	resp, getReqErr := http.Get(server.URL)
	if getReqErr != nil {
		t.Fatal("Error when requesting: " + getReqErr.Error())
	}
	var testList []structs.HistoricalRSE
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

	currentYear := getCurrentYear(testList)

	for _, v := range testList {
		if v.Year != currentYear {
			t.Fatal("Value is not of the current year.")
		}
	}
}
