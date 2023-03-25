package uptime

// This file contains provides the functionality for tracking the uptime of a service.
// Was made with inspiration from https://go.dev/play/p/by_nkvhzqD
import "time"

var startTime time.Time

// Init initializes the uptime tracking by setting the startTime variable to the current time.
func Init() {
	startTime = time.Now()
}

// GetUptime returns the uptime of the service in seconds.
func GetUptime() int {
	return int(time.Since(startTime).Seconds())
}
