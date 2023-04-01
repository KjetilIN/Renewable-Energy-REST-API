package handlers

import (
	"assignment-2/internal/constants"
	"net/http"
	"log"
	"strings"
)



func handleGetMethod(w http.ResponseWriter, r *http.Request){
	//Get any parameter received 
	givenParameters := strings.TrimPrefix(r.URL.Path, constants.NOTIFICATIONS_PATH);
	urlParts := strings.Split(givenParameters, "/");	
	
	//Should only be given max one parameters 
	if (len(urlParts) > 1){
		log.Println("Error on get method for notification endpoint");
		http.Error(w, "Not correct usage of getting webhook information", http.StatusBadRequest);
		return;
	}else{
		//Get the webhook id
		var ID string = "";
		if (len(urlParts) == 1){
			ID = urlParts[0];
		}
		log.Println("ID is: " + ID)
	}

}




// HandlerNotifications is a handler for the /notifications endpoint.
func HandlerNotifications(w http.ResponseWriter, r *http.Request) {
	//Handle request based on the methods. 
	switch r.Method{
		case http.MethodPost: 
			break;
		case http.MethodDelete:
			break;
		case http.MethodGet:
			handleGetMethod(w,r);
			break;

		default:
			// Only allowed methods
			http.Error(w, "Method " + r.Method + "not supported", http.StatusMethodNotAllowed)
			break;
	}
}
