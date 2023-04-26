package handlers

import (
	"assignment-2/db"
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/structs"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Function for handling get request for the  notification endpoint
func handleGetRequest(w http.ResponseWriter, r *http.Request){
	// Get any parameter received
	component, urlError := utility.GetOneFirstComponentOnly(constants.NOTIFICATIONS_PATH, r.URL.Path)
	if urlError != nil{
		log.Println("Utility function error for getting one component: " + urlError.Error())
		http.Error(w, "Bad request: Endpoint was not correctly used. See readme...", http.StatusBadRequest)
		return 
	} 
	
	
	// Check if there was given a component 
	if component != "" {
		//Fetch only the webhook with 
		fetchedWebhook, err := db.FetchWebhookWithID(component, constants.FIRESTORE_COLLECTION)
		
		//Error on fetching the webhook with the id
		if err != nil{
			log.Println("Error on fetching webhook with id: " + err.Error())
			http.Error(w, "No existing webhook with the ID: " + component, http.StatusNotFound)
			return
		}

		//Encode the result
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Use encoder for the result 
		utility.Encoder(w, fetchedWebhook)
			
	}else{
		// Fetch all webhooks 
		allWebHooks, fetchError := db.FetchAllWebhooks(constants.FIRESTORE_COLLECTION);

		//Handle error 
		if fetchError != nil{
			log.Println("Error on fetching all webhooks: ", fetchError.Error())
			http.Error(w, "Could not fetch all webhooks. \nSee status endpoint to the status of webhook database...", http.StatusInternalServerError)
			return 
		}

		//Return message if there are no webhooks created yet
		if len(allWebHooks) == 0 {
			log.Println("No content in webhooks")
			http.Error(w, "No webhooks in storage", http.StatusNoContent)
			return
		}


		//Encode the result
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		utility.Encoder(w, allWebHooks)
		
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

	// Response object to JSON format, and handle any errors 
	resp := structs.IdResponse{ID: formattedWebhook.ID}

	// Set the output to be correct
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	
	utility.Encoder(w, resp)
}

// Function for handling the get request 
func handleDeleteRequest(w http.ResponseWriter, r *http.Request){

	// Get the id of from the path, and check if its done correctly 
	id , urlError := utility.GetOneFirstComponentOnly(constants.NOTIFICATIONS_PATH, r.URL.Path)
	if urlError != nil{
		log.Println("Utility function error for getting one component: " + urlError.Error())
		http.Error(w, "Bad request, please add a single ID", http.StatusBadRequest)
		return 
	} 

	// No id was given
	if id == ""{
		//Tell the endpoint was not correctly used 
		log.Println("No id was given")
		http.Error(w, "Error: no id was given. See readme for usage...", http.StatusBadRequest)
		return 
	}
	
	// Delete the webhook using the firebase methods
	deletedError := db.DeleteWebhook(id, constants.FIRESTORE_COLLECTION)
	if deletedError == nil{
		//Tell the end user that the process of deletion attempt did not lead to any errors. 
		http.Error(w, "Webhook has been deleted successfully", http.StatusOK)
		return 
	}else{
		//No webhook was found
		log.Print("Error on deletion of webhook: " + deletedError.Error())
		http.Error(w, "Internal error on attempt for deleting the deleting the webhook....", http.StatusInternalServerError)
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
