package utility

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// Tests check request function.
func TestCheckRequest(t *testing.T) {
	var req http.Request
	// Sets request method to GET.
	req.Method = http.MethodGet
	// r needs to be a pointer, so it points to req.
	r := &req
	// Checks if the method GET is passed.
	assert.True(t, CheckRequest(r, http.MethodGet))
	// Sets request to POST.
	req.Method = http.MethodPost
	// Checks if method is GET, expects failure.
	assert.False(t, CheckRequest(r, http.MethodGet))
}
