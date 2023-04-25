package structs

import "time"

// Status This file defines structs to work with data.
type Status struct {
	CountriesApi   int    `json:"countries_api"`
	NotificationDB int    `json:"notification_db"`
	Webhooks       int    `json:"webhooks"`
	Version        string `json:"version"`
	Uptime         string `json:"uptime"`
	//AverageSystemLoad string `json:"average_system_load"`
	TotalMemoryUsage string `json:"total_memory_usage"`
}

type Webhook struct {
	Url     string `json:"url"`
	Country string `json:"country"`
	Calls   int    `json:"calls"`
	Event string `json:"event"`
}

type WebhookID struct {
	ID string `json:"webhook_id"`
	Webhook
	Created     time.Time `json:"created_timestamp"`
	Invocations int       `json:"invocations"`
}

// The call response for any given webhook 
type WebhookCallResponse struct{
	ID string `json:"webhook_id"`
	Webhook
	Invocations int `json:"invocations"`
	Message string `json:"message"`
}

// RenewableShareEnergyElement Struct to parse historical data into.
type RenewableShareEnergyElement struct {
	Name       string  `json:"name"`
	IsoCode    string  `json:"isoCode"`
	Year       int     `json:"year,omitempty"`
	Percentage float64 `json:"percentage"`
}

// Country Struct to collect information about countries.
type Country struct {
	Name        map[string]interface{} `json:"name"`
	CountryCode string                 `json:"cca3"`
	Borders     []string               `json:"borders"`
	Cache       time.Time              // Time in cache.
}

type IdResponse struct {
	ID string `json:"webhook_id"`
}
