package main

import (
	"net/http"
	"strings"
	"testing"
)

var serviceUrl = "http://localhost:8081"

func Test_EchoEverything(t *testing.T) {
	dataJson := `
{
	"body_field":{"int_field":123,"string_field":"the value"},
	"float_query":1.23,
	"bool_query":true,
	"uuid_header":"123e4567-e89b-12d3-a456-426655440000",
	"datetime_header":"2019-11-30T17:45:55",
	"date_url":"2020-01-01",
	"decimal_url":12345
}`

	req, _ := http.NewRequest("POST", serviceUrl+`/echo/everything/2020-01-01/12345`, strings.NewReader(`{"int_field":123,"string_field":"the value"}`))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Uuid-Header", "123e4567-e89b-12d3-a456-426655440000")
	req.Header.Add("Datetime-Header", "2019-11-30T17:45:55")
	q := req.URL.Query()
	q.Add("float_query", "1.23")
	q.Add("bool_query", "true")
	req.URL.RawQuery = q.Encode()

	assertJsonResponse(t, req, 200, dataJson, nil)
}

func Test_CheckEmpty(t *testing.T) {
	req, _ := http.NewRequest("GET", serviceUrl+`/check/empty`, nil)
	assertResponseStatus(t, req, 200)
}

func Test_CheckEmpty_Response(t *testing.T) {
	dataJson := `{"int_field":123,"string_field":"the value"}`

	req, _ := http.NewRequest("POST", serviceUrl+`/check/empty_response`, strings.NewReader(dataJson))
	req.Header.Add("Content-Type", "application/json")

	assertResponseStatus(t, req, 200)
}

func Test_CheckForbidden(t *testing.T) {
	req, _ := http.NewRequest("GET", serviceUrl+`/check/forbidden`, nil)
	assertResponseStatus(t, req, 403)
}

func Test_V2_EchoBodyModel(t *testing.T) {
	dataJson := `{"bool_field":true,"string_field":"the value"}`

	req, _ := http.NewRequest("POST", serviceUrl+`/v2/echo/body_model`, strings.NewReader(dataJson))
	req.Header.Add("Content-Type", "application/json")

	assertJsonResponse(t, req, 200, dataJson, nil)
}
