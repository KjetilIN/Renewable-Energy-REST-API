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

	STATUS_PATH        = "/energy/v1/status/"
	NOTIFICATIONS_PATH = "/energy/v1/notifications/"

	FIRESTORE_COLLECTION      = "webhooks"        // Name of the main collection for the webhooks
	FIRESTORE_COLLECTION_TEST = "test_collection" // Name of the collection which the test make use of
	FIREBASE_FILE_NAME = "cloud-assignment-2.json" // Name of the firebase credentials file 
	FIREBASE_CREDENTIALS_FILE_PATH = "./" + FIREBASE_FILE_NAME  // Path of credentials in the root folder
	FIREBASE_CREDENTIALS_FILE_PATH_FOR_DB_TESTS = "../" + FIREBASE_FILE_NAME // Same credentials but different locations
	MAX_WEBHOOK_COUNT = 40 // The max amount of webhooks allowed

	// COUNTRIES_API_URL Used to send head request.
	COUNTRIES_API_URL = "http://129.241.150.113:8080/v3.1/"

	// COUNTRYCODE_API_ADDRESS URL for GET request of country code.
	COUNTRYCODE_API_ADDRESS = "http://129.241.150.113:8080/v3.1/alpha/"
	// COUNTRYNAME_API_ADDRESS URL for GET request of country names.
	COUNTRYNAME_API_ADDRESS = "http://129.241.150.113:8080/v3.1/name/"

	// ASCENDING Used to address way of sorting.
	ASCENDING = 1
	// DESCENDING Used to address way of sorting.
	DESCENDING = 2
)

// The different events types. See readme for more details
const (
	COUNTRY_API_EVENT = "COUNTRY_DOWN" // Event: country api is down
	CALLS_EVENT       = "CALLS"        // Event: based on invocations
	PURGE_EVENT       = "PURGE"        // Event: webhooks has been purged
)
