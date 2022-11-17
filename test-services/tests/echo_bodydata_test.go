package main

import (
	"net/http"
	"strings"
	"testing"
)

func Test_EchoBodyArray(t *testing.T) {
	dataJson := `["the str1","the str2"]`

	req, _ := http.NewRequest("POST", serviceUrl+`/echo/body_array`, strings.NewReader(dataJson))
	req.Header.Add("Content-Type", "application/json")

	assertJsonResponse(t, req, 200, dataJson, nil)
}

func Test_EchoBodyMap(t *testing.T) {
	dataJson := `{"string_field":"the value","string_field_2":"the value_2"}`

	req, _ := http.NewRequest("POST", serviceUrl+`/echo/body_map`, strings.NewReader(dataJson))
	req.Header.Add("Content-Type", "application/json")

	assertJsonResponse(t, req, 200, dataJson, nil)
}
