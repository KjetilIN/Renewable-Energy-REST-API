package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleGetRequest(t *testing.T) {

	// Create the test server
	server := httptest.NewServer(http.HandlerFunc(HandlerNotifications))

	// Do the get request
	resp, err := http.Get(server.URL + constants.NOTIFICATIONS_PATH)
	if err != nil {
		t.Fatal("Error on get request to notification endpoint")
		return
	}
	defer resp.Body.Close()

	// Status code should be 200
	if resp.StatusCode != 200 {
		t.Fatal("Error on doing a request for an endpoint")
		return
	}

	defer server.Close()

}

func TestHandlePostGetAndDelete(t *testing.T) {

	// Create the test server
	server := httptest.NewServer(http.HandlerFunc(HandlerNotifications))
	defer server.Close()

	// Webhook to be deleted
	webhook := structs.WebhookID{
		ID: "WebhookID",
		Webhook: structs.Webhook{
			Url:     "http://fakeUrl.com",
			Country: "NOR",
			Calls:   3,
			Event:   "CALLS",
		},
		Created:     time.Now(),
		Invocations: 2,
	}

	// Encode the webhook
	webhookJSON, err := json.Marshal(webhook)
	if err != nil {
		t.Fatal("Error on encoding test webhook to JSON:" + err.Error())
		return
	}

	// Create a new HTTP request with the POST method and the request body
	req, err := http.NewRequest(http.MethodPost, server.URL+constants.NOTIFICATIONS_PATH, bytes.NewReader(webhookJSON))
	if err != nil {
		t.Fatal("Error on the post method:" + err.Error())
		return
	}

	// Send the request and get the response
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
		return
	}
	defer resp.Body.Close()

	// Check that the response status code is what we expect
	if status := resp.StatusCode; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check that we can decode the struct into a struct
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	var registration structs.IdResponse
	unmarshalError := json.Unmarshal(bodyBytes, &registration)
	if unmarshalError != nil {
		log.Fatal("Error on unmarshal the results from creating a webhook: " + unmarshalError.Error())
		return
	}

	//Sleep before doing another request
	time.Sleep(time.Second * 5)

	// Check if we can get the webhook
	getRequest, err := http.Get(server.URL + constants.NOTIFICATIONS_PATH + registration.ID)
	if err != nil {
		t.Fatal("Error on doing the get request")
		return
	}

	// UnMarshaling the response
	getBodyBytes, err := ioutil.ReadAll(getRequest.Body)

	var webhookRetrieved structs.WebhookID
	unmarshalError = json.Unmarshal(getBodyBytes, &webhookRetrieved)
	if unmarshalError != nil {
		log.Fatal("Error on unmarshal the results from creating a webhook: " + unmarshalError.Error())
		return
	}

	// Check the webhook matched some details
	if webhookRetrieved.ID != registration.ID {
		log.Fatal("Webhook retrieved did not match the initial created webhook")
		return
	}

	// Delete the webhook

	// Create a new HTTP request with the DELETE method and no parameters
	delRequest, delError := http.NewRequest(http.MethodDelete, server.URL+constants.NOTIFICATIONS_PATH+registration.ID, nil)
	if delError != nil {
		t.Fatal("Error on creating the delete request;" + delError.Error())
		return
	}

	// Send the request using the default HTTP client
	resp, delError = http.DefaultClient.Do(delRequest)
	if err != nil {
		t.Fatal("Error on calling the delete request:" + delError.Error())
		return
	}
	defer resp.Body.Close()

	// Check the response status code to make sure the request was successful
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Error: Response Code should be 200 but was %v", resp.StatusCode)
		return
	}

	// Wait before checking
	time.Sleep(time.Second * 5)

	// This time we should not get the deleted webhook
	newGetRequest, err := http.Get(server.URL + constants.NOTIFICATIONS_PATH + registration.ID)
	if err != nil {
		t.Fatal("Error on doing the get request")
		return
	}

	// Should get no content back
	if newGetRequest.StatusCode != http.StatusNotFound {
		t.Fatalf("Error: Status code should be 404 but was %s", newGetRequest.Status)
		return
	}
}

func TestCreateWebhookID(t *testing.T) {
	webhook := structs.Webhook{
		Url:     "https://example.com",
		Country: "USA",
	}

	// Call the function to create the webhook ID
	webhookID := createWebhookID(webhook)

	// Verify that the webhook ID is not empty
	// This is the only test that we can do, because we don't know exactly when it was created
	if webhookID == "" {
		t.Errorf("Webhook ID should not be empty")
	}
}
