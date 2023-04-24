package db

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	firestore "cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
			Url:     "TEST_URL",
			Country: "TEST_COUNTRY",
			Calls:   5,
			Event: constants.CALLS_EVENT,
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
			ID:      fmt.Sprintf("test-%d", i),
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
		ID:      "test-webhook-id",
		Webhook: structs.Webhook{
			Url: "THIS IS URL",
			Country: "NOR",
			Calls: 0,
			Event: constants.COUNTRY_API_EVENT,
		},
		Created: time.Now(),
	}
	err = AddWebhook(testWebhook, constants.FIRESTORE_COLLECTION_TEST)
	if err != nil {
		t.Error("AddWebHook test failure caused FetchWebhookWithID to fail....")
		t.Fatal("Failed to add webhook for testing: ", err)
	}

	// Fetch the webhook using its ID
	fetchedWebhook, err := FetchWebhookWithID("test-webhook-id", constants.FIRESTORE_COLLECTION_TEST)
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
		ID:      "12345",
		Webhook: structs.Webhook{},
		Created: time.Now(),
	}
	webhook2 := structs.WebhookID{
		ID:      "67890",
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
	if len(expectedWebhooks) != len(webhooksFromFirestore) {
		t.Fatal("Expected the length of all webhooks to be the same as the one list given, but they were not")
	}
}

func TestDeleteWebhook(t *testing.T) {
	// Get the Firestore client
	client, err := getFirestoreClient()
	if err != nil {
		t.Fatalf("Failed to get Firestore client: %v", err)
	}
	defer client.Close()

	// Clear the collection before the test
	clearFirestoreCollection(t, client)

	// Create a list of mock webhooks to add to Firestore
	webhookToDelete := structs.WebhookID{
		ID:      "ID_TO_BE_DELETED",
		Webhook: structs.Webhook{},
		Created: time.Now(),
	}
	webhookToKeep := structs.WebhookID{
		ID:      "ID_TO_KEEP",
		Webhook: structs.Webhook{},
		Created: time.Now(),
	}
	expectedWebhooks := []structs.WebhookID{webhookToDelete, webhookToKeep}

	// Add the webhooks to Firestore
	for _, webhook := range expectedWebhooks {
		_, err = client.Collection(constants.FIRESTORE_COLLECTION_TEST).Doc(webhook.ID).Set(context.Background(), webhook)
		if err != nil {
			t.Fatalf("Error adding webhook to Firestore: %v", err)
		}
	}

	// Delete the webhook from the firestore
	deletionError := DeleteWebhook(webhookToDelete.ID, constants.FIRESTORE_COLLECTION_TEST)
	if deletionError != nil {
		t.Errorf("Unexpected error while deleting webhook: %s", err)
	}

	// Loop until the webhook is no longer found in Firestore or the maximum number of attempts is reached
	// We do this because the webhook could still be in the collection when we do the get request for the deleted document
	// However, it should still be deleted after some moments after. This is why we do get request until we don't find it.
	// Limit is set to make sure that we don't over
	maxAttempts := 10
	attempts := 0
	for {
		// Do a attempt at retrieving the document
		_, err := client.Collection(constants.FIRESTORE_COLLECTION_TEST).Doc(webhookToDelete.ID).Get(context.Background())
		if err != nil {
			// There was no error but we still need to check if the status is okay
			// Use the gRPC package to execute this
			status, ok := status.FromError(err)
			if ok && status.Code() == codes.NotFound {
				// The document has been successfully deleted
				break
			} else {
				// An unexpected error occurred, return an error and fail the test
				t.Errorf("Unexpected error while retrieving webhook: %s", err)
				return
			}
		}

		// Sleep for a short time before checking again
		time.Sleep(time.Second)
		attempts++

		// If we have gone over the limit of attempts, then the document was not deleted correctly
		if attempts >= maxAttempts {
			t.Errorf("Webhook was not deleted after %d attempts", maxAttempts)
			return
		}
	}

	// Check if the webhook that was not deleted is still in the collection
	keptDoc, err := client.Collection(constants.FIRESTORE_COLLECTION_TEST).Doc(webhookToKeep.ID).Get(context.Background())
	if keptDoc == nil {
		t.Error("Webhook that was not supposed to be deleted was deleted")
	}
}

func TestPurgeWebhooks(t *testing.T) {
	// Get the Firestore client
	client, err := getFirestoreClient()
	if err != nil {
		t.Fatalf("Failed to get Firestore client: %v", err)
	}
	defer client.Close()

	// Clear the collection before the test
	clearFirestoreCollection(t, client)

	// Add more than the maximum allowed webhooks to the collection
	maxWebhookCount := 3
	for i := 0; i < maxWebhookCount+1; i++ {
		webhook := structs.WebhookID{
			ID:      "ID_OF_WEBHOOK_NR_" + strconv.Itoa(i),
			Webhook: structs.Webhook{},
			Created: time.Now(),
		}
		err := AddWebhook(webhook, constants.FIRESTORE_COLLECTION_TEST)
		if err != nil {
			t.Errorf("Error adding webhook to test collection: %s", err)
		}
	}

	// Purge the webhooks in the collection
	err = PurgeWebhooks(constants.FIRESTORE_COLLECTION_TEST, maxWebhookCount)
	if err != nil {
		t.Errorf("Error purging webhooks from test collection: %s", err)
	}

	// Do 5 attempt at getting the amount.
	// Same error as deletion test, were you need to redo the get request, in this case the amount of webhooks
	var numberOfWebhooks int
	attempts := 15
	for i := 0; i < attempts; i++ {
		// Get the number of webhooks and check
		numberOfWebhooks = GetNumberOfWebhooks(constants.FIRESTORE_COLLECTION_TEST)
		if numberOfWebhooks == maxWebhookCount {
			break
		}
		// Sleep before next attempt
		time.Sleep(time.Second)
	}
	// After the attempts, check if we have the expected amount of webhooks
	if numberOfWebhooks != maxWebhookCount {
		// Fatal error so we return to make it easier to decode
		t.Errorf("Expected %d webhooks, got %d", maxWebhookCount, numberOfWebhooks)
		return
	}

	// Try again but this time we go under the limit
	// First we clear the database with webhooks:
	clearFirestoreCollection(t, client)

	// Then we add less webhooks than the limit.
	// The limit is set to 3, but we add 2 webhooks
	for i := 0; i < maxWebhookCount-1; i++ {
		webhook := structs.WebhookID{
			ID:      "ID_OF_WEBHOOK_NR_" + strconv.Itoa(i),
			Webhook: structs.Webhook{},
			Created: time.Now(),
		}
		err := AddWebhook(webhook, constants.FIRESTORE_COLLECTION_TEST)
		if err != nil {
			t.Errorf("Error adding webhook to test collection: %s", err)
		}
	}

	// Purge the webhooks in the collection
	err = PurgeWebhooks(constants.FIRESTORE_COLLECTION_TEST, maxWebhookCount)
	if err != nil {
		t.Errorf("Error purging webhooks from test collection: %s", err)
	}

	// Do 10 attempts at getting the amount, same as before
	// This time we check if there is the same amount of webhooks in the collection
	expectedAmountOfWebhooks := maxWebhookCount - 1
	for i := 0; i < attempts; i++ {
		// Get the number of webhooks and check
		numberOfWebhooks = GetNumberOfWebhooks(constants.FIRESTORE_COLLECTION_TEST)
		if numberOfWebhooks == expectedAmountOfWebhooks {
			break
		}
		// Sleep before next attempt
		time.Sleep(time.Second)
	}
	// After the attempts, check if we have the expected amount of webhooks
	if numberOfWebhooks != expectedAmountOfWebhooks {
		t.Errorf("Expected %d webhooks, got %d", expectedAmountOfWebhooks, numberOfWebhooks)
		return
	}

}

func TestInvocate(t *testing.T) {
	// Get the Firestore client
	client, err := getFirestoreClient()
	if err != nil {
		t.Fatalf("Failed to get Firestore client: %v", err)
	}
	defer client.Close()

	// Clear the collection before the test
	clearFirestoreCollection(t, client)

	// Add information to test on. Two of every Country
	countries := []string{"USA", "NOR", "SWE"}
	for _, country := range countries {
		for i := 0; i < 2; i++ {
			webhook := structs.WebhookID{
				ID: "ID_OF_WEBHOOK_NR_" + strconv.Itoa(i) + country,
				Webhook: structs.Webhook{
					Url:     "https://exsample.no",
					Country: country, // add webhook to one of three countries
					Calls:   1,
				},
				Created:     time.Now(),
				Invocations: 0,
			}
			err := AddWebhook(webhook, constants.FIRESTORE_COLLECTION_TEST)
			if err != nil {
				t.Errorf("Error adding webhook to test collection: %s", err)
			}
		}
	}

	//IncrementInvocations USA 2 times, NOR 1 time and SWE 0 times
	IncrementInvocations("USA", constants.FIRESTORE_COLLECTION_TEST)
	IncrementInvocations("USA", constants.FIRESTORE_COLLECTION_TEST)
	IncrementInvocations("NOR", constants.FIRESTORE_COLLECTION_TEST)

	// Fetch all webhooks from Firestore
	webhooksFromFirestore, err := FetchAllWebhooks(constants.FIRESTORE_COLLECTION_TEST)
	if err != nil {
		t.Fatalf("Error fetching webhooks from Firestore in the testing of log invocation: %v", err)
	}

	// Check if each country has the correct amount of invocations
	for _, webhook := range webhooksFromFirestore {
		switch webhook.Country {
		case "NOR":
			if webhook.Invocations != 1 {
				t.Errorf("Error in webhooks given: NOR was expected to have %v invocations but had %v", 1, webhook.Invocations)
				return
			}
			break
		case "USA":
			if webhook.Invocations != 2 {
				t.Errorf("Error in webhooks given: USA was expected to have %v invocations but had %v", 2, webhook.Invocations)
				return
			}
			break
		case "SWE":
			if webhook.Invocations != 0 {
				t.Errorf("Error in webhooks given: NOR was expected to have %v invocations but had %v", 0, webhook.Invocations)
				return
			}
			break

		default:
			t.Errorf("Found a webhook that was not supposed to be here: %s", webhook.Country)
			return

		}
	}

}

func TestCallUrl(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a 200 status code and a message
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "This is a test response")
	}))
	defer server.Close()

	// Create a webhook with the mock server's URL
	webhook := structs.WebhookID{
		ID: "test-webhook-id",
		Webhook: structs.Webhook{
			Url: server.URL,
		},
		Created: time.Now(),
	}

	// Call the URL
	err := CallUrl(webhook)
	if err != nil {
		t.Error("Error calling URL: " + err.Error())
	}
}
