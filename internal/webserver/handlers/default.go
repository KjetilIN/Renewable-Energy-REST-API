package handlers

import (
	"assignment-2/internal/constants"
	"fmt"
	"log"
	"net/http"
)

// HandlerDefault is a handler for the /default endpoint.
func HandlerDefault(w http.ResponseWriter, r *http.Request) {
	// Ensure interpretation as HTML by client (browser)
	w.Header().Set("content-type", "text/html")

	// Offer information for redirection to paths
	output := "<h1>Welcome!<h1><h3>This service does not provide any functionality on root path level." +
		" Please try one of the paths below<h3>" +
		"<h5 style=\"background-color: lightblue; width: 250px;\">Current percentage of renewables:<br>" +
		"<a href=\"" + constants.CURRENT_PATH + "\">" + constants.CURRENT_PATH + "</a></h5>" +
		"<h5 style=\"background-color: lightblue; width: 250px;\">Historical percentages of renewables:<br>" +
		"<a href=\"" + constants.HISTORY_PATH + "\">" + constants.HISTORY_PATH + "</a></h5>" +
		"<h5 style=\"background-color: lightblue; width: 250px;\">Notification for webhooks:<br>" +
		"<a href=\"" + constants.NOTIFICATIONS_PATH + "\">" + constants.NOTIFICATIONS_PATH + "</a></h5>" +
		"<h5 style=\"background-color: lightblue; width: 250px;\">For status:<br>" +
		"<a href=\"" + constants.STATUS_PATH + "\">" + constants.STATUS_PATH + "</a></h5>"

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when returning output.")
	}
}
