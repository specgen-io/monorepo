package main

import (
	"encoding/json"
	"github.com/yalp/jsonpath"
	"gotest.tools/assert/cmp"
	"reflect"
	"strings"
	"testing"
)

func assertEqualJson(t *testing.T, actual []byte, expected string) {
	if strings.HasPrefix(strings.TrimSpace(expected), "[") {
		actualData := []interface{}{}
		err := json.Unmarshal(actual, &actualData)
		if err != nil {
			t.Fatalf(`failed to read actual json: "%s", error: %s`, string(actual), err.Error())
			return
		}

		expectedData := []interface{}{}
		err = json.Unmarshal([]byte(expected), &expectedData)
		if err != nil {
			t.Fatalf(`failed to read actual json: "%s", error: %s`, string(actual), err.Error())
			return
		}

		if !reflect.DeepEqual(actualData, expectedData) {
			t.Errorf("\nexpected: %s\nactual: %s", expectedData, actualData)
		}
	} else {
		actualData := map[string]interface{}{}
		err := json.Unmarshal(actual, &actualData)
		if err != nil {
			t.Fatalf(`failed to read actual json: "%s", error: %s`, string(actual), err.Error())
			return
		}

		expectedData := map[string]interface{}{}
		err = json.Unmarshal([]byte(expected), &expectedData)
		if err != nil {
			t.Fatalf(`failed to read actual json: "%s", error: %s`, string(actual), err.Error())
			return
		}

		if !reflect.DeepEqual(actualData, expectedData) {
			t.Errorf("\nexpected: %s\nactual: %s", expected, actual)
		}
	}
}

func assertJsonPaths(t *testing.T, expectedPaths map[string]interface{}, actual []byte) {
	var actualData interface{}
	json.Unmarshal(actual, &actualData)
	for path, expectedValue := range expectedPaths {
		actualValue, err := jsonpath.Read(actualData, path)
		if err != nil {
			t.Errorf(`failed to find json path: "%s"`, path)
			return
		}
		if !cmp.Equal(expectedValue, actualValue)().Success() {
			t.Errorf(`values are different at path: "%s", expected: %s, actual: %s`, path, expectedValue, actualValue)
			return
		}
	}
}
