package handlers

import (
	"assignment-2/db"
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/iterator"
)

// Webhooks DB
var webhooks []structs.WebhookID;

// Init empty list of webhooks 
func InitWebhookRegistrations(){
	webhooks = []structs.WebhookID{};
}

//Get number of webhooks. 
// Note that if the service is down there will be not handled this function. 
// The user has to see the status endpoint
func GetNumberOfWebhooks() int{
	//Create a client for the 
	client, err := db.GetFirestoreClient()
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
func fetchWebhookWithID(id string) (structs.WebhookID, error) {
	//Create a client for the 
	client, err := db.GetFirestoreClient()
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
func fetchAllWebhooks() ([]structs.WebhookID, error){
	//Create a client
	client, err := db.GetFirestoreClient()
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
	


// Function for handling get request for the  
func handleGetRequest(w http.ResponseWriter, r *http.Request){
	// Get any parameter received
	givenParameters := strings.TrimPrefix(r.URL.Path, constants.NOTIFICATIONS_PATH)
	givenParameters = strings.ReplaceAll(givenParameters, " ", "")
	urlParts := strings.Split(givenParameters, "/")

	// Remove empty strings from urlParts slice
	var params []string
	for _, part := range urlParts {
		if part != "" {
			params = append(params, part)
		}
	}
	
	//Should only be given max one parameters 
	if (len(params) > 1){
		log.Println("Error on get method for notification endpoint")
		http.Error(w, "Not correct usage of getting webhook information", http.StatusBadRequest)
		return;
	}else if (len(params) != 0 ){
		//Fetch only the webhook with 
		fetchedWebhook, err := fetchWebhookWithID(params[0])
		
		//Error on fetching the webhook with the id
		if(err != nil){
			log.Println("Error on fetching webhook with id: " + err.Error())
			http.Error(w, "No existing webhook with the ID: " + params[0], http.StatusNotFound)
			return
		}

		//Encode the result
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encodeError := json.NewEncoder(w).Encode(&fetchedWebhook)
		if(encodeError != nil){
			log.Println("Error on encoding fetched webhook with ID", encodeError.Error())
			http.Error(w, "Error on encoding struct", http.StatusInternalServerError)
			return 
		}
			
	}else{
		// Fetch all webhooks 
		allWebHooks, fetchError := fetchAllWebhooks();

		//Handle error 
		if(fetchError != nil){
			log.Println("Error on fetching all webhooks: ", fetchError.Error())
			http.Error(w, "Could not fetch all webhooks. \nSee status endpoint to the status of webhook database...", http.StatusInternalServerError)
			return 
		}

		//Return message if there are no webhooks created yet
		if(len(allWebHooks) == 0){
			log.Println("No content in webhooks")
			http.Error(w, "No webhooks in storage", http.StatusNoContent)
			return
		}


		//Encode the result
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encodeError := json.NewEncoder(w).Encode(&allWebHooks)
		if(encodeError != nil){
			log.Println("Error on encoding fetched webhooks", encodeError.Error())
			http.Error(w, "Error on encoding struct", http.StatusInternalServerError)
			return 
		}
		
	}

}

// Create a random webhook id based on the content and time of the creation 
func createWebhookID(webhook structs.Webhook) string{
	// Secret that the hash generation is based on
	secret := []byte(time.Now().Local().String());
	hashContent := []byte(webhook.Url + webhook.Country)
	hash := hmac.New(sha256.New, secret)
	hash.Write(hashContent);
	return hex.EncodeToString(hash.Sum(nil));
}

func handlePostRequest(w http.ResponseWriter, r *http.Request){
	// Expects incoming body to be in correct format, so we encode it directly to a struct 
	givenHook := structs.Webhook{}
	err := json.NewDecoder(r.Body).Decode(&givenHook)
	if err != nil {
		// Was not in the correct format
		log.Print("Something went wrong: "+err.Error())
		http.Error(w, "Error: given body does not fit the schema", http.StatusBadRequest)
		return
	}

	// Create an ID and add it to the list of hooks
	id := createWebhookID(givenHook)
	webhook := structs.WebhookID{ID:id, Webhook: givenHook}
	webhooks = append(webhooks, webhook)
	
	// Logging that a new webhook has been
	log.Println("Webhook " + webhook.Url + " has been registered.")


	// Encode the Response object to JSON format, and handle any errors 
	resp := structs.IdResponse{ID: webhook.ID}
    jsonResponse, err := json.Marshal(&resp)
    if err != nil {
		log.Println("Error on marshal response: " + err.Error())
        http.Error(w, "Something went wrong when returning the webhook ID", http.StatusInternalServerError)
        return
    }

	// Set the output to be correct
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func handleDeleteRequest(w http.ResponseWriter, r *http.Request){
	// Get any parameter received
	givenParameters := strings.TrimPrefix(r.URL.Path, constants.NOTIFICATIONS_PATH)
	givenParameters = strings.ReplaceAll(givenParameters, " ", "")
	urlParts := strings.Split(givenParameters, "/")

	// Remove empty strings from urlParts slice
	var params []string
	for _, part := range urlParts {
		if part != "" {
			params = append(params, part)
		}
	}

	//Should always be one parameter left.
	if(len(params) != 1){
		log.Println("Error on params, should have been 1 was: " + strconv.Itoa(len(params)))
		http.Error(w, "Bad request, please add a single ID", http.StatusBadRequest)
		return
	}

	//Deleting the webhook with the given ID, or not doing anything. 
	givenID := params[0]
	newWebhooks := []structs.WebhookID{}
	isDeleted := false
	for _, hook := range webhooks{
		if(hook.ID !=  givenID){
			newWebhooks = append(newWebhooks, hook);
		}else{
			//Indicate that the webhook was found and deleted. 
			isDeleted = true;
		}
	}
	webhooks = newWebhooks

	if(isDeleted){
		//Tell the end user that the webhook was deleted
		http.Error(w, "Webhook with ID: " + givenID + " is deleted", http.StatusOK)
		return 
	}else{
		//No webhook was found
		http.Error(w, "Webhook not stored, did not delete a webhook", http.StatusOK)
		return 
	}


}




// HandlerNotifications is a handler for the /notifications endpoint.
func HandlerNotifications(w http.ResponseWriter, r *http.Request) {
	//Handle request based on the methods. 
	switch r.Method{
		case http.MethodPost:
			//Creating a new webhook
			handlePostRequest(w,r)
			break
		case http.MethodDelete:
			//Deleting a webhook with the given ID
			handleDeleteRequest(w,r)
			break
		case http.MethodGet:
			//Handle get request to the endpoint based on given parameters 
			handleGetRequest(w,r)
			break
		default:
			// Only allowed methods
			http.Error(w, "Method " + r.Method + "not supported", http.StatusMethodNotAllowed)
			break
	}
}
