package handlers

import (
	"encoding/csv"
	"net/http"
	"os"
)

// HandlerHistory is a handler for the /history endpoint.
func HandlerHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

}

// Function to read from a CSV file.
func readCSV(filePath string) ([][]string, error) {
	file, readErr := os.Open(filePath)
	if readErr != nil {
		return nil, readErr
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	information, parseError := csvReader.ReadAll()
	if parseError != nil {
		return nil, parseError
	}
	return information, nil
}
