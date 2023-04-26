package constants

import "assignment-2/internal/webserver/structs"

var HISTORY_QUERIES = []structs.Information{
	{
		Title:   "No parameters",
		Example: "Use: /energy/v1/renewables/history",
		Description: "When you use endpoint with no parameters or queries, it will print all historical countries' mean value from" +
			"renewable share energy, which contains the percentage of renewable energy sources a given country has.",
	},

	{
		Title:       "Country code/name",
		Example:     "Use: /energy/v1/renewables/history/nor or /norway",
		Description: "When adding another slash and a country code or name, all data from that country will be printed.",
	},

	{
		Title:       "Begin query",
		Example:     "Use: /energy/v1/renewables/history?begin=2010",
		Description: "Using the begin query, it will give all countries' historical record from that year.",
	},

	{
		Title:       "End query",
		Example:     "Use: /energy/v1/renewables/history?end=1970",
		Description: "Using the end query, it will give all countries' historical records up until to that year.",
	},

	{
		Title:       "Begin & End query",
		Example:     "Use: /energy/v1/renewables/history?begin=1996&end=2002",
		Description: "Using both begin and end, will return the mean of all countries between the specific year.",
	},

	{
		Title:       "Year query",
		Example:     "Use: /energy/v1/renewables/history?year=2020",
		Description: "Using the year query, it will return all countries' specific year written.",
	},

	{
		Title:   "Mean query",
		Example: "Use: /energy/v1/renewables/history?mean=true",
		Description: "The mean query, will not do anything if used without other queries such as country code, or being/end." +
			"\nIf used it will calculate the mean of the data which would be returned.",
	},

	{
		Title:       "Information query",
		Example:     "Use: /energy/v1/renewables/history?information=true",
		Description: "Will present information about endpoint.",
	},

	{
		Title:   "Sort by value query",
		Example: "Use: /energy/v1/renewables/history?sortbyvalue=true",
		Description: "Will sort data ascending by percentage. An additional parameter &descending may also be used to sort" +
			"it descending.",
	},

	{
		Title:   "Sort alphabetically query",
		Example: "Use: /energy/v1/renewables/history?sortalphabetically=true",
		Description: "Will sort data ascending alphabetically. An additional parameter &descending may also be used to sort" +
			"it descending.",
	},
}

var STATUS_QUERIES = []structs.Information{
	{
		Title:   "No parameters",
		Example: "Use: /energy/v1/status",
		Description: "Will print all information about the availability of all individual services this service " +
			"depends on. The reporting occurs based on status codes returned by the dependent services. Furthermore" +
			" information about the number of registered webhooks, the uptime of the service, and the total memory " +
			"usage of the computer in use.",
	},

	{
		Title:       "Information query",
		Example:     "Use: /energy/v1/status?information=true",
		Description: "Will present information about endpoint.",
	},
}

var NOTIFICATION_QUERIES = []structs.Information{
	{
		Title:   "No parameters",
		Example: "Use: /energy/v1/notification",
		Description: "When you use endpoint with no parameters or queries, with a GET request, it will print a list " +
			"of all webhooks. Could also be empty if non are registered yet.",
	},

	{
		Title:       "Information query",
		Example:     "Use: /energy/v1/notification?information=true",
		Description: "Will present information about endpoint.",
	},

	{
		Title:   "Add Notification",
		Example: "Use: /energy/v1/notification",
		Description: "When you use endpoint with no parameters or queries, with a POST request and correct body, " +
			"it will add a notification to the database.",
	},

	{
		Title:   "Delete notification",
		Example: "Use: /energy/v1/notifications/{webhook_id}",
		Description: "When you use endpoint with this optional parameter, with a DELETE request and including the ID " +
			"of the webhook in the url it delete the notification.",
	},

	{
		Title:   "Get notification",
		Example: "Use: /energy/v1/notifications/{webhook_id}",
		Description: "When you use endpoint with this optional parameter, with a GET request and including the ID " +
			"of the webhook in the url it return the body of the notification.",
	},
}

var CURRENT_QUERIES = []structs.Information{
	{
		Title:   "No parameters",
		Example: "Use: /energy/v1/renewables/current",
		Description: "When you use endpoint with no parameters or queries, it will print all countries of the " +
			"year=2021, to the client. The year found is based on the highest year found in the csv file.",
	},

	{
		Title:       "Country code/name",
		Example:     "Use: /energy/v1/renewables/current/nor or /norway",
		Description: "When adding another slash and a country code or name, all data from that country will be printed.",
	},

	{
		Title:   "Neighbour parameter",
		Example: "Use: /energy/v1/renewables/current?neighbours=bool?",
		Description: "It is an optional query parameter which will print information about the neighbouring " +
			"countries of the country passed. It therefore, dependent on the optional parameter: country.",
	},

	{
		Title:       "Information query",
		Example:     "Use: /energy/v1/renewables/current?information=true",
		Description: "Will present information about endpoint.",
	},

	{
		Title:   "Sort by value query",
		Example: "Use: /energy/v1/renewables/current?sortbyvalue=true",
		Description: "Will sort data ascending by percentage. An additional parameter &descending may also be used to sort" +
			"it descending.",
	},

	{
		Title:   "Sort alphabetically query",
		Example: "Use: /energy/v1/renewables/current?sortalphabetically=true",
		Description: "Will sort data ascending alphabetically. An additional parameter &descending may also be used to sort" +
			"it descending.",
	},
}

var DEFAULT_QUERIES = []structs.Information{
	{
		Title:       "No parameters",
		Example:     "Use: /energy/",
		Description: "Will provide an overview of the functionality of the web application",
	},

	{
		Title:       "Information query",
		Example:     "Use: /energy/?information=true",
		Description: "Will present information about endpoint.",
	},
}
