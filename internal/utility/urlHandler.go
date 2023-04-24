package utility

import "strings"

// ReplaceSpaces Function to remove formatted space: %20, so it is a normal space.
func ReplaceSpaces(url string) string {
	return strings.ReplaceAll(url, "%20", " ")
}

// GetParams Returns a string slice of parameters.
func GetParams(url string, endpoint string) string {
	// Checks if url or endpoint is empty. Returns an empty string if so.
	if url == "" || endpoint == "" {
		return ""
	}

	basisParams := strings.Split(endpoint, "/")
	params := strings.Split(url, "/") //Used to split the / in path to collect search parameters.
	var param string

	for i, v := range basisParams {
		if strings.Contains(strings.ToLower(v), strings.ToLower(params[i])) {
			// If basisParams correspond, it will continue.
			continue
		} else { // If a parameter does not match the basis parameter, it will set it to param string and return it.
			param = params[i]
			param = ReplaceSpaces(param)
		}
	}
	return param
}
