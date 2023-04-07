package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

// TestHandlerStatus_ValidMethod checks if the correct HTTP status code is returned by the handler
// when using a valid method (GET).
// Returns: http.StatusOK, or an error message.
func TestHandlerStatus_ValidMethod(t *testing.T) {
	InitWebhookRegistrations()

	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerStatus)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// TestHandlerStatus_InvalidMethod checks the handlers behaviour when an invalid HTTP method is used to access the endpoint.
// Returns: http.StatusMethodNotAllowed, or an error message.
func TestHandlerStatus_InvalidMethod(t *testing.T) {
	InitWebhookRegistrations()

	req, err := http.NewRequest("POST", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerStatus)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

// TestHandlerStatus_GetStatusError checks if the getStatus() function returns an error when accessing the country API fails.
// Returns: an error.
func TestHandlerStatus_GetStatusError(t *testing.T) {
	InitWebhookRegistrations()

	// Create a mock http client that returns an error when accessing the country API
	mockErrClient := &http.Client{
		Transport: &mockRoundTripper{
			roundTripFunc: func(*http.Request) (*http.Response, error) {
				return nil, errors.New("error accessing country API")
			},
		},
	}

	// Replace the global client with the mock client for this test case
	client = mockErrClient

	_, err := getStatus()
	if err == nil {
		t.Error("getStatus() did not return an error when accessing country API failed")
	}
}

func TestHandlerStatus_GetStatusSuccess(t *testing.T) {
}

func TestHandlerStatus_JSONEncoding(t *testing.T) {
}

func TestHandlerStatus_UnavailableThirdParty(t *testing.T) {
}
