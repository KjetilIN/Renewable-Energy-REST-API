package db

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"context"
	"log"
	"net/http"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
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
