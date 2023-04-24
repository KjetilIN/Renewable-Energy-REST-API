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
