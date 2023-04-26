package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"net/http"
	"strings"
)

// HandlerDefault is a handler for the /default endpoint.
func HandlerDefault(w http.ResponseWriter, r *http.Request) {
	// Query for printing information about endpoint.
	if r.URL.Query().Has("information") && strings.Contains(strings.ToLower(r.URL.Query().Get("information")), "true") {
		_, writeErr := w.Write([]byte("To use API, remove ?information=true, from the URL.\n"))
		if writeErr != nil {
			return
		}
		utility.Encoder(w, constants.DEFAULT_QUERIES)
		return
	}
}
