package app

import (
	"path/filepath"
	"runtime"
)

// RootPath always returns the root path of this module. This utility function is intended to help during testing when
// specific file paths relative to the root of the project are needed.
func RootPath() string {
	_, file, _, _ := runtime.Caller(0)

	// Root folder of this project
	return filepath.Join(filepath.Dir(file), "..", "..")
}
