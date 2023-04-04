package handlers

import (
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/structs"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// HandlerHistory is a handler for the /history endpoint.
func HandlerHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	// Reads from csv and returns json list.
	listOfRSE, jsonError := rseToJSON()
	if jsonError != nil {
		http.Error(w, jsonError.Error(), http.StatusInternalServerError)
		fmt.Errorf("%s", jsonError.Error())
		return
	}

	// Collects parameters, separated by /
	params := strings.Split(r.URL.Path, "/") //Used to split the / in path to collect search parameters.

	// Checks if an optional parameter is passed.
	if len(params) == 6 {
		listOfRSE = countryCodeLimiter(listOfRSE, params[5])
	}

	// Checks for queries.
	if r.URL.Query().Has("begin") || r.URL.Query().Has("end") {
		var queryError error // Initialises a potential error.
		beginQuery := r.URL.Query().Get("begin")
		endQuery := r.URL.Query().Get("end")
		// Calls function to include begin and end checking.
		listOfRSE, queryError = beginEndLimiter(beginQuery, endQuery, listOfRSE)
		if queryError != nil {
			http.Error(w, "Error using queries: "+queryError.Error(), http.StatusBadRequest)
		}
	}
	// Checks if sortByValue query is passed.
	if r.URL.Query().Has("sortbyvalue") && strings.Contains(strings.ToLower(r.URL.Query().Get("sortbyvalue")), "true") {
		listOfRSE = sortingListPercentage(listOfRSE)
	}

	// Checks if list is empty.
	if len(listOfRSE) == 0 {
		http.Error(w, "Nothing matching your search terms.", http.StatusBadRequest)
		return
	}

	// If Query: mean=true, a different struct type will be encoded to client. It calculates the mean of grouped countries.
	if r.URL.Query().Has("mean") && strings.Contains(strings.ToLower(r.URL.Query().Get("mean")), "true") {
		meanList := meanCalculation(listOfRSE)
		utility.Encoder(w, meanList)
	} else {
		// Encodes list and prints to console.
		utility.Encoder(w, listOfRSE)
	}
}

// rseToJSON is an internal function to use a 2D string and input it into a struct.
func rseToJSON() ([]structs.HistoricalRSE, error) {
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

// countryCodeLimiter Method to limit a list based on country code.
func countryCodeLimiter(listToIterate []structs.HistoricalRSE, countryCode string) []structs.HistoricalRSE {
	var limitedList []structs.HistoricalRSE
	for i, v := range listToIterate { // Iterates through input list.
		if strings.Contains(strings.ToLower(listToIterate[i].IsoCode), countryCode) { // If country code match it is
			// appended to new list.
			limitedList = append(limitedList, v)
		}
	}
	return limitedList // Returns list containing all matching countries.
}

// beginEndLimiter Function to allow for searching to and from a year.
func beginEndLimiter(begin string, end string, listToIterate []structs.HistoricalRSE) ([]structs.HistoricalRSE, error) {
	var newlist []structs.HistoricalRSE
	var convErr error // Potential error.
	var convBegin int // Variable to store str turned to int.
	var convEnd int   // Variable to store str turned to int.
	toFromOr := 0     // Functions as a boolean.

	// Switch case to make it possible to check for begin and end, or just begin/end.
	switch {
	case len(begin) > 0 && len(end) > 0: // Both begin and end exists.
		toFromOr = 3
		convBegin, convErr = strconv.Atoi(begin)
		convEnd, convErr = strconv.Atoi(end)
	case len(begin) > 0: // Only begin exists.
		toFromOr = 1
		convBegin, convErr = strconv.Atoi(begin)
	case len(end) > 0: // Only end exists.
		toFromOr = 2
		convEnd, convErr = strconv.Atoi(end)
	}
	// If a conversion error occurred.
	if convErr != nil {
		return nil, convErr
	}
	// Append json objects fitting conditions to newlist.
	for i, v := range listToIterate {
		relevantYear := listToIterate[i].Year
		if toFromOr == 3 && relevantYear <= convEnd && convBegin <= relevantYear {
			newlist = append(newlist, v)
		} else if toFromOr == 1 && convBegin <= relevantYear {
			newlist = append(newlist, v)
		} else if toFromOr == 2 && relevantYear <= convEnd {
			newlist = append(newlist, v)
		}
	}
	return newlist, nil
}

// meanCalculation Function to calculate the mean of percentage per country, from the inputted list.
func meanCalculation(listToIterate []structs.HistoricalRSE) []structs.HistoricalRSEMean {
	newList := []structs.HistoricalRSEMean{}
	meanList := []float64{} // Initiates an empty float slice.
	sum, mean := 0.0, 0.0
	// Iterates through input list to calculate mean.
	for i := 1; i < len(listToIterate); i++ {
		if listToIterate[i].Name == listToIterate[i-1].Name { // If name is the same as previous, add value to meanList.
			meanList = append(meanList, listToIterate[i].Percentage)
		} else { // If it is not the same, we have jumped to a new country. Then the mean should be calculated.
			// Add up all floats.
			for _, v := range meanList {
				sum = sum + v
			}
			mean = sum / float64(len(meanList))

			// Potential bug: duplicate names and isocode.
			newEntry := structs.HistoricalRSEMean{
				Name:       listToIterate[i-1].Name,
				IsoCode:    listToIterate[i-1].IsoCode,
				Percentage: mean,
			}
			// Resets the lists and variables.
			newList = append(newList, newEntry)
			mean, sum = 0.0, 0.0
			meanList = []float64{}
		}
	}
	return newList
}

// sortingListPercentage a function which sorts a json list based on percentage. The function is not very
// efficient.
func sortingListPercentage(listToIterate []structs.HistoricalRSE) []structs.HistoricalRSE {
	var sortedList []structs.HistoricalRSE
	HighestValIndex := 0
	HighestVal := 0.0
	sorted := false
	count := 0

	// Loop which iterates until sorted is true.
	for !sorted {
		// Iterates through all elements in listToIterate.
		for i, v := range listToIterate {
			// If the current percentage is highest.
			if v.Percentage > HighestVal {
				HighestVal = v.Percentage
				HighestValIndex = i
			}
			// Checks if i is at the end of the list.
			if i == len(listToIterate)-1 {
				sortedList = append(sortedList, listToIterate[HighestValIndex])
				// Resets values for another loop.
				listToIterate[HighestValIndex].Percentage = 0.0
				HighestVal = 0
				HighestValIndex = 0
			}
		}
		// Counts amount of times iterated through list.
		count = count + 1
		// If count is as long as the passed list, the sorting is done.
		if count == len(listToIterate) {
			sorted = true
		}
	}
	return sortedList
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
