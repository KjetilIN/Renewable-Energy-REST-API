package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
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
	//Get any parameter received 
	givenParameters := strings.TrimPrefix(r.URL.Path, constants.NOTIFICATIONS_PATH)
	urlParts := strings.Split(givenParameters, "/")

	
	//Should only be given max one parameters 
	if (len(urlParts) > 1){
		log.Println("Error on get method for notification endpoint")
		http.Error(w, "Not correct usage of getting webhook information", http.StatusBadRequest)
		return;
	}else if (len(urlParts) == 1){
		//Fetch only the webhook with 
		fetchedWebhook, err := fetchWebhookWithID(urlParts[0])
		
		//Error on fetching the webhook with the id
		if(err != nil){
			http.Error(w, "No existing webhook with the ID: " + urlParts[0], http.StatusNotFound)
			return
		}

		//Encode the result
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


		//Encode the result
		encodeError := json.NewEncoder(w).Encode(&allWebHooks)
		if(encodeError != nil){
			log.Println("Error on encoding fetched webhooks", encodeError.Error())
			http.Error(w, "Error on encoding struct", http.StatusInternalServerError)
			return 
		}
		
	}

}




// HandlerNotifications is a handler for the /notifications endpoint.
func HandlerNotifications(w http.ResponseWriter, r *http.Request) {
	//Handle request based on the methods. 
	switch r.Method{
		case http.MethodPost:
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
