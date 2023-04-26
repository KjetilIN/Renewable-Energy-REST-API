package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"net/http"
	"strings"
)

// HandlerDefault is a handler for the /default endpoint.
func HandlerDefault(w http.ResponseWriter, _ *http.Request) {
	// Load HTML and CSS files
	html, err := loadFile("templates/default.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when reading HTML file.")
		return
	}
	css, err := loadFile("templates/style.css")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when reading CSS file.")
		return
	}

	// Parse HTML and add CSS styles
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when parsing HTML.")
		return
	}
	style := fmt.Sprintf("<style>%s</style>", css)
	doc.Find("head").AppendHtml(style)

	// Write the modified HTML to the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	docHtml, err := doc.Html()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when converting HTML to string.")
		return
	}

	if _, err := w.Write([]byte(docHtml)); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when writing HTML to response writer.")
		return
	}
}

// loadFile takes a filename as a string and returns the contents of the file as a string.
// Returns: a string, or an error and an empty string.
func loadFile(filename string) (string, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
