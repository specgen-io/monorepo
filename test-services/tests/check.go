package main

import "os"

const ERRORS = "TEST_ERRORS"

func check(name string) bool {
	return os.Getenv(name) != "false"
}
