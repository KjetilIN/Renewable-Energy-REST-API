package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/structs"
	"net/http"
	"strconv"
	"strings"
)

// HandlerHistory is a handler for the /history endpoint.
func HandlerHistory(w http.ResponseWriter, r *http.Request) {
	// Query for printing information about endpoint.
	if r.URL.Query().Has("information") && strings.Contains(strings.ToLower(r.URL.Query().Get("information")), "true") {
		_, writeErr := w.Write([]byte("To use API, remove ?information=true, from the URL.\n"))
		if writeErr != nil {
			return
		}
		utility.Encoder(w, constants.HISTORY_QUERIES)
		return
	}

	// Runs initialise method for handler.
	listOfRSE, initError := InitHandler(w, r)
	if initError != nil {
		return
	}

	// Boolean if all countries are to be shown.
	allCountries := true

	// Collects parameter from url path. If empty, an empty string is returned.
	countryIdentifier := utility.GetParams(r.URL.Path, constants.HISTORY_PATH)
	if countryIdentifier != "" {
		var filterErr error

		listOfRSE, filterErr = CountryFilterer(w, listOfRSE, countryIdentifier)
		if filterErr != nil {
			return
		}
		// It will not show all countries.
		allCountries = false
	}

	// Checks for begin and end queries.
	if r.URL.Query().Has("begin") || r.URL.Query().Has("end") {
		var queryError error // Initialises a potential error.
		beginQuery := r.URL.Query().Get("begin")
		endQuery := r.URL.Query().Get("end")
		// Checks if queries is empty.
		if beginQuery != "" || endQuery != "" {
			// Calls function to include begin and end checking.
			listOfRSE, queryError = beginEndLimiter(beginQuery, endQuery, listOfRSE)
			if queryError != nil {
				http.Error(w, "Begin or end query faulty. It should be of type number.", http.StatusBadRequest)
				return
			}
			// Mean of each country should not be calculated.
			allCountries = false
		} else { // If year is empty, it tells user it is empty, but still allows for code to run.
			http.Error(w, "Begin or end query should not be empty.", http.StatusLengthRequired)
			return
		}
	}

	// Year query, which returns a specific year.
	if r.URL.Query().Has("year") {
		year := r.URL.Query().Get("year")
		var queryErr error
		if year != "" {
			listOfRSE, queryErr = beginEndLimiter(year, year, listOfRSE)
			if queryErr != nil {
				http.Error(w, "Year query faulty, it should be a number.", http.StatusBadRequest)
				return
			}
			// Mean of each country should not be calculated.
			allCountries = false
		} else { // If year is empty, it tells user it is empty, but still allows for code to run.
			http.Error(w, "Year query should not be empty.", http.StatusLengthRequired)
			return
		}
	}

	// Calculates the mean of grouped countries. Overrides allCountries=true.
	if r.URL.Query().Has("mean") {
		if strings.Contains(strings.ToLower(r.URL.Query().Get("mean")), "true") {
			listOfRSE = meanCalculation(listOfRSE)
			listOfRSE = utility.SortRSEList(listOfRSE, true, constants.ASCENDING)
		} else { // To inform the user that mean query must be true to use.
			http.Error(w, "?mean=true, required to use mean.", http.StatusMethodNotAllowed)
			return
		}
	}

	// If all countries is to be printed, it will calculate the mean first, then sort it alphabetically.
	if allCountries {
		listOfRSE = meanCalculation(listOfRSE)
		// Sorts alphabetically as meanCalculation uses maps, which randomizes entries.
		listOfRSE = utility.SortRSEList(listOfRSE, true, constants.ASCENDING)
	}

	// Handles sort query.
	var sortErr error
	listOfRSE, sortErr = SortQueryHandler(r, listOfRSE)
	if sortErr != nil {
		http.Error(w, sortErr.Error(), http.StatusBadRequest)
		return
	}

	// Checks if list is empty.
	if len(listOfRSE) == 0 {
		http.Error(w, "Nothing matching your search terms.", http.StatusBadRequest)
		return
	}

	// Resets country identifier.
	countryIdentifier = ""
	// Encodes list and prints to console.
	utility.Encoder(w, listOfRSE)
}

// beginEndLimiter Function to allow for searching to and from a year.
func beginEndLimiter(begin string, end string, listToIterate []structs.RenewableShareEnergyElement) ([]structs.RenewableShareEnergyElement, error) {
	var newList []structs.RenewableShareEnergyElement
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
	// Append json objects fitting conditions to newList.
	for _, v := range listToIterate {
		relevantYear := v.Year
		if toFromOr == 3 && relevantYear <= convEnd && convBegin <= relevantYear {
			newList = append(newList, v)
		} else if toFromOr == 1 && convBegin <= relevantYear {
			newList = append(newList, v)
		} else if toFromOr == 2 && relevantYear <= convEnd {
			newList = append(newList, v)
		}
	}
	// Returns mean of years between, as long as begin and end is not the same.
	if toFromOr == 3 && convBegin != convEnd {
		newList = meanCalculation(newList)
		// Sorts newList, as mean calculation randomizes entries.
		newList = utility.SortRSEList(newList, true, constants.ASCENDING)
	}
	return newList, nil
}

// meanCalculation Function to calculate the mean of percentage per country, from the inputted list.
func meanCalculation(listToIterate []structs.RenewableShareEnergyElement) []structs.RenewableShareEnergyElement {
	// If listToIterate is empty, nothing is done.
	if len(listToIterate) == 0 {
		return []structs.RenewableShareEnergyElement{}
	}
	// Creates a map for counting and collecting percentages.
	meanMap := make(map[string]structs.RenewableShareEnergyElement)
	countMap := make(map[string]int)

	// Loops through listToIterate and inserts into newly created maps.
	for _, v := range listToIterate {
		key := v.Name
		// Value returned is not relevant, exits is a bool if it exists in map.
		_, exists := meanMap[key]
		// Adds new entry if it does not exist.
		if !exists {
			meanMap[key] = structs.RenewableShareEnergyElement{
				Name:       v.Name,
				IsoCode:    v.IsoCode,
				Percentage: 0,
			}
		}
		// Cannot modify map values directly, has to extract and then reassign.
		mapValueExtracted := meanMap[key]
		mapValueExtracted.Percentage = mapValueExtracted.Percentage + v.Percentage
		meanMap[key] = mapValueExtracted
		// Increments count to be used to calculate mean.
		countMap[key]++
	}

	// Create a new listToIterate to be appended to.
	resultCalc := make([]structs.RenewableShareEnergyElement, len(meanMap))
	i := 0
	for _, v := range meanMap {
		amount := countMap[v.Name]
		// Removes the possibility for division by 0.
		if amount == 0 {
			continue
		}
		// Calculates the mean.
		v.Percentage /= float64(amount)
		resultCalc[i] = v
		// Increments, to append to next index.
		i++
	}
	// Returns the results, year is not added to result list and therefore omitted.
	return resultCalc
}
