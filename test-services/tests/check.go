package main

import (
	"os"
	"testing"
)

const ERRORS = "TEST_ERRORS"
const PARAMETERS_MODE = "TEST_PARAMETERS_MODE"
const COMMA_SEPARATED_FORM_PARAMS_MODE = "TEST_COMMA_SEPARATED_FORM_PARAMS_MODE"
const NO_FORM_DATA = "TEST_NO_FORM_DATA"

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
