package utility

import (
	"errors"
	"strconv"
	"strings"
)

// ReplaceSpaces Function to remove formatted space: %20, so it is a normal space.
func ReplaceSpaces(url string) string {
	return strings.ReplaceAll(url, "%20", " ")
}

// GetParams Returns a string slice of parameters.
func GetParams(url string, endpoint string) string {
	// Checks if url or endpoint is empty. Returns an empty string if so.
	if url == "/" || url == "" || endpoint == "" {
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

// Function for getting the first component only from an given path.
// Takes the prefix of the path and the full raw path
// Return the component and any error if there is are more then one component or incorrect prefix path given
func GetOneFirstComponentOnly(prefixPath string, givenPath string) (string, error) {
	// Check if the given prefix exists 
    if !strings.HasPrefix(givenPath, prefixPath) {
        return "", errors.New("Prefix given did not match the url")
    }

	// Remove the prefix
    relativePath := strings.TrimPrefix(givenPath, prefixPath)

	//Split on the relative path 
	componentsList := strings.Split(relativePath , "/");

	//Check if the length is one as expected:
	if len(componentsList) != 1{
		return "", errors.New("Given list of components was expected to be 1 was " + strconv.Itoa(len(componentsList)))
	}

	// Retrieve the first component
	component := strings.ReplaceAll(componentsList[0], " ", "");
    return component, nil
}