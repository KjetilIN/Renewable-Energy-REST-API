package structs

// Status This file defines structs to work with data.
type Status struct {
	CasesApi       int    `json:"countries_api"`
	NotificationDB int    `json:"notification_db"`
	Webhooks       int    `json:"webhooks"`
	Version        string `json:"version"`
	Uptime         int    `json:"uptime"`
}

// RenewableShareEnergyElement Struct to parse historical data into.
type RenewableShareEnergyElement struct {
	Name       string  `json:"name"`
	IsoCode    string  `json:"isoCode"`
	Year       int     `json:"year"`
	Percentage float64 `json:"percentage"`
}

// RenewableShareEnergyElementMean Struct to parse historical data into. Used when calculating mean percentage of countries over time.
type RenewableShareEnergyElementMean struct {
	Name       string  `json:"name"`
	IsoCode    string  `json:"isoCode"`
	Percentage float64 `json:"percentage"`
}

// Country Struct to collect information about countries.
type Country struct {
	Name        map[string]interface{} `json:"name"`
	CountryCode string                 `json:"cca3"`
	Borders     []string               `json:"borders"`
}
