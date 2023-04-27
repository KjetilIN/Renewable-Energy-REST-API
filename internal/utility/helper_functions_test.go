package utility

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestPluralize tests if the correct pluralization of the word is returned.
func TestPluralize(t *testing.T) {
	tests := []struct {
		input    int    // input value for Pluralize function
		expected string // expected output from Pluralize function
	}{
		{1, ""},
		{0, "s"},
		{2, "s"},
		{100, "s"},
	}

	for _, tt := range tests {
		actual := Pluralize(tt.input)
		if actual != tt.expected {
			t.Errorf("Pluralize(%d): expected '%s', but got '%s'", tt.input, tt.expected, actual)
		}
	}
}

// TestDirChanger Tests the change directory function.
func TestDirChanger(t *testing.T) {
	// Retrieves the original working directory path.
	earlierPath, _ := os.Getwd()
	// Jumps two directories up.
	dirErr := DirChanger(1)
	if dirErr != nil {
		t.Fatal("Error when changing directory.")
	}
	// Checks the working directory.
	newPath, _ := os.Getwd()
	// Checks if the directory has changed.
	assert.NotEqual(t, earlierPath, newPath, "Path has not changed.")

	earlierPath = newPath
	dirErr = DirChanger(2)
	if dirErr != nil {
		t.Fatal("Error when changing directory.")
	}
	newPath, _ = os.Getwd()
	assert.NotEqual(t, earlierPath, newPath, "Path has not changed.")

	earlierPath = newPath
	dirErr = DirChanger(3)
	if dirErr != nil {
		t.Fatal("Error when changing directory.")
	}
	newPath, _ = os.Getwd()
	assert.NotEqual(t, earlierPath, newPath, "Path has not changed.")
}
