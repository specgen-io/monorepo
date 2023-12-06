package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_EchoFormUrlencoded(t *testing.T) {
	skipIf(t, NO_FORM_DATA)
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

	var body io.Reader
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		body = strings.NewReader("int_param=123&long_param=12345&float_param=1.23&double_param=12.345&decimal_param=12345&bool_param=true&string_param=the value&string_opt_param=the value&string_defaulted_param=value&string_array_param=the str1,the str2&uuid_param=123e4567-e89b-12d3-a456-426655440000&date_param=2020-01-01&date_array_param=2020-01-01,2020-01-02&datetime_param=2019-11-30T17:45:55&enum_param=SECOND_CHOICE")
	} else {
		body = strings.NewReader("int_param=123&long_param=12345&float_param=1.23&double_param=12.345&decimal_param=12345&bool_param=true&string_param=the value&string_opt_param=the value&string_defaulted_param=value&string_array_param=the str1&string_array_param=the str2&uuid_param=123e4567-e89b-12d3-a456-426655440000&date_param=2020-01-01&date_array_param=2020-01-01&date_array_param=2020-01-02&datetime_param=2019-11-30T17:45:55&enum_param=SECOND_CHOICE")
	}
	req, _ := http.NewRequest("POST", serviceUrl+`/echo/form_urlencoded`, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	assertJsonResponse(t, req, 200, dataJson, nil)
}

func Test_EchoFormUrlencoded_Missing_Required_Param(t *testing.T) {
	skipIf(t, NO_FORM_DATA)
	var body io.Reader
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		body = strings.NewReader("int_param=value&long_param=12345&float_param=1.23&double_param=12.345&decimal_param=12345&bool_param=true&string_opt_param=the value&string_defaulted_param=value&string_array_param=the str1,the str2&uuid_param=123e4567-e89b-12d3-a456-426655440000&date_param=2020-01-01&date_array_param=2020-01-01,2020-01-02&datetime_param=2019-11-30T17:45:55&enum_param=SECOND_CHOICE")
	} else {
		body = strings.NewReader("int_param=value&long_param=12345&float_param=1.23&double_param=12.345&decimal_param=12345&bool_param=true&string_param=the value&string_opt_param=the value&string_defaulted_param=value&string_array_param=the str1&string_array_param=the str2&uuid_param=123e4567-e89b-12d3-a456-426655440000&date_param=2020-01-01&date_array_param=2020-01-01&date_array_param=2020-01-02&datetime_param=2019-11-30T17:45:55&enum_param=SECOND_CHOICE")
	}
	req, _ := http.NewRequest("POST", serviceUrl+`/echo/form_urlencoded`, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if check(PARAMETERS_MODE) {
		assertJsonResponse(t, req, 400, "", map[string]interface{}{
			"$.message":        "Failed to parse parameters",
			"$.location":       "parameters",
			"$.errors[0].path": "int_param",
			"$.errors[0].code": "parsing_failed",
		})
	} else {
		assertJsonResponse(t, req, 400, "", map[string]interface{}{
			"$.message":        "Failed to parse body",
			"$.location":       "body",
			"$.errors[0].path": "int_param",
			"$.errors[0].code": "parsing_failed",
		})
	}
}

func Test_EchoFormUrlencoded_Missing_Optional_Param(t *testing.T) {
	skipIf(t, NO_FORM_DATA)
	var body io.Reader
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		body = strings.NewReader("int_param=value&long_param=12345&float_param=1.23&double_param=12.345&decimal_param=12345&bool_param=true&string_param=the value&string_defaulted_param=value&string_array_param=the str1,the str2&uuid_param=123e4567-e89b-12d3-a456-426655440000&date_param=2020-01-01&date_array_param=2020-01-01,2020-01-02&datetime_param=2019-11-30T17:45:55&enum_param=SECOND_CHOICE")
	} else {
		body = strings.NewReader("int_param=value&long_param=12345&float_param=1.23&double_param=12.345&decimal_param=12345&bool_param=true&string_param=the value&string_opt_param=the value&string_defaulted_param=value&string_array_param=the str1&string_array_param=the str2&uuid_param=123e4567-e89b-12d3-a456-426655440000&date_param=2020-01-01&date_array_param=2020-01-01&date_array_param=2020-01-02&datetime_param=2019-11-30T17:45:55&enum_param=SECOND_CHOICE")
	}
	req, _ := http.NewRequest("POST", serviceUrl+`/echo/form_urlencoded`, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if check(PARAMETERS_MODE) {
		assertJsonResponse(t, req, 400, "", map[string]interface{}{
			"$.message":        "Failed to parse parameters",
			"$.location":       "parameters",
			"$.errors[0].path": "int_param",
			"$.errors[0].code": "parsing_failed",
		})
	} else {
		assertJsonResponse(t, req, 400, "", map[string]interface{}{
			"$.message":        "Failed to parse body",
			"$.location":       "body",
			"$.errors[0].path": "int_param",
			"$.errors[0].code": "parsing_failed",
		})
	}
}

func Test_EchoFormUrlencoded_Missing_Defaulted_Param(t *testing.T) {
	skipIf(t, NO_FORM_DATA)
	var body io.Reader
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		body = strings.NewReader("int_param=value&long_param=12345&float_param=1.23&double_param=12.345&decimal_param=12345&bool_param=true&string_param=the value&string_opt_param=the value&string_array_param=the str1,the str2&uuid_param=123e4567-e89b-12d3-a456-426655440000&date_param=2020-01-01&date_array_param=2020-01-01,2020-01-02&datetime_param=2019-11-30T17:45:55&enum_param=SECOND_CHOICE")
	} else {
		body = strings.NewReader("int_param=value&long_param=12345&float_param=1.23&double_param=12.345&decimal_param=12345&bool_param=true&string_param=the value&string_opt_param=the value&string_defaulted_param=value&string_array_param=the str1&string_array_param=the str2&uuid_param=123e4567-e89b-12d3-a456-426655440000&date_param=2020-01-01&date_array_param=2020-01-01&date_array_param=2020-01-02&datetime_param=2019-11-30T17:45:55&enum_param=SECOND_CHOICE")
	}
	req, _ := http.NewRequest("POST", serviceUrl+`/echo/form_urlencoded`, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if check(PARAMETERS_MODE) {
		assertJsonResponse(t, req, 400, "", map[string]interface{}{
			"$.message":        "Failed to parse parameters",
			"$.location":       "parameters",
			"$.errors[0].path": "int_param",
			"$.errors[0].code": "parsing_failed",
		})
	} else {
		assertJsonResponse(t, req, 400, "", map[string]interface{}{
			"$.message":        "Failed to parse body",
			"$.location":       "body",
			"$.errors[0].path": "int_param",
			"$.errors[0].code": "parsing_failed",
		})
	}
}

func Test_EchoFormUrlencoded_WrongFormat(t *testing.T) {
	skipIf(t, NO_FORM_DATA)
	var body io.Reader
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		body = strings.NewReader("int_param=value&long_param=12345&float_param=1.23&double_param=12.345&decimal_param=12345&bool_param=true&string_param=the value&string_opt_param=the value&string_defaulted_param=value&string_array_param=the str1,the str2&uuid_param=123e4567-e89b-12d3-a456-426655440000&date_param=2020-01-01&date_array_param=2020-01-01,2020-01-02&datetime_param=2019-11-30T17:45:55&enum_param=SECOND_CHOICE")
	} else {
		body = strings.NewReader("int_param=value&long_param=12345&float_param=1.23&double_param=12.345&decimal_param=12345&bool_param=true&string_param=the value&string_opt_param=the value&string_defaulted_param=value&string_array_param=the str1&string_array_param=the str2&uuid_param=123e4567-e89b-12d3-a456-426655440000&date_param=2020-01-01&date_array_param=2020-01-01&date_array_param=2020-01-02&datetime_param=2019-11-30T17:45:55&enum_param=SECOND_CHOICE")
	}
	req, _ := http.NewRequest("POST", serviceUrl+`/echo/form_urlencoded`, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if check(PARAMETERS_MODE) {
		assertJsonResponse(t, req, 400, "", map[string]interface{}{
			"$.message":        "Failed to parse parameters",
			"$.location":       "parameters",
			"$.errors[0].path": "int_param",
			"$.errors[0].code": "parsing_failed",
		})
	} else {
		assertJsonResponse(t, req, 400, "", map[string]interface{}{
			"$.message":        "Failed to parse body",
			"$.location":       "body",
			"$.errors[0].path": "int_param",
			"$.errors[0].code": "parsing_failed",
		})
	}
}
