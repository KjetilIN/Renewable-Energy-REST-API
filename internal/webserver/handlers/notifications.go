package handlers

import (
	"assignment-2/db"
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Function for handling get request for the  notification endpoint
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
		fetchedWebhook, err := db.FetchWebhookWithID(params[0], constants.FIRESTORE_COLLECTION)
		
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
		allWebHooks, fetchError := db.FetchAllWebhooks(constants.FIRESTORE_COLLECTION);

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

	// Format the webhook for the firestore and add it to the database as a new document 
	id := createWebhookID(givenHook)
	formattedWebhook := structs.WebhookID{ID: id, Webhook: givenHook, Created: time.Now()}
	db.AddWebhook(formattedWebhook, constants.FIRESTORE_COLLECTION)
	
	// Logging that a new webhook has been
	log.Println("Request for adding webhook with url: " + formattedWebhook.Url)

	// Encode the Response object to JSON format, and handle any errors 
	resp := structs.IdResponse{ID: formattedWebhook.ID}
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

// Function for handling the get request 
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
	
	deletedError := db.DeleteWebhook(givenID)

	if(deletedError == nil){
		//Tell the end user that the webhook was deleted
		http.Error(w, "Webhook with ID: " + givenID + " has been deleted", http.StatusOK)
		return 
	}else{
		//No webhook was found
		log.Print("Error on deletion of webhook: " + deletedError.Error())
		http.Error(w, "Webhook was not deleted. Internal error...", http.StatusInternalServerError)
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
