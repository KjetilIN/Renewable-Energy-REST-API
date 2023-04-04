package handlers

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/structs"
	"assignment-2/internal/webserver/uptime"
	"encoding/json"
	"net/http"
)

// Webhooks DB
var webhooks []structs.WebhookID

// Init empty list of webhooks
func InitWebhookRegistrations() {
	webhooks = []structs.WebhookID{}
}

// Get number of webhooks
func GetNumberOfWebhooks() int {
	return len(webhooks)
}

// HandlerStatus is a handler for the /status endpoint.
// TODO: fill in notification_db URL
func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	// define dependent services and their URLs
	services := map[string]string{
		"country_api":     constants.COUNTRIES_API,
		"notification_db": "fill in firebase project url",
	}

	// check the status of each service
	status := make(map[string]int)
	for name, url := range services {
		resp, err := http.Get(url)
		if err != nil {
			status[name] = http.StatusInternalServerError
		} else {
			status[name] = resp.StatusCode
		}
	}

	// get number of registered webhooks
	numWebhooks := GetNumberOfWebhooks()

	// get uptime in seconds since service restart
	serviceUptime := uptime.GetUptime()

	// build response JSON
	response := structs.Status{
		CountryApi:     status["countries_api"],
		NotificationDB: status["notification_db"],
		Webhooks:       numWebhooks,
		Version:        "v1",
		Uptime:         serviceUptime,
	}

	// encode and send response JSON
	w.Header().Set("Content-Type", "application/json")
	if response.CountryApi != http.StatusOK || response.NotificationDB != http.StatusOK {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
