package main

import (
	"net/http"
	"testing"
)

func Test_EchoHeader(t *testing.T) {
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

	req, _ := http.NewRequest("GET", serviceUrl+`/echo/header`, nil)
	h := req.Header
	h.Add("Int-Header", "123")
	h.Add("Long-Header", "12345")
	h.Add("Float-Header", "1.23")
	h.Add("Double-Header", "12.345")
	h.Add("Decimal-Header", "12345")
	h.Add("Bool-Header", "true")
	h.Add("String-Header", "the value")
	h.Add("String-Opt-Header", "the value")
	h.Add("String-Defaulted-Header", "value")
	h.Add("String-Array-Header", "the str1")
	h.Add("String-Array-Header", "the str2")
	h.Add("Uuid-Header", "123e4567-e89b-12d3-a456-426655440000")
	h.Add("Date-Header", "2020-01-01")
	h.Add("Date-Array-Header", "2020-01-01")
	h.Add("Date-Array-Header", "2020-01-02")
	h.Add("Datetime-Header", "2019-11-30T17:45:55")
	h.Add("Enum-Header", "SECOND_CHOICE")

	assertJsonResponse(t, req, 200, dataJson, nil)
}

func Test_EchoHeader_Missing_Defaulted_Param(t *testing.T) {
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
	"string_defaulted_field":"the default value",
	"string_array_field":["the str1","the str2"],
	"uuid_field":"123e4567-e89b-12d3-a456-426655440000",
	"date_field":"2020-01-01",
	"date_array_field":["2020-01-01","2020-01-02"],
	"datetime_field":"2019-11-30T17:45:55",
	"enum_field":"SECOND_CHOICE"
}`

	req, _ := http.NewRequest("GET", serviceUrl+`/echo/header`, nil)
	h := req.Header
	h.Add("Int-Header", "123")
	h.Add("Long-Header", "12345")
	h.Add("Float-Header", "1.23")
	h.Add("Double-Header", "12.345")
	h.Add("Decimal-Header", "12345")
	h.Add("Bool-Header", "true")
	h.Add("String-Header", "the value")
	h.Add("String-Opt-Header", "the value")
	// h.Add("String-Defaulted-Header", "value")  //<- this is defaulted param and it's not provided
	h.Add("String-Array-Header", "the str1")
	h.Add("String-Array-Header", "the str2")
	h.Add("Uuid-Header", "123e4567-e89b-12d3-a456-426655440000")
	h.Add("Date-Header", "2020-01-01")
	h.Add("Date-Array-Header", "2020-01-01")
	h.Add("Date-Array-Header", "2020-01-02")
	h.Add("Datetime-Header", "2019-11-30T17:45:55")
	h.Add("Enum-Header", "SECOND_CHOICE")

	assertJsonResponse(t, req, 200, dataJson, nil)
}

func Test_EchoHeader_Missing_Optional_Param(t *testing.T) {
	req, _ := http.NewRequest("GET", serviceUrl+`/echo/header`, nil)
	h := req.Header
	h.Add("Int-Header", "123")
	h.Add("Long-Header", "12345")
	h.Add("Float-Header", "1.23")
	h.Add("Double-Header", "12.345")
	h.Add("Decimal-Header", "12345")
	h.Add("Bool-Header", "true")
	h.Add("String-Header", "the value")
	//h.Add("String-Opt-Header", "the value")  //<- this is defaulted param and it's not provided
	h.Add("String-Defaulted-Header", "value")
	h.Add("String-Array-Header", "the str1")
	h.Add("String-Array-Header", "the str2")
	h.Add("Uuid-Header", "123e4567-e89b-12d3-a456-426655440000")
	h.Add("Date-Header", "2020-01-01")
	h.Add("Date-Array-Header", "2020-01-01")
	h.Add("Date-Array-Header", "2020-01-02")
	h.Add("Datetime-Header", "2019-11-30T17:45:55")
	h.Add("Enum-Header", "SECOND_CHOICE")

	assertJsonResponse(t, req, 200, "", nil)
}

func Test_EchoHeader_Missing_Required_Param(t *testing.T) {
	req, _ := http.NewRequest("GET", serviceUrl+`/echo/header`, nil)
	h := req.Header
	//h.Add("Int-Header", "123")  // this param is required and it's missing
	h.Add("Long-Header", "12345")
	h.Add("Float-Header", "1.23")
	h.Add("Double-Header", "12.345")
	h.Add("Decimal-Header", "12345")
	h.Add("Bool-Header", "true")
	h.Add("String-Header", "the value")
	h.Add("String-Opt-Header", "the value")
	h.Add("String-Defaulted-Header", "value")
	h.Add("String-Array-Header", "the str1")
	h.Add("String-Array-Header", "the str2")
	h.Add("Uuid-Header", "123e4567-e89b-12d3-a456-426655440000")
	h.Add("Date-Header", "2020-01-01")
	h.Add("Date-Array-Header", "2020-01-01")
	h.Add("Date-Array-Header", "2020-01-02")
	h.Add("Datetime-Header", "2019-11-30T17:45:55")
	h.Add("Enum-Header", "SECOND_CHOICE")

	assertJsonResponse(t, req, 400, "", map[string]interface{}{
		"$.message":        "Failed to parse header",
		"$.location":       "header",
		"$.errors[0].path": "Int-Header",
		"$.errors[0].code": "missing",
	})
}
