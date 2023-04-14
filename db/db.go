package db

import (
	"assignment-2/internal/constants"
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
	sa := option.WithCredentialsFile("path/to/serviceAccount.json")
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
