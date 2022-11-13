package main

import (
	"net/http"
	"strings"
	"testing"
)

const ECHO_BODY_MODEL_JSON = `{"int_field":123,"string_field":"the value"}`

func Test_EchoBodyModel(t *testing.T) {
	req, _ := http.NewRequest("POST", serviceUrl+`/echo/body_model`, strings.NewReader(ECHO_BODY_MODEL_JSON))
	req.Header.Add("Content-Type", "application/json")

	assertJsonResponse(t, req, 200, ECHO_BODY_MODEL_JSON, nil)
}

func Test_EchoBodyModel_Request_ContentType_Charset(t *testing.T) {
	req, _ := http.NewRequest("POST", serviceUrl+`/echo/body_model`, strings.NewReader(ECHO_BODY_MODEL_JSON))
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	assertJsonResponse(t, req, 200, ECHO_BODY_MODEL_JSON, nil)
}

func Test_EchoBodyModel_Request_ContentType_Empty(t *testing.T) {
	req, _ := http.NewRequest("POST", serviceUrl+`/echo/body_model`, strings.NewReader(ECHO_BODY_MODEL_JSON))
	req.Header.Add("Content-Type", "")

	assertJsonResponse(t, req, 400, "", map[string]interface{}{
		"$.location":       "header",
		"$.message":        "Failed to parse header",
		"$.errors[0].code": "missing",
		"$.errors[0].path": "Content-Type",
	})
}

func Test_EchoBodyModel_Request_ContentType_Missing(t *testing.T) {
	req, _ := http.NewRequest("POST", serviceUrl+`/echo/body_model`, strings.NewReader(ECHO_BODY_MODEL_JSON))
	assertJsonResponse(t, req, 400, "", map[string]interface{}{
		"$.location":       "header",
		"$.message":        "Failed to parse header",
		"$.errors[0].code": "missing",
		"$.errors[0].path": "Content-Type",
	})
}

func Test_EchoBodyModel_Bad_Json(t *testing.T) {
	dataJson := `{"int_field":"the string","string_field":"the value"}` // <- int_field has string value instead of int

	req, _ := http.NewRequest("POST", serviceUrl+`/echo/body_model`, strings.NewReader(dataJson))
	req.Header.Add("Content-Type", "application/json")

	assertJsonResponse(t, req, 400, "", map[string]interface{}{
		"$.message":        "Failed to parse body",
		"$.location":       "body",
		"$.errors[0].path": "int_field",
		"$.errors[0].code": "parsing_failed",
	})
}
