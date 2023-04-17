package db

import (
	"net/http"
	"strconv"
	"testing"
)

func TestCheckFirestoreConnection(t *testing.T) {
    // Call the function to check the Firestore connection.
    statusCode := CheckFirestoreConnection()

    // Verify that the status code is 200 (OK).
	if statusCode != http.StatusOK{
		t.Fatal("Expected status code 200, got " + strconv.Itoa(statusCode))
	}
    
}
