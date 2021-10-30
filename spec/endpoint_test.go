package spec

import (
	"gotest.tools/assert"
	"reflect"
	"strings"
	"testing"
)

func Test_ParseEndpoint_NoParams(t *testing.T) {
	endpoint, err := parseEndpoint("GET /some/url", nil)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, endpoint.Method, "GET")
	assert.Equal(t, endpoint.Url, "/some/url")
	assert.Equal(t, reflect.DeepEqual(endpoint.UrlParams, UrlParams{}), true)
}

func Test_ParseEndpoint_Param(t *testing.T) {
	endpoint, err := parseEndpoint("POST /some/url/{id:str}", nil)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, endpoint.Method, "POST")
	assert.Equal(t, endpoint.Url, "/some/url/{id}")
	assert.Equal(t, len(endpoint.UrlParams), 1)
	idParam := endpoint.UrlParams[0]
	assert.Equal(t, idParam.Name.Source, "id")
	assert.Equal(t, idParam.Type.Definition, ParseType("str"))
}

func Test_ParseEndpoint_MultipleParams(t *testing.T) {
	endpoint, err := parseEndpoint("GET /some/url/{some_id:str}/{the_name:str}", nil)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, endpoint.Method, "GET")
	assert.Equal(t, endpoint.Url, "/some/url/{some_id}/{the_name}")
	assert.Equal(t, len(endpoint.UrlParams), 2)
	idParam := endpoint.UrlParams[0]
	nameParam := endpoint.UrlParams[1]
	assert.Equal(t, idParam.Name.Source, "some_id")
	assert.Equal(t, idParam.Type.Definition, ParseType("str"))
	assert.Equal(t, nameParam.Name.Source, "the_name")
	assert.Equal(t, nameParam.Type.Definition, ParseType("str"))
}

func Test_ParseEndpoint_EnumType(t *testing.T) {
	endpoint, err := parseEndpoint("GET /some/url/{some_id:CustomEnum}", nil)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, endpoint.Method, "GET")
	assert.Equal(t, endpoint.Url, "/some/url/{some_id}")
	assert.Equal(t, len(endpoint.UrlParams), 1)
	idParam := endpoint.UrlParams[0]
	assert.Equal(t, idParam.Name.Source, "some_id")
	assert.Equal(t, idParam.Type.Definition, ParseType("CustomEnum"))
}

func Test_ParseEndpoint_ToManyParts(t *testing.T) {
	_, err := parseEndpoint("GET /some/ url/", nil)
	assert.Equal(t, err != nil, true)
	assert.Equal(t, strings.Contains(err.Error(), "endpoint"), true)
}

func Test_ParseEndpoint_WrongMethod(t *testing.T) {
	_, err := parseEndpoint("METHOD /some/url/", nil)
	assert.Equal(t, err != nil, true)
	assert.Equal(t, strings.Contains(err.Error(), "METHOD"), true)
}

func Test_Endpoint_EnumParam_Marshal(t *testing.T) {
	expectedYaml := "GET /some/url/{some_id:CustomEnum}\n"
	var endpoint Endpoint
	checkUnmarshalMarshal(t, expectedYaml, &endpoint)
}
