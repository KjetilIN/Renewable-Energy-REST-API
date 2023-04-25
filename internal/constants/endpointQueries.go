package constants

import "assignment-2/internal/webserver/structs"

var HISTORY_QUERIES = []structs.Information{
	{
		Title:   "No parameters",
		Example: "Use: /energy/v1/renewables/history",
		Description: "When you use endpoint with no parameters or queries, it will print all historical countries' mean value from" +
			"renewable share energy, which contains the percentage of renewable energy sources a given country has.",
	},
}

var STATUS_QUERIES = []structs.Information{
	{
		Title:   "No parameters",
		Example: "Use: /energy/v1/status",
		Description: "This endpoint are used with no parameters or queries, it will print all information about the " +
			"availability of all individual services this service depends on. The reporting occurs based on status " +
			"codes returned by the dependent services. The status interface further provides information about the " +
			"number of registered webhooks, and the uptime of the service. It also provides the total memory usage " +
			"of the computer in use.",
	},
}

var NOTIFICATION_QUERIES = []structs.Information{}

var CURRENT_QUERIES = []structs.Information{}
