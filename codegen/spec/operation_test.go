package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_Operation_Unmarshal(t *testing.T) {
	data := `
endpoint: GET /some/url
response:
  ok: empty
`

	var operation Operation
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &operation)
	assert.Equal(t, err, nil)

	assert.Equal(t, operation.Endpoint.Method, "GET")
	assert.Equal(t, operation.Endpoint.Url, "/some/url")
	assert.Assert(t, operation.Body.IsEmpty())
	assert.Equal(t, len(operation.Responses), 1)
	response := operation.Responses[0]
	assert.Equal(t, response.Name.Source, "ok")
	assert.Assert(t, response.Body.IsEmpty())
}

func Test_Operation_Unmarshal_QueryParams(t *testing.T) {
	data := `
endpoint: GET /ping
query:
  message: string?
response:
  ok: empty
`

	var operation Operation
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &operation)
	assert.Equal(t, err, nil)

	assert.Equal(t, operation.Endpoint.Method, "GET")
	assert.Equal(t, operation.Endpoint.Url, "/ping")
	assert.Equal(t, len(operation.QueryParams), 1)
	queryParam := operation.QueryParams[0]
	assert.Equal(t, queryParam.Name.Source, "message")
	assert.Equal(t, queryParam.Type.Definition.Name, "string?")
}

func Test_Operation_Unmarshal_BodyDescription(t *testing.T) {
	data := `
endpoint: GET /some/url
body: Some  # body description
response:
  ok: empty
`

	var operation Operation
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &operation)
	assert.Equal(t, err, nil)

	assert.Equal(t, *operation.Body.Description, "body description")
}

func Test_Operation_Marshal(t *testing.T) {
	expectedYaml := strings.TrimLeft(`
endpoint: GET /some/url
header:
  Message: string?
query:
  message: string?
body: Some # body description
response:
  ok: empty
`, "\n")
	var operation Operation
	checkUnmarshalMarshal(t, expectedYaml, &operation)
}

func Test_Operations_Marshal(t *testing.T) {
	expectedYaml := strings.TrimLeft(`
some_url:
  endpoint: GET /some/url
  description: some url description
  response:
    ok: empty
ping:
  endpoint: GET /ping
  description: ping description
  response:
    ok: empty
`, "\n")
	var operations Operations
	checkUnmarshalMarshal(t, expectedYaml, &operations)
}
