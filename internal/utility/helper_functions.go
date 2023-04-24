package utility

// This file contains helper_functions used to increase the functionality of the service.

// Pluralize Helper function to pluralize words as appropriate.
// Returns: an empty string, or "s".
func Pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
