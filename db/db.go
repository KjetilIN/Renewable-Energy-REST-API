package db

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"context"
	"errors"
	"log"
	"net/http"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Function for getting the Firestore client
func GetFirestoreClient() (*firestore.Client, error) {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("cloud-assigment-2-36e8e-5557620affae.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

//Function for getting status code of the connection to the firestore 
func CheckFirestoreConnection() int {
	// Connect to to the firestore with the client
	client, err := GetFirestoreClient();
	defer client.Close()

	//check for errors on connection 
	if err != nil{
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


//Function that adds a webhook to the firestore, using the given webhook struct and a generated ID.
//Return an error if it could not add the webhook, returns nil if everything went okay 
func AddWebhook(webhook structs.WebhookID) error{
	// Get the client for the firestore
	client, clientErr := GetFirestoreClient()
	defer client.Close()
	if clientErr != nil{
		return clientErr
	}

	// Create a new doc in the 
	_ , err := client.Collection(constants.FIRESTORE_COLLECTION).Doc(webhook.ID).Set(context.Background(),webhook)
	if err != nil{
		return err
	}

	return nil
}


//Get number of webhooks. 
// Note that if the service is down there will be not handled this function. 
// The user has to see the status endpoint
func GetNumberOfWebhooks() int{
	//Create a client for the 
	client, err := GetFirestoreClient()
	defer client.Close()
	if err != nil{
		return 0;
	}

	// Init a iterator
	docsIter := client.Collection(constants.FIRESTORE_COLLECTION).Documents(context.Background())

	// Initialize a counter variable
	count := 0

	// Iterate over the documents and increment the counter
	for {
		// Ignore the document itself
		_, err := docsIter.Next()
		if err == iterator.Done {
			break
		}
		// If there is now error we add 1
		if err == nil {
			count++
		}

	}
	// Return the count
	return count
}


// Fetch a webhook using its ID. The webhook id has to be the same as the document id
func FetchWebhookWithID(id string) (structs.WebhookID, error) {
	//Create a client for the 
	client, err := GetFirestoreClient()
	defer client.Close()
	if err != nil{
		return structs.WebhookID{}, err;
	}

	var webhook structs.WebhookID
	iter := client.Collection(constants.FIRESTORE_COLLECTION).Documents(context.Background());

	//Loop through each document 
	for{
		//Get the document and check if it is done 
		doc, err := iter.Next()
		if err == iterator.Done {
			// Break if no more docs to get
			break
		}

		//if the ID is the same as the docs 
		if doc.Ref.ID == id {
			log.Println("Webhook found: " + id)
			err := doc.DataTo(&webhook);
			if err != nil{
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

//Fetch all webhooks
func FetchAllWebhooks() ([]structs.WebhookID, error){
	//Create a client
	client, err := GetFirestoreClient()
	defer client.Close()
	if err != nil{
		return nil, err;
	}

	//Iterate through all docs and decode them into the list of structs
	var webhooks []structs.WebhookID
	iter := client.Collection(constants.FIRESTORE_COLLECTION).Documents(context.Background())
	for{
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
		}else{
			// No error so we append the webhook
			webhooks = append(webhooks, webhook);
		}

	}
	
	// Returns either an empty list or a list of webhooks
	return webhooks, nil

}

// Delete a webhook from a given webhook id
// No error returns indicates that the process was okay 
func DeleteWebhook(webhookID string) error{
	// Get the client
	client, clientError := GetFirestoreClient()
	if clientError != nil{
		return clientError
	}


	// Delete the document based on the id given
	_ , err := client.Collection(constants.FIRESTORE_COLLECTION).Doc(webhookID).Delete(context.Background())
	if err != nil{
		return err
	}
	// No error and we return nil 
	return nil
}
