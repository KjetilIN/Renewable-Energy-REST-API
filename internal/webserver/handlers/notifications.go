package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

// Webhooks DB
var webhooks []structs.WebhookID;

// Init empty list of webhooks 
func InitWebhookRegistrations(){
	webhooks = []structs.WebhookID{};
}


// Fetch a webhook using its ID
func fetchWebhookWithID(id string) (structs.WebhookID, error) {
	for _, hook := range webhooks{
		if(hook.ID == id){
			return hook, nil;
		}
	}
	return structs.WebhookID{}, errors.New("Fetch with id: could not find")
}

//Fetch all webhooks
func fetchAllWebhooks() ([]structs.WebhookID, error){
	// At this stage we already have all the webhooks 
	return webhooks, nil
}


// Function for handling get request for the  
func handleGetMethod(w http.ResponseWriter, r *http.Request){
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
			http.Error(w, "Could not fetch all webhooks", http.StatusInternalServerError)
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




// HandlerNotifications is a handler for the /notifications endpoint.
func HandlerNotifications(w http.ResponseWriter, r *http.Request) {
	//Handle request based on the methods. 
	switch r.Method{
		case http.MethodPost:
			handlePostRequest(w,r);
			break
		case http.MethodDelete:
			break
		case http.MethodGet:
			//Handle get request to the endpoint based on given parameters 
			handleGetMethod(w,r)
			break
		default:
			// Only allowed methods
			http.Error(w, "Method " + r.Method + "not supported", http.StatusMethodNotAllowed)
			break
	}
}
