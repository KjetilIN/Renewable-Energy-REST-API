package webserver

import (
	"assignment-2/internal/constants"
	"assignment-2/internal/webserver/handlers"
	"assignment-2/internal/webserver/mock"
	"assignment-2/internal/webserver/uptime"
	"log"
	"net/http"
)

// InitServer sets up handler endpoints and starts the HTTP-server
func InitServer() {

	// Points the different URL-paths to the correct handler
	http.HandleFunc(constants.DEFAULT_PATH, handlers.HandlerDefault)
	http.HandleFunc(constants.CURRENT_PATH, handlers.HandlerCurrent)
	http.HandleFunc(constants.HISTORY_PATH, handlers.HandlerHistory)
	http.HandleFunc(constants.STATUS_PATH, handlers.HandlerStatus)
	http.HandleFunc(constants.NOTIFICATIONS_PATH, handlers.HandlerNotifications)

	// Points the different URL-paths to the correct stubHandler
	http.HandleFunc(constants.MOCK_CURRENT_API_URL, mock.StubHandlerCurrent)
	http.HandleFunc(constants.MOCK_HISTORY_API_URL, mock.StubHandlerHistory)

	// Starting HTTP-server
	log.Println("Starting server on port " + constants.DEFAULT_PORT + " ...")
	uptime.Init()
	log.Fatal(http.ListenAndServe(":"+constants.DEFAULT_PORT, nil))
}
