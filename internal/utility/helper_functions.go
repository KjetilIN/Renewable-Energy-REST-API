package utility

import (
	"os"
	"path"
	"runtime"
)

// This file contains helper_functions used to increase the functionality of the service.

// Pluralize Helper function to pluralize words as appropriate.
// Returns: an empty string, or "s".
func Pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// DirChanger Changes the directory to project root.
func DirChanger(amountOfFolderToJump int) error {
	// Gets the filepath of file in question.
	_, filename, _, _ := runtime.Caller(0)
	var dir string

	switch amountOfFolderToJump {
	case 0:
		dir = path.Join(path.Dir(filename), ".")
		break
	case 1:
		dir = path.Join(path.Dir(filename), "..")
		break
	case 2:
		dir = path.Join(path.Dir(filename), "..", "..")
		break
	case 3:
		dir = path.Join(path.Dir(filename), "..", "..", "..")
		break
	case 4:
		dir = path.Join(path.Dir(filename), "..", "..", "..", "..")
		break
	}
	// Changes to the new dir structure.
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	return nil
}
