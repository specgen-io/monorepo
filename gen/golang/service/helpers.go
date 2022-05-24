package service

import (
	"path/filepath"
)

//TODO: Kill this - replace with Module
func createPath(args ...string) string {
	parts := []string{}
	for _, arg := range args {
		if arg != "" {
			parts = append(parts, arg)
		}
	}
	return filepath.Join(parts...)
}
