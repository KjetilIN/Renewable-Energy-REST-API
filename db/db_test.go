package db

import (
	"net/http"
	"testing"
)

func TestCheckFirestoreConnection(t *testing.T) {
	// Call the function to check the Firestore connection.
	statusCode := CheckFirestoreConnection()

	// Verify that the status code is 200 (OK).
	if statusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", statusCode)
	}
}
