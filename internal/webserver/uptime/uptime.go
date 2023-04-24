package uptime

// This file contains provides the functionality for tracking the uptime of a service.
// Was made with inspiration from https://go.dev/play/p/by_nkvhzqD
import (
	"assignment-2/internal/utility"
	"fmt"
	"time"
)

var startTime time.Time

// Init initializes the uptime tracking by setting the startTime variable to the current time.
func Init() {
	startTime = time.Now()
}

// GetUptime returns the uptime of the service as a formatted string showing days, hours, minutes and seconds.
// Returns: A formatted string in the format of "x day(s), x hour(s), x minute(s), x second(s)".
func GetUptime() string {
	uptimeSeconds := int(time.Since(startTime).Seconds())
	uptimeMinutes := uptimeSeconds / 60
	uptimeHours := uptimeMinutes / 60
	uptimeDays := uptimeHours / 24

	var uptimeStr string
	if uptimeDays > 0 {
		uptimeStr += fmt.Sprintf("%d day%s, ", uptimeDays, utility.Pluralize(uptimeDays))
	}
	if uptimeHours%24 > 0 {
		uptimeStr += fmt.Sprintf("%d hour%s, ", uptimeHours%24, utility.Pluralize(uptimeHours%24))
	}
	if uptimeMinutes%60 > 0 {
		uptimeStr += fmt.Sprintf("%d minute%s, ", uptimeMinutes%60, utility.Pluralize(uptimeMinutes%60))
	}
	if uptimeSeconds%60 > 0 || uptimeStr == "" {
		uptimeStr += fmt.Sprintf("%d second%s", uptimeSeconds%60, utility.Pluralize(uptimeSeconds%60))
	}
	return uptimeStr
}
