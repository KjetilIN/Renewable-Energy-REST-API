package handlers

import (
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/structs"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// ASCENDING Used in sorting method to sort ascending.
const ASCENDING = 1

// DESCENDING Used in sorting method to sort descending.
const DESCENDING = 2

// HandlerHistory is a handler for the /history endpoint.
func HandlerHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	// Reads from csv and returns json list.
	listOfRSE, jsonError := utility.RSEToJSON()
	if jsonError != nil {
		http.Error(w, jsonError.Error(), http.StatusInternalServerError)
		return
	}

	// Collects parameters, separated by /
	params := strings.Split(r.URL.Path, "/") //Used to split the / in path to collect search parameters.

	// Checks if an optional parameter is passed.
	if len(params) == 6 {
		countryIdentifier := params[5]
		var filteredList []structs.RenewableShareEnergyElement
		filteredList = countryCodeLimiter(listOfRSE, countryIdentifier)

		// If list is empty, could not find country by country code.
		if len(filteredList) == 0 {
			// Checks if country searched for is a full common name.
			country, countryConversionErr := utility.GetCountry(countryIdentifier, false)
			if countryConversionErr != nil {
				http.Error(w, "Did not find country based on search parameters: "+countryConversionErr.Error(), http.StatusNoContent)
				return
			}
			// Checks if country code is empty.
			if country.CountryCode != "" {
				filteredList = countryCodeLimiter(listOfRSE, country.CountryCode)
			}
		}
		listOfRSE = filteredList
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

	// If Query: mean=true, a different struct type will be encoded to client. It calculates the mean of grouped countries.
	if r.URL.Query().Has("mean") && strings.Contains(strings.ToLower(r.URL.Query().Get("mean")), "true") {
		listOfRSE = meanCalculation(listOfRSE)
	}

	// Checks if sortByValue query is passed.
	if r.URL.Query().Has("sortbyvalue") && strings.Contains(strings.ToLower(r.URL.Query().Get("sortbyvalue")), "true") {
		// Sorts percentage descending if descending query is true.
		if strings.Contains(strings.ToLower(r.URL.Query().Get("descending")), "true") {
			listOfRSE = sliceSortingByValue(listOfRSE, false, DESCENDING)
		} else { // Sorting standard is ascending if nothing else is passed.
			listOfRSE = sliceSortingByValue(listOfRSE, false, ASCENDING)
		}
	}

	// Checks if sortAlphabetically query is passed.
	if r.URL.Query().Has("sortalphabetically") && strings.Contains(strings.ToLower(r.URL.Query().Get("sortalphabetically")), "true") {
		// Sorts list descending if descending query is true.
		if strings.Contains(strings.ToLower(r.URL.Query().Get("descending")), "true") {
			listOfRSE = sliceSortingByValue(listOfRSE, true, DESCENDING)
		} else { // Sorting standard is ascending if nothing else is passed.
			listOfRSE = sliceSortingByValue(listOfRSE, true, ASCENDING)
		}
	}

	// Checks if list is empty.
	if len(listOfRSE) == 0 {
		http.Error(w, "Nothing matching your search terms.", http.StatusNoContent)
		return
	}

	// Encodes list and prints to console.
	utility.Encoder(w, listOfRSE)
}

// countryCodeLimiter Method to limit a list based on country code.
func countryCodeLimiter(listToIterate []structs.RenewableShareEnergyElement, countryCode string) []structs.RenewableShareEnergyElement {
	var limitedList []structs.RenewableShareEnergyElement
	for i, v := range listToIterate { // Iterates through input list.
		if strings.Contains(strings.ToLower(listToIterate[i].IsoCode), strings.ToLower(countryCode)) { // If country code match it is
			// appended to new list.
			limitedList = append(limitedList, v)
		}
	}
	return limitedList // Returns list containing all matching countries.
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
	// Returns mean of years between.
	if toFromOr == 3 {
		newList = meanCalculation(newList)
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

// sliceSortingByValue A function which sorts a json list based on value, using inbuilt sort method.
func sliceSortingByValue(listToIterate []structs.RenewableShareEnergyElement, alphabetical bool, sortingMethod int) []structs.RenewableShareEnergyElement {
	// Sorts list, based on alphabetical boolean and sortingMethods value.
	switch {
	// Sorts by percentage, ascending.
	case sortingMethod == ASCENDING && !alphabetical:
		sort.Slice(listToIterate, func(i, j int) bool {
			return listToIterate[j].Percentage < listToIterate[i].Percentage
		})
	// Sorts by percentage, descending.
	case sortingMethod == DESCENDING && !alphabetical:
		sort.Slice(listToIterate, func(i, j int) bool {
			return listToIterate[i].Percentage < listToIterate[j].Percentage
		})
	// Sorts alphabetically, ascending.
	case sortingMethod == ASCENDING && alphabetical:
		sort.Slice(listToIterate, func(i, j int) bool {
			return listToIterate[i].Name < listToIterate[j].Name
		})
	// Sorts alphabetically, descending.
	case sortingMethod == DESCENDING && alphabetical:
		sort.Slice(listToIterate, func(i, j int) bool {
			return listToIterate[j].Name < listToIterate[i].Name
		})
	}
	return listToIterate
}
