package main

import (
	"os"
	"testing"
)

const ERRORS = "TEST_ERRORS"
const PARAMETERS_MODE = "TEST_PARAMETERS_MODE"
const FORM_PARAMS = "TEST_FORM_PARAMS"

func check(name string) bool {
	return os.Getenv(name) == "true"
}

func skipIf(t *testing.T, name string) {
	if check(name) {
		t.SkipNow()
	}
}

func skipIfNot(t *testing.T, name string) {
	if !check(name) {
		t.SkipNow()
	}
}
