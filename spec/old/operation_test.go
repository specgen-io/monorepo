package old

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
	assert.Equal(t, operation.Body == nil, true)
	assert.Equal(t, len(operation.Responses), 1)
	response := operation.Responses[0]
	assert.Equal(t, response.Name.Source, "ok")
	assert.Equal(t, response.Type.Definition, ParseType("empty"))
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

func Test_Operations_Unmarshal(t *testing.T) {
	data := `
some_url:
  endpoint: GET /some/url
  response:
    ok: empty
ping:
  endpoint: GET /ping
  response:
    ok: empty
`

	var operations Operations
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &operations)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(operations), 2)
	operation1 := operations[0]
	operation2 := operations[1]
	assert.Equal(t, operation1.Name.Source, "some_url")
	assert.Equal(t, operation2.Name.Source, "ping")
}

func Test_Operations_Unmarshal_Description(t *testing.T) {
	data := `
some_url:     # some url description
  endpoint: GET /some/url
  response:
    ok: empty
ping:         # ping description
  endpoint: GET /ping
  response:
    ok: empty
`

	var operations Operations
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &operations)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(operations), 2)
	operation1 := operations[0]
	operation2 := operations[1]

	assert.Equal(t, operation1.Name.Source, "some_url")
	assert.Equal(t, *operation1.Description, "some url description")
	assert.Equal(t, operation2.Name.Source, "ping")
	assert.Equal(t, *operation2.Description, "ping description")
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
