package uptime

// This file contains provides the functionality for tracking the uptime of a service.
// Was made with inspiration from https://go.dev/play/p/by_nkvhzqD
import (
	"strconv"
	"time"
)

var startTime time.Time

// Init initializes the uptime tracking by setting the startTime variable to the current time.
func Init() {
	startTime = time.Now()
}

// GetUptime returns the uptime of the service in seconds.
func GetUptime() int {
	return int(time.Since(startTime).Seconds())
}

// ConvertUptime converts the given uptime in seconds to a human-readable string.
func ConvertUptime(uptimeInSeconds int) string {
	if uptimeInSeconds < 60 {
		return strconv.Itoa(uptimeInSeconds) + " seconds"
	}

	uptimeInMinutes := uptimeInSeconds / 60
	if uptimeInMinutes < 60 {
		return strconv.Itoa(uptimeInMinutes) + " minutes"
	}

	uptimeInHours := uptimeInMinutes / 60
	if uptimeInHours < 24 {
		return strconv.Itoa(uptimeInHours) + " hours"
	}

	uptimeInDays := uptimeInHours / 24
	return strconv.Itoa(uptimeInDays) + " days"
}
