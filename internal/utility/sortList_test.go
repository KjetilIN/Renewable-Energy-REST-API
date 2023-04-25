package utility

import (
	"assignment-2/internal/webserver/structs"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestSortRSEList Tests sorting method in utility.
func TestSortRSEList(t *testing.T) {
	listForTesting := setupList()

	// Sorts list by percentage, in ascending order.
	listForTesting = SortRSEList(listForTesting, false, 1)
	for i := 1; i < len(listForTesting); i++ {
		assert.LessOrEqual(t, listForTesting[i].Percentage, listForTesting[i-1].Percentage, "List is not sorted.")
	}
	// Sorts list by percentage, in descending order.
	listForTesting = SortRSEList(listForTesting, false, 2) // Descending percentage
	for i := 1; i < len(listForTesting); i++ {
		assert.LessOrEqual(t, listForTesting[i-1].Percentage, listForTesting[i].Percentage, "List is not sorted.")
	}
	// Sorts list alphabetically, in ascending order.
	listForTesting = SortRSEList(listForTesting, true, 1)
	for i := 1; i < len(listForTesting); i++ {
		assert.GreaterOrEqual(t, listForTesting[i].Name, listForTesting[i-1].Name, "List is not sorted.")
	}
	// Sorts list alphabetically, in descending order.
	listForTesting = SortRSEList(listForTesting, true, 2)
	for i := 1; i < len(listForTesting); i++ {
		assert.LessOrEqual(t, listForTesting[i].Name, listForTesting[i-1].Name, "List is not sorted.")
	}
}

// setupList Prepares a list used for testing.
func setupList() []structs.RenewableShareEnergyElement {
	elementOne := structs.RenewableShareEnergyElement{
		Name:       "XYZ",
		IsoCode:    "XYZ",
		Year:       2020,
		Percentage: 45,
	}
	var listForTesting []structs.RenewableShareEnergyElement
	listForTesting = append(listForTesting, elementOne)

	elementTwo := structs.RenewableShareEnergyElement{
		Name:       "ABC",
		IsoCode:    "ABC",
		Year:       2020,
		Percentage: 33,
	}
	listForTesting = append(listForTesting, elementTwo)

	elementThree := structs.RenewableShareEnergyElement{
		Name:       "GHI",
		IsoCode:    "GHI",
		Year:       2020,
		Percentage: 90,
	}
	listForTesting = append(listForTesting, elementThree)

	return listForTesting
}
