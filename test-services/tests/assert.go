package main

import (
	"gotest.tools/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func contains(what string, where []string) bool {
	for _, item := range where {
		if what == item {
			return true
		}
	}
	return false
}

func assertResponseStatus(t *testing.T, req *http.Request, expectedStatusCode int) {
	resp, err := http.DefaultClient.Do(req)
	assert.NilError(t, err)

	assert.Equal(t, resp.StatusCode, expectedStatusCode)
}

func assertContentType(t *testing.T, resp *http.Response, expectedContentType string) {
	actualContentTypeValue := resp.Header.Get("Content-Type")
	actualContentType := []string{}
	for _, part := range strings.Split(actualContentTypeValue, ";") {
		actualContentType = append(actualContentType, strings.ToLower(strings.TrimSpace(part)))
	}

	if !contains(expectedContentType, actualContentType) {
		t.Errorf(`Content-Type should contain %s, but received %s'`, expectedContentType, actualContentTypeValue)
	}
}

func assertTextResponse(t *testing.T, req *http.Request, expectedStatusCode int, expectedBody string) {
	resp, err := http.DefaultClient.Do(req)
	assert.NilError(t, err)
	assert.Equal(t, resp.StatusCode, expectedStatusCode)
	assertContentType(t, resp, "text/plain")
	if expectedBody != "" {
		actualBody, err := ioutil.ReadAll(resp.Body)
		assert.NilError(t, err)
		err = resp.Body.Close()
		assert.NilError(t, err)

		assert.Equal(t, strings.TrimSuffix(string(actualBody), "\n"), expectedBody)
	}
}

func assertJsonResponse(t *testing.T, req *http.Request, expectedStatusCode int, expectedBody string, expectedBodyPaths map[string]interface{}) {
	resp, err := http.DefaultClient.Do(req)
	assert.NilError(t, err)
	assert.Equal(t, resp.StatusCode, expectedStatusCode)
	assertContentType(t, resp, "application/json")
	if expectedBody != "" || expectedBodyPaths != nil {
		actualBody, err := ioutil.ReadAll(resp.Body)
		assert.NilError(t, err)
		err = resp.Body.Close()
		assert.NilError(t, err)
		if expectedBody != "" {
			assertEqualJson(t, actualBody, expectedBody)
		}
		if expectedBodyPaths != nil {
			assertJsonPaths(t, expectedBodyPaths, actualBody)
		}
	}
}
