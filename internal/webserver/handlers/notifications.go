package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"encoding/json"
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


func handleGetMethod(w http.ResponseWriter, r *http.Request){
	//Get any parameter received 
	givenParameters := strings.TrimPrefix(r.URL.Path, constants.NOTIFICATIONS_PATH)
	urlParts := strings.Split(givenParameters, "/")
	
	//Should only be given max one parameters 
	if (len(urlParts) > 1){
		log.Println("Error on get method for notification endpoint")
		http.Error(w, "Not correct usage of getting webhook information", http.StatusBadRequest)
		return;
	}else{
		//Get the webhook if there is an id 
		if (len(urlParts) == 1){

			//An id was given 
			ID := urlParts[0]

			//Loop through and try to find the webhook with the given id
			var foundWebhook structs.WebhookID
			for _, hook := range webhooks{
				if(hook.ID == ID){
					foundWebhook = hook
					break // found the hook, and break out of loop 
				}
			}

			//check if we found the webhook
			if(len(foundWebhook.ID) == 0 ){
				err := json.NewEncoder(w).Encode(&foundWebhook)
				if(err != nil){
					return
				}

			}else{
				//Encode the result 
				err := json.NewEncoder(w).Encode(&foundWebhook)
				if(err != nil){
					return
				}
			}

			
		}else{

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
			handleGetMethod(w,r)
			break

		default:
			// Only allowed methods
			http.Error(w, "Method " + r.Method + "not supported", http.StatusMethodNotAllowed)
			break
	}
}
