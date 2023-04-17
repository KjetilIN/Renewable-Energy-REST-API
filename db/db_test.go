package db

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"context"
	"fmt"
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


func TestGetNumberOfWebhooks(t *testing.T) {
    // Get the Firestore client
    client, err := getFirestoreClient()
    if err != nil {
        t.Fatalf("Failed to get Firestore client: %v", err)
    }
    defer client.Close()

	// Clear the collection before the test
    clearFirestoreCollection(t, client)

    // Add some 5 test documents to the collection
    for i := 0; i < 5; i++ {
        webhook := structs.WebhookID{
            ID:   fmt.Sprintf("test-%d", i),
            Webhook: structs.Webhook{},
			Created: time.Now(),
        }
        err = AddWebhook(webhook, constants.FIRESTORE_COLLECTION_TEST)
        if err != nil {
			t.Error("AddWebHook test failure caused GetNumberOfWebhooks to fail....")
            t.Fatalf("Error adding test webhook: %v", err)
        }
    }

    // Get the number of webhooks in the collection
    numWebhooks := GetNumberOfWebhooks(constants.FIRESTORE_COLLECTION_TEST)

    // Verify that the number of webhooks is correct
    if numWebhooks != 5 {
        t.Fatalf("Expected 5 webhooks, but got %d", numWebhooks)
    }
}


func TestFetchWebhookWithID(t *testing.T) {

	// Get the Firestore client
    client, err := getFirestoreClient()
    if err != nil {
        t.Fatalf("Failed to get Firestore client: %v", err)
    }
    defer client.Close()

	// Clear the collection before the test
    clearFirestoreCollection(t, client)


	// Add a webhook to Firestore
	testWebhook := structs.WebhookID{
		ID: "test-webhook-id",
		Webhook: structs.Webhook{},
		Created: time.Now(),
	}
	err = AddWebhook(testWebhook, constants.FIRESTORE_COLLECTION_TEST)
	if err != nil {
		t.Error("AddWebHook test failure caused FetchWebhookWithID to fail....")
		t.Fatal("Failed to add webhook for testing: ", err)
	}

	// Fetch the webhook using its ID
	fetchedWebhook, err := FetchWebhookWithID("test-webhook-id",constants.FIRESTORE_COLLECTION_TEST)
	if err != nil {
		t.Fatal("Failed to fetch webhook: ", err)
	}

	// Check that the fetched webhook matches the added webhook
	if fetchedWebhook.ID != testWebhook.ID || fetchedWebhook.Webhook != testWebhook.Webhook {
		t.Fatalf("Fetched webhook does not match added webhook.\nAdded: %v\nFetched: %v", testWebhook, fetchedWebhook)
	}

	// Try to fetch a non-existent webhook
	_, err = FetchWebhookWithID("non-existent-webhook-id", constants.FIRESTORE_COLLECTION_TEST)
	if err == nil {
		t.Fatal("Expected an error when fetching non-existent webhook, but no error was returned.")
	}
}

func TestFetchAllWebhooks(t *testing.T) {

	// Get the Firestore client
    client, err := getFirestoreClient()
    if err != nil {
        t.Fatalf("Failed to get Firestore client: %v", err)
    }
    defer client.Close()

	// Clear the collection before the test
    clearFirestoreCollection(t, client)


    // Create a list of mock webhooks to add to Firestore
	 webhook1 := structs.WebhookID{
        ID: "12345",
        Webhook: structs.Webhook{},
        Created: time.Now(),
    }
    webhook2 := structs.WebhookID{
        ID: "67890",
        Webhook: structs.Webhook{},
        Created: time.Now(),
    }
    expectedWebhooks := []structs.WebhookID{webhook1, webhook2}

    // Add the webhooks to Firestore
    for _, webhook := range expectedWebhooks {
        _, err = client.Collection(constants.FIRESTORE_COLLECTION_TEST).Doc(webhook.ID).Set(context.Background(), webhook)
        if err != nil {
            t.Fatalf("Error adding webhook to Firestore: %v", err)
        }
    }

    // Fetch all webhooks from Firestore
    webhooksFromFirestore, err := FetchAllWebhooks(constants.FIRESTORE_COLLECTION_TEST)
    if err != nil {
        t.Fatalf("Error fetching webhooks from Firestore: %v", err)
    }

    // Check that at least one of the expected webhooks is in the database
    found := false
	expectedWebhook := expectedWebhooks[0]
	for _, actualWebhook := range webhooksFromFirestore {
		if expectedWebhook.ID == actualWebhook.ID {
			found = true
			break
		}
	}
    
	// Check if the a single webhook was found
    if !found {
        t.Fatalf("Expected webhook not found in database")
    }
	
	// Check if the length of the webhooks list match
	if len(expectedWebhooks) != len(webhooksFromFirestore){
		t.Fatal("Expected the length of all webhooks to be the same as the one list given, but they were not")
	}
	

}
