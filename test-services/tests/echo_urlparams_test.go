package main

import (
	"net/http"
	"testing"
)

func Test_EchoUrlParams(t *testing.T) {
	dataJson := `
{
	"int_field":123,
	"long_field":12345,
	"float_field":1.23,
	"double_field":12.345,
	"decimal_field":12345,
	"bool_field":true,
	"string_field":"the value",
	"uuid_field":"123e4567-e89b-12d3-a456-426655440000",
	"date_field":"2020-01-01",
	"datetime_field":"2019-11-30T17:45:55",
	"enum_field":"SECOND_CHOICE"
}`

	req, _ := http.NewRequest("GET", serviceUrl+`/echo/url_params/123/12345/1.23/12.345/12345/true/the value/123e4567-e89b-12d3-a456-426655440000/2020-01-01/2019-11-30T17:45:55/SECOND_CHOICE`, nil)

	assertJsonResponse(t, req, 200, dataJson, nil)
}

func Test_EchoUrlParams_Unparsable(t *testing.T) {
	req, _ := http.NewRequest("GET", serviceUrl+`/echo/url_params/value/12345/1.23/12.345/12345/true/the value/123e4567-e89b-12d3-a456-426655440000/2020-01-01/2019-11-30T17:45:55/SECOND_CHOICE`, nil)

	assertJsonResponse(t, req, 404, `{"message":"Failed to parse url parameters"}`, nil)
}
