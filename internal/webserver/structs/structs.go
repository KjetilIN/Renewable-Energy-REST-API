package structs

// This file defines structs to work with data.
type Status struct {
	CasesApi       int    `json:"countries_api"`
	NotificationDB int    `json:"notification_db"`
	Webhooks       int    `json:"webhooks"`
	Version        string `json:"version"`
	Uptime         int    `json:"uptime"`
}

type Webhook struct{
	Url string `json:"url"`
	Country string `json:"country"`
	Calls int `json:"calls"`
}

type WebhookID struct{
	ID string `json:"webhook_id"`
	Webhook
}