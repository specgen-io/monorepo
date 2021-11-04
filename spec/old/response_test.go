package old

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"reflect"
	"strings"
	"testing"
)

func Test_Responses_Unmarshal_Long(t *testing.T) {
	data := `
ok:
  type: empty
  description: success
bad_request:
  type: empty
  description: invalid request
`
	var responses Responses
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &responses)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(responses), 2)
	response1 := responses[0]
	response2 := responses[1]
	assert.Equal(t, response1.Name.Source, "ok")
	assert.Equal(t, reflect.DeepEqual(response1.Type.Definition, ParseType("empty")), true)
	assert.Equal(t, *response1.Description, "success")
	assert.Equal(t, response2.Name.Source, "bad_request")
	assert.Equal(t, reflect.DeepEqual(response2.Type.Definition, ParseType("empty")), true)
	assert.Equal(t, *response2.Description, "invalid request")
}

func Test_Responses_Unmarshal_Short(t *testing.T) {
	data := `
ok: empty            # success
bad_request: empty   # invalid request
`
	var responses Responses
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &responses)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(responses), 2)
	response1 := responses[0]
	response2 := responses[1]
	assert.Equal(t, response1.Name.Source, "ok")
	assert.Equal(t, reflect.DeepEqual(response1.Type.Definition, ParseType("empty")), true)
	assert.Equal(t, *response1.Description, "success")
	assert.Equal(t, response2.Name.Source, "bad_request")
	assert.Equal(t, reflect.DeepEqual(response2.Type.Definition, ParseType("empty")), true)
	assert.Equal(t, *response2.Description, "invalid request")
}

func Test_Response_WrongName_Error(t *testing.T) {
	data := `bla: empty`
	var responses Responses
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &responses)
	assert.Equal(t, err != nil, true)
	assert.Equal(t, strings.Contains(err.Error(), "bla"), true)
}

func Test_Responses_Marshal(t *testing.T) {
	expectedYaml := strings.TrimLeft(`
ok: empty # success
bad_request: empty # invalid request
`, "\n")
	var responses Responses
	checkUnmarshalMarshal(t, expectedYaml, &responses)
}
