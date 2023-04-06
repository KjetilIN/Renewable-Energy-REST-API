package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestHandlerStatus_InvalidMethod(t *testing.T) {
}

func TestHandlerStatus_GetStatusError(t *testing.T) {
}

func TestHandlerStatus_GetStatusSuccess(t *testing.T) {
}

func TestHandlerStatus_JSONEncoding(t *testing.T) {
}

func TestHandlerStatus_UnavailableThirdParty(t *testing.T) {
}
