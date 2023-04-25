package constants

import "assignment-2/internal/webserver/structs"

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
)

// HistoryEndpointInformation Method to print information about endpoint
func HistoryEndpointInformation() []structs.Information {
	var historyInformation []structs.Information

	infoNoParam := structs.Information{
		Title:   "No parameters",
		Example: "Use: /energy/v1/renewables/history",
		Description: "When you use endpoint with no parameters or queries, it will print all historical countries' mean value from" +
			"renewable share energy, which contains the percentage of renewable energy sources a given country has.",
	}
	historyInformation = append(historyInformation, infoNoParam)

	infoOptionalParam := structs.Information{
		Title:       "Country code/name",
		Example:     "Use: /energy/v1/renewables/history/nor or /norway",
		Description: "When adding another slash and a country code or name, all data from that country will be printed.",
	}
	historyInformation = append(historyInformation, infoOptionalParam)

	infoBeginQuery := structs.Information{
		Title:       "Begin query",
		Example:     "Use: /energy/v1/renewables/history?begin=2010",
		Description: "Using the begin query, it will give all countries' historical record from that year.",
	}
	historyInformation = append(historyInformation, infoBeginQuery)

	infoEndQuery := structs.Information{
		Title:       "End query",
		Example:     "Use: /energy/v1/renewables/history?end=1970",
		Description: "Using the end query, it will give all countries' historical records up until to that year.",
	}
	historyInformation = append(historyInformation, infoEndQuery)

	infoBeginEndQuery := structs.Information{
		Title:       "Begin & End query",
		Example:     "Use: /energy/v1/renewables/history?begin=1996&end=2002",
		Description: "Using both begin and end, will return the mean of all countries between the specific year.",
	}
	historyInformation = append(historyInformation, infoBeginEndQuery)

	infoYearQuery := structs.Information{
		Title:       "Year query",
		Example:     "Use: /energy/v1/renewables/history?year=2020",
		Description: "Using the year query, it will return all countries' specific year written.",
	}
	historyInformation = append(historyInformation, infoYearQuery)

	infoMeanQuery := structs.Information{
		Title:   "Mean query",
		Example: "Use: /energy/v1/renewables/history?mean=true",
		Description: "The mean query, will not do anything if used without other queries such as country code, or being/end." +
			"\nIf used it will calculate the mean of the data which would be returned.",
	}
	historyInformation = append(historyInformation, infoMeanQuery)

	infoQuery := structs.Information{
		Title:       "Information query",
		Example:     "Use: /energy/v1/renewables/history?information=true",
		Description: "Will present information about endpoint.",
	}
	historyInformation = append(historyInformation, infoQuery)

	sortingbyValQuery := structs.Information{
		Title:   "Sort by value query",
		Example: "Use: /energy/v1/renewables/history?sortbyvalue=true",
		Description: "Will sort data ascending by percentage. An additional parameter &descending may also be used to sort" +
			"it descending.",
	}
	historyInformation = append(historyInformation, sortingbyValQuery)

	sortAlphabeticallyQuery := structs.Information{
		Title:   "Sort alphabetically query",
		Example: "Use: /energy/v1/renewables/history?sortalphabetically=true",
		Description: "Will sort data ascending alphabetically. An additional parameter &descending may also be used to sort" +
			"it descending.",
	}
	historyInformation = append(historyInformation, sortAlphabeticallyQuery)

	return historyInformation
}
