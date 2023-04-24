package db

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Function for getting the Firestore client
// Private for security reasons
func getFirestoreClient() (*firestore.Client, error) {
	// Load cred
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return nil, err
	}

	// Use a service account
	ctx := context.Background()
	credentialsPath := os.Getenv("CREDENTIALS_PATH")
	sa := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Println("Credentials not found: " + credentialsPath)
		log.Println("Error on getting the application")
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Println("Credentials file: '" + credentialsPath + "'")
		return nil, err
	}
	return client, nil
}

// Function for getting status code of the connection to the firestore
func CheckFirestoreConnection() int {
	// Connect to to the firestore with the client
	client, err := getFirestoreClient()
	defer client.Close()

	//check for errors on connection
	if err != nil {
		log.Fatal("Error on creating the connection: " + err.Error())
		return http.StatusInternalServerError
	}

	// Test the connection by querying a collection
	docRef := client.Collection(constants.FIRESTORE_COLLECTION)
	if docRef == nil {
		// If there was an error querying the webhook collection, return a 500 status code
		log.Fatal("No collection for the webhooks was found")
		return http.StatusInternalServerError
	}

	// If everything worked, return a 200 status code
	return http.StatusOK
}

// Function that adds a webhook to the firestore, using the given webhook struct and a generated ID.
// Takes the webhook struct and the collection name as parameters.
// Return an error if it could not add the webhook, returns nil if everything went okay .
func AddWebhook(webhook structs.WebhookID, collection string) error {
	// Get the client for the firestore
	client, clientErr := getFirestoreClient()
	defer client.Close()
	if clientErr != nil {
		return clientErr
	}

	// Set default event
	if webhook.Event == ""{
		webhook.Event = constants.CALLS_EVENT
	}
	
	// Turn it into upper case 
	webhook.Event = strings.ToUpper(webhook.Event)
	

	// Create a new doc in the
	_, err := client.Collection(collection).Doc(webhook.ID).Set(context.Background(), webhook)
	if err != nil {
		return err
	}

	return nil
}

// Get number of webhooks.
// Note that if the service is down there will be not handled this function, and 0 wil be returned
// The user has to see the status endpoint
// Takes the name of the collection as a parameter
func GetNumberOfWebhooks(collection string) int {
	//Create a client for the
	client, err := getFirestoreClient()
	defer client.Close()
	if err != nil {
		return 0
	}

	// Get the number of webhooks in the collection.
	allWebHooks, err := client.Collection(collection).Documents(context.Background()).GetAll()
	if err != nil {
		// There was an error but we return 0
		log.Println("Error on getting all webhooks in the GetNumberOfWebhooks method: " + err.Error())
		return 0
	}
	// Return the length of the retrieved given webhooks
	return len(allWebHooks)
}

// Fetch a webhook using its ID.
// The webhook id has to be the same as the document id.
// Takes the id and collection name as parameters.
// Return an error if something went wrong.
func FetchWebhookWithID(id string, collection string) (structs.WebhookID, error) {
	//Create a client for the
	client, err := getFirestoreClient()
	defer client.Close()
	if err != nil {
		return structs.WebhookID{}, err
	}

	var webhook structs.WebhookID
	iter := client.Collection(collection).Documents(context.Background())

	//Loop through each document
	for {
		//Get the document and check if it is done
		doc, err := iter.Next()
		if err == iterator.Done {
			// Break if no more docs to get
			break
		}

		//if the ID is the same as the docs
		if doc.Ref.ID == id {
			log.Println("Webhook found: " + id)
			err := doc.DataTo(&webhook)
			if err != nil {
				log.Println("Webhook with id: " + id + " was found but not decodable")
				return structs.WebhookID{}, err
			}
			//No error on decoding and webhook that matched the id was returned
			return webhook, nil

		}
	}

	// Correctly went through the method but did not find a webhook
	return structs.WebhookID{}, errors.New("No webhook was found in that matched the id: " + id)
}

// Fetch all webhooks
// Returns a list of webhooks with id from the database.
// Takes the name of the collection as parameter.
// Returns an error if something went wrong
func FetchAllWebhooks(collection string) ([]structs.WebhookID, error) {
	//Create a client
	client, err := getFirestoreClient()
	defer client.Close()
	if err != nil {
		return nil, err
	}

	//Iterate through all docs and decode them into the list of structs
	var webhooks []structs.WebhookID
	iter := client.Collection(collection).Documents(context.Background())
	for {
		//Get the document and check if it is done
		doc, err := iter.Next()
		if err == iterator.Done {
			// Break if now more docs to get
			break
		}

		//Check for errors on iterator
		if err != nil {
			//Log the error if any
			log.Println("Failed to iterate: " + err.Error())
		}

		// Decode the webhook into a struct if possible
		var webhook structs.WebhookID
		if err := doc.DataTo(&webhook); err != nil {
			log.Println("Error during data decoding")
		} else {
			// No error so we append the webhook
			webhooks = append(webhooks, webhook)
		}

	}

	// Returns either an empty list or a list of webhooks
	return webhooks, nil

}

// Delete a webhook from a given webhook id
// Takes the given webhook as id and the name of the collection
// No error returns indicates that the process was okay
func DeleteWebhook(webhookID string, collection string) error {
	// Get the client
	client, clientError := getFirestoreClient()
	if clientError != nil {
		return clientError
	}

	// Delete the document based on the id given
	_, err := client.Collection(collection).Doc(webhookID).Delete(context.Background())
	if err != nil {
		return err
	}
	// No error and we return nil
	return nil
}

// Function that removes webhook when we have stored over a set limit
// Takes a collection name as argument.
// Optionally takes the max amount of webhooks, else uses the predefined limit
func PurgeWebhooks(collection string, maxWebhookCount ...int) error {
	// Get the client
	client, clientError := getFirestoreClient()
	if clientError != nil {
		return clientError
	}

	// Get the amount of webhooks
	numberOfWebhooks := GetNumberOfWebhooks(collection)

	// Determine the maximum webhook count
	var webhookLimit int
	if len(maxWebhookCount) > 0 {
		webhookLimit = maxWebhookCount[0]
	} else {
		webhookLimit = constants.MAX_WEBHOOK_COUNT
	}

	// Check if we need to purge
	if numberOfWebhooks <= webhookLimit {
		return nil
	}

	// Notify all subscribers for the event of purging 
	go NotifyForEvent(constants.PURGE_EVENT, constants.FIRESTORE_COLLECTION)

	// Calculate how many of the webhooks we can delete
	numberOfWebhooksToDelete := numberOfWebhooks - webhookLimit

	// Get all the webhooks
	querySnapshot, err := client.Collection(collection).Documents(context.Background()).GetAll()
	if err != nil {
		return err
	}

	// Add the webhooks to a list and sort them by the timestamp
	var webhooks []*structs.WebhookID
	for _, doc := range querySnapshot {
		var webhook structs.WebhookID
		if err := doc.DataTo(&webhook); err != nil {
			return err
		}
		webhooks = append(webhooks, &webhook)
	}

	// Sort based on the oldest first by using indexes to compare the created time
	sort.Slice(webhooks, func(i, j int) bool {
		return webhooks[i].Created.Before(webhooks[j].Created)
	})

	// Use the sorted list of webhooks to delete webhooks from the firestore
	for i := 0; i < numberOfWebhooksToDelete; i++ {
		_, err := client.Collection(collection).Doc(webhooks[i].ID).Delete(context.Background())
		if err != nil {
			log.Println("Error on deleting in purging mechanism: " + err.Error())
			return err
		}
	}

	return nil
}

// Function to call when a alpha code of a country has been used.
// Returns an error if something went wrong
func IncrementInvocations(alphaCode string, collection string) error {
	// Get the client
	client, clientError := getFirestoreClient()
	if clientError != nil {
		return clientError
	}

	// Create a query that filters documents based on the specified alpha code
	query := client.Collection(collection).Where("Country", "==", alphaCode).Documents(context.Background())

	//Iterate thought each document
	for {
		// Get the next document or quit
		doc, err := query.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println("Error on iterating a doc")
			return err
		}

		// Increment the invocations field of the retrieved document by one
		_, err = doc.Ref.Update(context.Background(), []firestore.Update{
			{
				Path:  "Invocations",
				Value: firestore.Increment(1)},
		})
		if err != nil {
			return errors.New("Error on trying to increment the invocations number")
		}

		// Check if we need to call a webhook...
		var webhook structs.WebhookID
		err = doc.DataTo(&webhook)
		if err != nil {
			log.Println("Error on deconstruct the webhook")
			return err
		}
		webhook.Invocations++ // Increment the local version

		// Invocation that has been updated is multiple of two
		if webhook.Calls != 0 && webhook.Invocations%webhook.Calls == 0 && strings.ToUpper(webhook.Event) == "CALLS" {
			// Using the call url method as a go routine
			go CallUrl(webhook)

		}
	}
	return nil
}

// Function using for calling a url
// Takes a webhook and uses its information when calling
// Return an error if something went wrong
func CallUrl(webhook structs.WebhookID) error {
	// Log the attempt for calling an url
	log.Println("Calling the url: " + webhook.Url + "...")

	// Message to return to the user
	var message string
	switch (webhook.Event){
		case constants.COUNTRY_API_EVENT:
			message = "Notification triggered: Country API is down!"
			break
		case constants.PURGE_EVENT:
			message = "Notification triggered: Webhooks has been purged!"
		default:
			message = "Notification triggered: " + strconv.Itoa(webhook.Calls) + " invocations made to " + strings.ToUpper(webhook.Country) + " endpoint."
			break

	}
	
	//Creating the response struct:
	responseWebhook := structs.WebhookCallResponse{
		ID: webhook.ID,
		Webhook: webhook.Webhook,
		Invocations: webhook.Invocations,
		Message: message,
	}

	// Add the struct to the response struct 
	requestBody, err := json.Marshal(responseWebhook)
	if err != nil {
		log.Println("Error encoding responseWebhook to JSON")
		return err
	}

	// Create a new request with the JSON-encoded webhook in the body
	request, creatingError := http.NewRequest(http.MethodPost, webhook.Url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error creating request")
		return creatingError
	}

	// Set the ID as a signature in the request and the application type
	request.Header.Set("Signature-x", webhook.ID)
	request.Header.Set("Content-Type", "application/json")


	// Creating a client and executing the request
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error while trying to execute the request")
		return err
	}

	// Logging that a webhook has been invocated
	log.Println("Webhook with ID: " + webhook.ID + " was invoked. Status code for response is " + strconv.Itoa(response.StatusCode))
	return nil
}


// Function to call when a given event 
func NotifyForEvent(event string, collection string) error {
	// Get the client
	client, clientError := getFirestoreClient()
	if clientError != nil {
		return clientError
	}

	// Create a query that filters documents based on the specified alpha code
	query := client.Collection(collection).Where("Event", "==", strings.ToUpper(event) ).Documents(context.Background())

	//Iterate thought each document
	for {
		// Get the next document or quit
		doc, err := query.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println("Error on iterating a doc")
			return err
		}

		// Check if we need to call a webhook...
		var webhook structs.WebhookID
		err = doc.DataTo(&webhook)
		if err != nil {
			log.Println("Error on deconstruct the webhook")
			return err
		}
		
		// Using the call url method as a go routine
		go CallUrl(webhook)

	
	}
	return nil
}