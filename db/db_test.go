package db

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	firestore "cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Function for the test package to clear the collection
// Used before each test
func clearFirestoreCollection(t *testing.T, client *firestore.Client) {
    iter := client.Collection(constants.FIRESTORE_COLLECTION_TEST).Documents(context.Background())
    for {
		// If there is a new document, get it
        doc, err := iter.Next()
        if err == iterator.Done {
            break
        }
        if err != nil {
            t.Fatalf("Failed to iterate over Firestore documents: %v", err)
        }

		// Delete the document 
        _, err = doc.Ref.Delete(context.Background())
        if err != nil {
            t.Fatalf("Failed to delete Firestore document %s: %v", doc.Ref.ID, err)
        }
    }
}

func TestCheckFirestoreConnection(t *testing.T) {
	// Call the function to check the Firestore connection.
	statusCode := CheckFirestoreConnection()

	// Verify that the status code is 200 (OK).
	if statusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", statusCode)
	}
}

func TestAddWebhook(t *testing.T) {
    // Create a new webhook with some default values
    webhook := structs.WebhookID{
		ID: "TEST_ID" + time.Now().Local().String(),
		Webhook: structs.Webhook{
			Url: "TEST_URL",
			Country: "TEST_COUNTRY",
			Calls: 5,
		},
		Created: time.Now(),
			
    }

    // Get the Firestore client
    client, err := getFirestoreClient()
    if err != nil {
        t.Fatalf("Failed to get Firestore client: %v", err)
    }
    defer client.Close()

	// Clear the collection before the test
    clearFirestoreCollection(t, client)
    

    // Add the webhook to Firestore test collection
    err = AddWebhook(webhook, constants.FIRESTORE_COLLECTION_TEST)
    if err != nil {
        t.Fatalf("Failed to add webhook to Firestore: %v", err)
    }

    // Verify that the webhook was added to Firestore
    doc, findDocError := client.Collection(constants.FIRESTORE_COLLECTION_TEST).Doc(webhook.ID).Get(context.Background())
    if findDocError != nil {
        t.Fatalf("Failed to get webhook from Firestore: %v", err)
    }

	// We assume that most variables are the same but not the timestamp.
	// When testing this is usually the result 
	// Timestamp of the created object: 2023-04-17 19:45:42.709753 +0000 UTC
	// Timestamp of the retrieved document : 2023-04-17 21:45:42.7097538 +0200 CEST m=+0.013129501 
	// For some reason they have different timezones. 

	// Solution: only look at the valuable webhook information
	var storedWebhook structs.WebhookID
    doc.DataTo(&storedWebhook)
    if !reflect.DeepEqual(storedWebhook.Webhook, webhook.Webhook) {
        t.Errorf("Stored webhook %v does not match expected webhook %v", storedWebhook, webhook)
    }

}
