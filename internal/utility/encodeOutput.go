package utility

import (
	"encoding/json"
	"net/http"
)

// Encoder Encodes into pretty json, writes to client.
// Function retrieved from Sanders' assignment 1.
func Encoder(w http.ResponseWriter, template interface{}) {
	enc := json.NewEncoder(w)           //Creates a new encoder to be sent to the client.
	enc.SetIndent("", "\t")             //Prints a formatted json.
	encodeError := enc.Encode(template) //Actual encoding.
	if encodeError != nil {
		http.Error(w, "Error when encoding to client: "+encodeError.Error(), http.StatusUnprocessableEntity)
		return
	}
}
