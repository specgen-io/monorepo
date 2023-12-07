package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"
)

func Test_EchoFormData(t *testing.T) {
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

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("int_param", "123")
	_ = writer.WriteField("long_param", "12345")
	_ = writer.WriteField("float_param", "1.23")
	_ = writer.WriteField("double_param", "12.345")
	_ = writer.WriteField("decimal_param", "12345")
	_ = writer.WriteField("bool_param", "true")
	_ = writer.WriteField("string_param", "the value")
	_ = writer.WriteField("string_opt_param", "the value")
	_ = writer.WriteField("string_defaulted_param", "value")
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		_ = writer.WriteField("string_array_param", "the str1,the str2")
	} else {
		_ = writer.WriteField("string_array_param", "the str1")
		_ = writer.WriteField("string_array_param", "the str2")
	}
	_ = writer.WriteField("uuid_param", "123e4567-e89b-12d3-a456-426655440000")
	_ = writer.WriteField("date_param", "2020-01-01")
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		_ = writer.WriteField("date_array_param", "2020-01-01,2020-01-02")
	} else {
		_ = writer.WriteField("date_array_param", "2020-01-01")
		_ = writer.WriteField("date_array_param", "2020-01-02")
	}
	_ = writer.WriteField("datetime_param", "2019-11-30T17:45:55")
	_ = writer.WriteField("enum_param", "SECOND_CHOICE")
	_ = writer.Close()

	req, _ := http.NewRequest("POST", serviceUrl+`/echo/form_data`, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	assertJsonResponse(t, req, 200, dataJson, nil)
}

func Test_EchoFormData_Missing_Required_Param(t *testing.T) {
	skipIf(t, NO_FORM_DATA)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("int_param", "123")
	_ = writer.WriteField("long_param", "12345")
	_ = writer.WriteField("float_param", "1.23")
	_ = writer.WriteField("double_param", "12.345")
	_ = writer.WriteField("decimal_param", "12345")
	_ = writer.WriteField("bool_param", "true")
	//_ = writer.WriteField("string_param", "the value") //<- this is required param and it's not provided
	_ = writer.WriteField("string_opt_param", "the value")
	_ = writer.WriteField("string_defaulted_param", "value")
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		_ = writer.WriteField("string_array_param", "the str1,the str2")
	} else {
		_ = writer.WriteField("string_array_param", "the str1")
		_ = writer.WriteField("string_array_param", "the str2")
	}
	_ = writer.WriteField("uuid_param", "123e4567-e89b-12d3-a456-426655440000")
	_ = writer.WriteField("date_param", "2020-01-01")
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		_ = writer.WriteField("date_array_param", "2020-01-01,2020-01-02")
	} else {
		_ = writer.WriteField("date_array_param", "2020-01-01")
		_ = writer.WriteField("date_array_param", "2020-01-02")
	}
	_ = writer.WriteField("datetime_param", "2019-11-30T17:45:55")
	_ = writer.WriteField("enum_param", "SECOND_CHOICE")
	_ = writer.Close()

	req, _ := http.NewRequest("POST", serviceUrl+`/echo/form_data`, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if check(PARAMETERS_MODE) {
		assertJsonResponse(t, req, 400, "", map[string]interface{}{
			"$.message":        "Failed to parse parameters",
			"$.location":       "parameters",
			"$.errors[0].path": "string_param",
			"$.errors[0].code": "missing",
		})
	} else {
		assertJsonResponse(t, req, 400, "", map[string]interface{}{
			"$.message":        "Failed to parse body",
			"$.location":       "body",
			"$.errors[0].path": "string_param",
			"$.errors[0].code": "missing",
		})
	}
}

func Test_EchoFormData_Missing_Optional_Param(t *testing.T) {
	skipIf(t, NO_FORM_DATA)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("int_param", "123")
	_ = writer.WriteField("long_param", "12345")
	_ = writer.WriteField("float_param", "1.23")
	_ = writer.WriteField("double_param", "12.345")
	_ = writer.WriteField("decimal_param", "12345")
	_ = writer.WriteField("bool_param", "true")
	_ = writer.WriteField("string_param", "the value")
	//_ = writer.WriteField("string_opt_param", "the value")	//<- this is optional param and it's not provided
	_ = writer.WriteField("string_defaulted_param", "value")
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		_ = writer.WriteField("string_array_param", "the str1,the str2")
	} else {
		_ = writer.WriteField("string_array_param", "the str1")
		_ = writer.WriteField("string_array_param", "the str2")
	}
	_ = writer.WriteField("uuid_param", "123e4567-e89b-12d3-a456-426655440000")
	_ = writer.WriteField("date_param", "2020-01-01")
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		_ = writer.WriteField("date_array_param", "2020-01-01,2020-01-02")
	} else {
		_ = writer.WriteField("date_array_param", "2020-01-01")
		_ = writer.WriteField("date_array_param", "2020-01-02")
	}
	_ = writer.WriteField("datetime_param", "2019-11-30T17:45:55")
	_ = writer.WriteField("enum_param", "SECOND_CHOICE")
	_ = writer.Close()

	req, _ := http.NewRequest("POST", serviceUrl+`/echo/form_data`, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	assertResponseStatus(t, req, 200)
}

func Test_EchoFormData_Missing_Defaulted_Param(t *testing.T) {
	skipIf(t, NO_FORM_DATA)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("int_param", "123")
	_ = writer.WriteField("long_param", "12345")
	_ = writer.WriteField("float_param", "1.23")
	_ = writer.WriteField("double_param", "12.345")
	_ = writer.WriteField("decimal_param", "12345")
	_ = writer.WriteField("bool_param", "true")
	_ = writer.WriteField("string_param", "the value")
	_ = writer.WriteField("string_opt_param", "the value")
	//_ = writer.WriteField("string_defaulted_param", "value")	//<- this is defaulted param and it's not provided
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		_ = writer.WriteField("string_array_param", "the str1,the str2")
	} else {
		_ = writer.WriteField("string_array_param", "the str1")
		_ = writer.WriteField("string_array_param", "the str2")
	}
	_ = writer.WriteField("uuid_param", "123e4567-e89b-12d3-a456-426655440000")
	_ = writer.WriteField("date_param", "2020-01-01")
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		_ = writer.WriteField("date_array_param", "2020-01-01,2020-01-02")
	} else {
		_ = writer.WriteField("date_array_param", "2020-01-01")
		_ = writer.WriteField("date_array_param", "2020-01-02")
	}
	_ = writer.WriteField("datetime_param", "2019-11-30T17:45:55")
	_ = writer.WriteField("enum_param", "SECOND_CHOICE")
	_ = writer.Close()

	req, _ := http.NewRequest("POST", serviceUrl+`/echo/form_data`, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	assertJsonResponse(t, req, 200, "", map[string]interface{}{
		"$.string_defaulted_field": "the default value",
	})
}

func Test_EchoFormData_WrongFormat(t *testing.T) {
	skipIf(t, NO_FORM_DATA)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("int_param", "abc")
	_ = writer.WriteField("long_param", "12345")
	_ = writer.WriteField("float_param", "1.23")
	_ = writer.WriteField("double_param", "12.345")
	_ = writer.WriteField("decimal_param", "12345")
	_ = writer.WriteField("bool_param", "true")
	_ = writer.WriteField("string_param", "the value")
	_ = writer.WriteField("string_opt_param", "the value")
	_ = writer.WriteField("string_defaulted_param", "value")
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		_ = writer.WriteField("string_array_param", "the str1,the str2")
	} else {
		_ = writer.WriteField("string_array_param", "the str1")
		_ = writer.WriteField("string_array_param", "the str2")
	}
	_ = writer.WriteField("uuid_param", "123e4567-e89b-12d3-a456-426655440000")
	_ = writer.WriteField("date_param", "2020-01-01")
	if check(COMMA_SEPARATED_FORM_PARAMS_MODE) {
		_ = writer.WriteField("date_array_param", "2020-01-01,2020-01-02")
	} else {
		_ = writer.WriteField("date_array_param", "2020-01-01")
		_ = writer.WriteField("date_array_param", "2020-01-02")
	}
	_ = writer.WriteField("datetime_param", "2019-11-30T17:45:55")
	_ = writer.WriteField("enum_param", "SECOND_CHOICE")
	_ = writer.Close()

	req, _ := http.NewRequest("POST", serviceUrl+`/echo/form_data`, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

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
