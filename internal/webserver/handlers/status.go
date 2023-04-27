package handlers

import (
	"assignment-2/db"
	"assignment-2/internal/constants"
	"assignment-2/internal/utility"
	"assignment-2/internal/webserver/structs"
	"assignment-2/internal/webserver/uptime"
	"errors"
	"github.com/shirou/gopsutil/mem"
	"net/http"
	"strconv"
	"strings"
)

// Webhooks DB
var webhooks []structs.WebhookID

// Init empty list of webhooks
func InitWebhookRegistrations() {
	webhooks = []structs.WebhookID{}
}

// HTTP client
var client = &http.Client{}

// HandlerStatus is a handler for the /status endpoint.
func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	// Query for printing information about endpoint.
	if r.URL.Query().Has("information") && strings.Contains(strings.ToLower(r.URL.Query().Get("information")), "true") {
		_, writeErr := w.Write([]byte("To use API, remove ?information=true, from the URL.\n"))
		if writeErr != nil {
			return
		}
		utility.Encoder(w, constants.STATUS_QUERIES)
		return
	}

	// Set the content-type header to indicate that the response contains JSON data
	w.Header().Set("content-type", "application/json")

	// Return an error if the HTTP method is not GET.
	if r.Method != http.MethodGet {
		http.Error(w, errors.New("method is not supported. Currently only GET are supported").Error(), http.StatusMethodNotAllowed)
		return
	}

	// Get status information.
	status, err := getStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the status information as JSON and send it in the response.
	utility.Encoder(w, status)
}

func getStatus() (structs.Status, error) {
	// Check the status of the country API.
	url := constants.COUNTRIES_API_URL + "all"
	countryApiRequest, _ := http.NewRequest(http.MethodHead, url, nil)

	// Set the content-type header to indicate that the response contains JSON data
	countryApiRequest.Header.Set("content-type", "application/json")

	res, err := client.Do(countryApiRequest)
	if err != nil {
		return structs.Status{}, err
	}

	// Status code of the country API
	countriesApiStatus := res.StatusCode

	// If the status code is not 200, notify all subscribers to that event
	if countriesApiStatus != 200 {
		// Start a go routine for notifying all subscribers that they have been notified for the country api is down.
		go db.NotifyForEvent(constants.COUNTRY_API_EVENT, constants.FIRESTORE_COLLECTION)
	}

	// Firebase status
	dbStatus := db.CheckFirestoreConnection()

	var memUsage string
	defer func() {
		if r := recover(); r != nil {
			memUsage = "N/A"
		}
	}()

	// Get the memory usage in percent.
	memory, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}
	memUsage = strconv.Itoa(int(memory.UsedPercent))

	// Return a status struct containing information about the uptime and status of the notificationDB and countries APIs.
	return structs.Status{
		CountriesApi:   countriesApiStatus,
		NotificationDB: dbStatus,
		Webhooks:       db.GetNumberOfWebhooks(constants.FIRESTORE_COLLECTION),
		Version:        constants.VERSION,
		Uptime:         uptime.GetUptime(),
		//AverageSystemLoad: loadAvg + " in the last minute",
		TotalMemoryUsage: memUsage + "%",
	}, nil
}
