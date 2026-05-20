package main

import (
	"os"
)

// exists checks if the given filename exists, using os.Stat
func exists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
