package utility

import (
	"os"
	"path"
	"runtime"
)

// dirChanger Changes the directory to project root.
func DirChanger() error {
	// Gets the filepath of history_test.go.
	_, filename, _, _ := runtime.Caller(0)
	// Jumps back 3 folders.
	dir := path.Join(path.Dir(filename), "..", "..", "..")
	// Changes to the new dir structure.
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	return nil
}
