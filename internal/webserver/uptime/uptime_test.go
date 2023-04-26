package uptime

import (
	"testing"
	"time"
)

// TestInit checks if the upTime is initialized.
func TestInit(t *testing.T) {
	Init()
	if startTime.IsZero() {
		t.Errorf("Init did not set startTime")
	}
}

// TestGetUptime checks if the uptime of the service is a string in the correct format
func TestGetUptime(t *testing.T) {
	startTime = time.Now().Add(-time.Hour * 24 * 3).Add(-time.Minute * 15).Add(-time.Second * 30)
	expected := "3 days, 15 minutes, 30 seconds"
	if actual := GetUptime(); actual != expected {
		t.Errorf("GetUptime returned incorrect uptime: expected %s, got %s", expected, actual)
	}

	startTime = time.Now().Add(-time.Hour * 2).Add(-time.Minute * 30).Add(-time.Second * 15)
	expected = "2 hours, 30 minutes, 15 seconds"
	if actual := GetUptime(); actual != expected {
		t.Errorf("GetUptime returned incorrect uptime: expected %s, got %s", expected, actual)
	}

	startTime = time.Now().Add(-time.Minute * 45).Add(-time.Second * 10)
	expected = "45 minutes, 10 seconds"
	if actual := GetUptime(); actual != expected {
		t.Errorf("GetUptime returned incorrect uptime: expected %s, got %s", expected, actual)
	}

	startTime = time.Now().Add(-time.Second * 10)
	expected = "10 seconds"
	if actual := GetUptime(); actual != expected {
		t.Errorf("GetUptime returned incorrect uptime: expected %s, got %s", expected, actual)
	}

	startTime = time.Now()
	expected = "0 seconds"
	if actual := GetUptime(); actual != expected {
		t.Errorf("GetUptime returned incorrect uptime: expected %s, got %s", expected, actual)
	}
}
