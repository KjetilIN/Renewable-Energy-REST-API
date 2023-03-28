package structs

// This file defines structs to work with data.
type Status struct {
	CasesApi       int    `json:"countries_api"`
	NotificationDB int    `json:"notification_db"`
	Webhooks       int    `json:"webhooks"`
	Version        string `json:"version"`
	Uptime         int    `json:"uptime"`
}

// Struct to parse historical data into.
type HistoricalRSE struct {
	Name       string  `json:"name"`
	IsoCode    string  `json:"isoCode"`
	Percentage float64 `json:"percentage"`
}
