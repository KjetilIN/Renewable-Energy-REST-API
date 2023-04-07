package handlers

import (
	"assignment-2/internal/constants"
	"html/template"
	"log"
	"net/http"
)

// DefaultPageData contains data for the default page.
type DefaultPageData struct {
	CurrentPath       string
	HistoryPath       string
	NotificationsPath string
	StatusPath        string
}

// HandlerDefault is a handler for the /default endpoint.
func HandlerDefault(w http.ResponseWriter, r *http.Request) {
	// Ensure interpretation as HTML by client (browser)
	w.Header().Set("content-type", "text/html")

	// Define the page data.
	pageData := DefaultPageData{
		CurrentPath:       constants.CURRENT_PATH,
		HistoryPath:       constants.HISTORY_PATH,
		NotificationsPath: constants.NOTIFICATIONS_PATH,
		StatusPath:        constants.STATUS_PATH,
	}

	// Parse the template.
	tmpl, err := template.ParseFiles("templates/default.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when parsing template.")
		return
	}

	// Render the template with the page data.
	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when rendering template.")
	}
}
