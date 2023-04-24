package constants

// This file defines constants used throughout the program.
const (
	// DEFAULT_PORT PORT Default port. If the port is not set by environment variables, set the port.
	DEFAULT_PORT = "8080"

	// The paths that will be handled by each handler
	DEFAULT_PATH = "/energy/"
	CURRENT_PATH = "/energy/v1/renewables/current/"
	HISTORY_PATH = "/energy/v1/renewables/history/"

	// LIMIT_CACHE_TIME The default limit of how long entry can be stored in cache.
	LIMIT_CACHE_TIME = 600

	STATUS_PATH               = "/energy/v1/status/"
	NOTIFICATIONS_PATH        = "/energy/v1/notifications/"
	FIRESTORE_COLLECTION      = "webhooks"                // Name of the main collection for the webhooks
	FIRESTORE_COLLECTION_TEST = "test_collection"         // Name of the collection which the test make use of
	FIREBASE_CREDENTIALS_FILE = "cloud-assignment-2.json" // Name of the credential file, see readme for how to use and where to place
	MAX_WEBHOOK_COUNT         = 40                        //

	COUNTRIES_API_URL  = "http://129.241.150.113:8080/v3.1/"
	NOTIFICATIONDB_URL = ""

	COUNTRYCODE_API_ADDRESS = "https://restcountries.com/v3/alpha/"
	COUNTRYNAME_API_ADDRESS = "https://restcountries.com/v3/name/"

	// ASCENDING Used to address way of sorting.
	ASCENDING = 1
	// DESCENDING Used to address way of sorting.
	DESCENDING = 2
)
