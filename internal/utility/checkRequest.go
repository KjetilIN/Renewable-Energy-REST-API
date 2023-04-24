package utility

import "net/http"

// CheckRequest Checks if request type is as expected.
func CheckRequest(r *http.Request, expected string) bool {
	if r.Method == expected {
		return true
	} else {
		return false
	}
}
