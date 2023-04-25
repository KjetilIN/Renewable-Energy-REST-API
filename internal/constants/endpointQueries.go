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

}

var NOTIFICATION_QUERIES = []structs.Information{
	
}

var CURRENT_QUERIES = []structs.Information{
	
}