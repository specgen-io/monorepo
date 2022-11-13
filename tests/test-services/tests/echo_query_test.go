package main

import (
	"net/http"
	"testing"
)

func Test_EchoQuery(t *testing.T) {
	dataJson := `
{
	"int_field":123,
	"long_field":12345,
	"float_field":1.23,
	"double_field":12.345,
	"decimal_field":12345,
	"bool_field":true,
	"string_field":"the value",
	"string_opt_field":"the value",
	"string_defaulted_field":"value",
	"string_array_field":["the str1","the str2"],
	"uuid_field":"123e4567-e89b-12d3-a456-426655440000",
	"date_field":"2020-01-01",
	"date_array_field":["2020-01-01","2020-01-02"],
	"datetime_field":"2019-11-30T17:45:55",
	"enum_field":"SECOND_CHOICE"
}`

	req, _ := http.NewRequest("GET", serviceUrl+`/echo/query`, nil)
	q := req.URL.Query()
	q.Add("int_query", "123")
	q.Add("long_query", "12345")
	q.Add("float_query", "1.23")
	q.Add("double_query", "12.345")
	q.Add("decimal_query", "12345")
	q.Add("bool_query", "true")
	q.Add("string_query", "the value")
	q.Add("string_opt_query", "the value")
	q.Add("string_defaulted_query", "value")
	q.Add("string_array_query", "the str1")
	q.Add("string_array_query", "the str2")
	q.Add("uuid_query", "123e4567-e89b-12d3-a456-426655440000")
	q.Add("date_query", "2020-01-01")
	q.Add("date_array_query", "2020-01-01")
	q.Add("date_array_query", "2020-01-02")
	q.Add("datetime_query", "2019-11-30T17:45:55")
	q.Add("enum_query", "SECOND_CHOICE")
	req.URL.RawQuery = q.Encode()

	assertJsonResponse(t, req, 200, dataJson, nil)
}

func Test_EchoQuery_Missing_Required_Param(t *testing.T) {
	req, _ := http.NewRequest("GET", serviceUrl+`/echo/query`, nil)
	q := req.URL.Query()
	q.Add("int_query", "123")
	q.Add("long_query", "12345")
	q.Add("float_query", "1.23")
	q.Add("double_query", "12.345")
	q.Add("decimal_query", "12345")
	q.Add("bool_query", "true")
	//q.Add("string_query", "the value")  //<- this is required param and it's not provided
	q.Add("string_opt_query", "the value")
	q.Add("string_defaulted_query", "value")
	q.Add("string_array_query", "the str1")
	q.Add("string_array_query", "the str2")
	q.Add("uuid_query", "123e4567-e89b-12d3-a456-426655440000")
	q.Add("date_query", "2020-01-01")
	q.Add("date_array_query", "2020-01-01")
	q.Add("date_array_query", "2020-01-02")
	q.Add("datetime_query", "2019-11-30T17:45:55")
	q.Add("enum_query", "SECOND_CHOICE")
	req.URL.RawQuery = q.Encode()

	assertJsonResponse(t, req, 400, "", map[string]interface{}{
		"$.message":        "Failed to parse query",
		"$.location":       "query",
		"$.errors[0].path": "string_query",
		"$.errors[0].code": "missing",
	})
}

func Test_EchoQuery_Missing_Optional_Param(t *testing.T) {
	req, _ := http.NewRequest("GET", serviceUrl+`/echo/query`, nil)
	q := req.URL.Query()
	q.Add("int_query", "123")
	q.Add("long_query", "12345")
	q.Add("float_query", "1.23")
	q.Add("double_query", "12.345")
	q.Add("decimal_query", "12345")
	q.Add("bool_query", "true")
	q.Add("string_query", "the value")
	//q.Add("string_opt_query", "the value")  //<- this is optional param and it's not provided
	q.Add("string_defaulted_query", "value")
	q.Add("string_array_query", "the str1")
	q.Add("string_array_query", "the str2")
	q.Add("uuid_query", "123e4567-e89b-12d3-a456-426655440000")
	q.Add("date_query", "2020-01-01")
	q.Add("date_array_query", "2020-01-01")
	q.Add("date_array_query", "2020-01-02")
	q.Add("datetime_query", "2019-11-30T17:45:55")
	q.Add("enum_query", "SECOND_CHOICE")
	req.URL.RawQuery = q.Encode()

	assertResponseStatus(t, req, 200)
}

func Test_EchoQuery_Missing_Defaulted_Param(t *testing.T) {
	req, _ := http.NewRequest("GET", serviceUrl+`/echo/query`, nil)
	q := req.URL.Query()
	q.Add("int_query", "123")
	q.Add("long_query", "12345")
	q.Add("float_query", "1.23")
	q.Add("double_query", "12.345")
	q.Add("decimal_query", "12345")
	q.Add("bool_query", "true")
	q.Add("string_query", "the value")
	q.Add("string_opt_query", "the value")
	//q.Add("string_defaulted_query", "value") //<- this is defaulted param and it's not provided
	q.Add("string_array_query", "the str1")
	q.Add("string_array_query", "the str2")
	q.Add("uuid_query", "123e4567-e89b-12d3-a456-426655440000")
	q.Add("date_query", "2020-01-01")
	q.Add("date_array_query", "2020-01-01")
	q.Add("date_array_query", "2020-01-02")
	q.Add("datetime_query", "2019-11-30T17:45:55")
	q.Add("enum_query", "SECOND_CHOICE")
	req.URL.RawQuery = q.Encode()

	assertJsonResponse(t, req, 200, "", map[string]interface{}{
		"$.string_defaulted_field": "the default value",
	})
}

func Test_EchoQuery_WrongFormat(t *testing.T) {
	req, _ := http.NewRequest("GET", serviceUrl+`/echo/query`, nil)
	q := req.URL.Query()
	q.Add("int_query", "abc")
	q.Add("long_query", "12345")
	q.Add("float_query", "1.23")
	q.Add("double_query", "12.345")
	q.Add("decimal_query", "12345")
	q.Add("bool_query", "true")
	q.Add("string_query", "the value")
	q.Add("string_opt_query", "the value")
	q.Add("string_defaulted_query", "value")
	q.Add("string_array_query", "the str1")
	q.Add("string_array_query", "the str2")
	q.Add("uuid_query", "123e4567-e89b-12d3-a456-426655440000")
	q.Add("date_query", "2020-01-01")
	q.Add("date_array_query", "2020-01-01")
	q.Add("date_array_query", "2020-01-02")
	q.Add("datetime_query", "2019-11-30T17:45:55")
	q.Add("enum_query", "SECOND_CHOICE")
	req.URL.RawQuery = q.Encode()

	assertJsonResponse(t, req, 400, "", map[string]interface{}{
		"$.message":        "Failed to parse query",
		"$.location":       "query",
		"$.errors[0].path": "int_query",
		"$.errors[0].code": "parsing_failed",
	})
}
