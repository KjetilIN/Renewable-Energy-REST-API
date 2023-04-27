package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// Define file paths
var htmlFilePath = "templates/default.html"
var cssFilePath = "templates/style.css"

func TestLoadFile(t *testing.T) {
	// Changes the working directory to the project directory.
	changeErr := utility.DirChanger(2)
	if changeErr != nil {
		t.Fatal(changeErr.Error())
	}

	// Test loading HTML and Css file that exists
	_, err := loadFile(htmlFilePath)
	if err != nil {
		t.Errorf("loadFile returned error when loading existing file: %v", err)
	}

	_, err = loadFile(cssFilePath)
	if err != nil {
		t.Errorf("loadFile returned error when loading existing file: %v", err)
	}

	// Test loading a file that doesn't exist
	_, err = loadFile("testdata/nonexistent.txt")
	if err == nil {
		t.Error("loadFile did not return an error when loading nonexistent file")
	}
}

func TestHandlerDefault(t *testing.T) {
	// Changes the working directory to the project directory.
	changeErr := utility.DirChanger(2)
	if changeErr != nil {
		t.Fatal(changeErr.Error())
	}

	// Load HTML and CSS files
	html, htmlErr := os.ReadFile(htmlFilePath)
	if htmlErr != nil {
		t.Fatal(htmlErr)
	}
	css, cssErr := os.ReadFile(cssFilePath)
	if cssErr != nil {
		t.Fatal(cssErr)
	}

	// Parse HTML and add CSS styles
	doc, docErr := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if docErr != nil {
		t.Fatal(docErr)
	}
	style := fmt.Sprintf("<style>%s</style>", string(css))
	doc.Find("head").AppendHtml(style)

	// Create request and response recorder
	req, reqErr := http.NewRequest("GET", constants.DEFAULT_PATH, nil)
	if reqErr != nil {
		t.Fatal(reqErr)
	}
	rr := httptest.NewRecorder()

	// Call handler function
	handler := http.HandlerFunc(HandlerDefault)
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check content type
	expectedContentType := "text/html; charset=utf-8"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}

	// Check body contains expected substrings
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
