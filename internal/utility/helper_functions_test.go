package utility

import "testing"

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
