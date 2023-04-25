package handlers

import (
	"assignment-2/internal/constants"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerDefault(t *testing.T) {
	// Changes the working directory to the project directory.
	changeErr := dirChanger()
	if changeErr != nil {
		t.Fatal(changeErr.Error())
	}

	req, err := http.NewRequest("GET", constants.DEFAULT_PATH, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerDefault)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedContentType := "text/html; charset=utf-8"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}

	expectedSubstring := "<div class=\"project-info\">"
	if body := rr.Body.String(); !strings.Contains(body, expectedSubstring) {
		t.Errorf("handler returned unexpected body: missing expected substring %v", expectedSubstring)
	}

	expectedSubstring2 := "<section class=\"endpoints\">"
	if body := rr.Body.String(); !strings.Contains(body, expectedSubstring2) {
		t.Errorf("handler returned unexpected body: missing expected substring %v", expectedSubstring2)
	}

	expectedLink := `<link rel="stylesheet" type="text/css" href="/style.css"/>`
	if body := rr.Body.String(); !strings.Contains(body, expectedLink) {
		t.Errorf("handler returned unexpected body: missing expected link tag")
	}

	expectedTitle := "<title>Renewable Energy REST Web Application</title>"
	if body := rr.Body.String(); !strings.Contains(body, expectedTitle) {
		t.Errorf("handler returned unexpected body: missing expected title tag")
	}

	expectedHeader := "<header>Welcome to the Renewable Energy REST Web Application!</header>"
	if body := rr.Body.String(); !strings.Contains(body, expectedHeader) {
		t.Errorf("handler returned unexpected body: missing expected header tag")
	}
}
