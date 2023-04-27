package utility

import (
	"assignment-2/internal/constants"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetParams Tests if generalized parameter retrieval works.
func TestGetParams(t *testing.T) {
	url := "/energy/v1/renewables/history/united%20states%20of%20america/tesy"
	param := GetParams(url, constants.HISTORY_PATH)
	assert.Equal(t, "united states of america", param, "Incorrect param returned.")
	url = "/energy/v1/renewables/current/norway"
	param = GetParams(url, constants.CURRENT_PATH)
	assert.Equal(t, "norway", param, "Incorrect param returned.")
}

func TestGetOneFirstComponentOnly_Correct_Use(t *testing.T) {

	// Correct given prefix and
	prefixPath := "http://localhost:8000/energy/v1/notifications/"
    givenPath := "http://localhost:8000/energy/v1/notifications/temp"

    // Test a valid path
    expected := "temp"
    result, err := GetOneFirstComponentOnly(prefixPath, givenPath)
    if err != nil {
        t.Error("Unexpected error: " + err.Error())
		return 
    }
    if result != expected {
        t.Error("Expected '"+ expected + "' but got '"+ result + "'")
		return 
    }

}

func TestGetOneFirstComponentOnly_Incorrect_Prefix(t *testing.T) {
	// Incorrect given prefix and
	prefixPath := "http://localhost:8000/v1/notifications/"
    givenPath := "http://localhost:8000/energy/v1/notifications/temp"

    // Test a valid path
    _, err := GetOneFirstComponentOnly(prefixPath, givenPath)
    if err == nil {
        t.Error("Expected error!")
		return 
    }
}


func TestGetOneFirstComponentOnly_Multiple_Components(t *testing.T) {
	// Correct given prefix, but with multiple components
	prefixPath := "http://localhost:8000/v1/notifications/"
    givenPath := "http://localhost:8000/energy/v1/notifications/temp/okay"

    // Test a valid path
    _, err := GetOneFirstComponentOnly(prefixPath, givenPath)
    if err == nil {
        t.Error("Expected error!")
		return 
    }
}


func TestGetOneFirstComponentOnly_With_Empty_Component(t *testing.T) {
	// Correct given prefix, but with an empty component
	prefixPath := "http://localhost:8000/v1/notifications/"
    givenPath := "http://localhost:8000/energy/v1/notifications/       "

    // Test a valid path
    component , err := GetOneFirstComponentOnly(prefixPath, givenPath)
    if err == nil {
        t.Error("Expected error!")
		return 
    }

	// Expect the component to be empty string 
	if component != ""{
		t.Error("Expected component to be empty, but instead was: '" + component + "'")
		return 
	}
}