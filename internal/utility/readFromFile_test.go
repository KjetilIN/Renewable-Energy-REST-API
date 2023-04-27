package utility

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"runtime"
	"testing"
)

// TestRSEToJSON Tests the read from file function.
func TestRSEToJSON(t *testing.T) {
	// Change directory.
	dirChangeErr := dirChanger()
	if dirChangeErr != nil {
		t.Fatal("Changing directory failed.")
	}
	readFromFile, err := RSEToJSON()
	if err != nil {
		t.Fatal("Error reading file.")
	}
	// Checks if list read from file is empty.
	assert.NotEmpty(t, readFromFile, "List is empty!")

}

// dirChanger Changes the directory to project root.
func dirChanger() error {
	// Gets the filepath of history_test.go.
	_, filename, _, _ := runtime.Caller(0)
	// Jumps back 3 folders.
	dir := path.Join(path.Dir(filename), "..", "..")
	// Changes to the new dir structure.
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	return nil
}
