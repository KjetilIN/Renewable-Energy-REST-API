package utility

import (
	"assignment-2/internal/webserver/structs"
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

// rseToJSON is an internal function to use a 2D string and input it into a struct.
func RSEToJSON() ([]structs.HistoricalRSE, error) {
	var jsonList []structs.HistoricalRSE
	var jsonObj structs.HistoricalRSE

	// readFromFile is a 2D string array.
	readFromFile, readErr := readCSV("./internal/res/renewable-share-energy.csv")
	if readErr != nil {
		return nil, readErr
	}
	var lineRead []string
	//Iterates through 1 dimension of readFromFile.
	for i := 1; i < len(readFromFile); i++ {
		// Stores a slice of values to be iterated through.
		lineRead = readFromFile[i]

		// Parses year from JSON to int, if failed error is handled.
		year, convErr := strconv.Atoi(lineRead[2]) // Converts string line to integer.
		if convErr != nil {
			log.Fatal(convErr)
			return nil, convErr
		}
		// Parses percentage from JSON to float, if failed error is handled.
		percentage, convErr := strconv.ParseFloat(lineRead[3], 6) // Converts string line to float og 6 decimals.
		if convErr != nil {
			log.Fatal(convErr)
			return nil, convErr
		}
		// Iterates through the lineRead slice, and appends to a new entity in HistoricalRSE slice.
		jsonObj = structs.HistoricalRSE{
			Name:       lineRead[0],
			IsoCode:    lineRead[1],
			Year:       year,
			Percentage: percentage,
		}
		jsonList = append(jsonList, jsonObj)

	}
	return jsonList, nil
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
