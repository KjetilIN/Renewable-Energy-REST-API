package structs

// This file defines structs to work with data.
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
}

type WebhookID struct {
	ID string `json:"webhook_id"`
	Webhook
}

type IdResponse struct {
	ID string `json:"webhook_id"`
}
