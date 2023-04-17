package constants

// This file defines constants used throughout the program.
const (
	// PORT Default port. If the port is not set by environment variables, set the port.
	DEFAULT_PORT = "8080"

	// The paths that will be handled by each handler
	DEFAULT_PATH       = "/energy/"
	CURRENT_PATH       = "/energy/v1/renewables/current/"
	HISTORY_PATH       = "/energy/v1/renewables/history/"
	STATUS_PATH        = "/energy/v1/status/"
	NOTIFICATIONS_PATH = "/energy/v1/notifications/"
	FIRESTORE_COLLECTION = "webhooks"
	FIRESTORE_COLLECTION_TEST = "test_collection"
	FIREBASE_CREDENTIALS_FILE = "cloud-assignment-2.json"
	MAX_WEBHOOK_COUNT = 40 // The number of max webhooks that are
)
