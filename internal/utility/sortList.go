package utility

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"sort"
)

// SortRSEList A function which sorts a json list based on value, using inbuilt sort method.
func SortRSEList(listToIterate []structs.RenewableShareEnergyElement, alphabetical bool, sortingMethod int) []structs.RenewableShareEnergyElement {
	// Sorts list, based on alphabetical boolean and sortingMethods value.
	switch {
	// Sorts by percentage, ascending.
	case sortingMethod == constants.ASCENDING && !alphabetical:
		sort.Slice(listToIterate, func(i, j int) bool {
			return listToIterate[j].Percentage > listToIterate[i].Percentage
		})
	// Sorts by percentage, descending.
	case sortingMethod == constants.DESCENDING && !alphabetical:
		sort.Slice(listToIterate, func(i, j int) bool {
			return listToIterate[i].Percentage > listToIterate[j].Percentage
		})
	// Sorts alphabetically, ascending.
	case sortingMethod == constants.ASCENDING && alphabetical:
		sort.Slice(listToIterate, func(i, j int) bool {
			return listToIterate[i].Name < listToIterate[j].Name
		})
	// Sorts alphabetically, descending.
	case sortingMethod == constants.DESCENDING && alphabetical:
		sort.Slice(listToIterate, func(i, j int) bool {
			return listToIterate[j].Name < listToIterate[i].Name
		})
	}
	return listToIterate
}
