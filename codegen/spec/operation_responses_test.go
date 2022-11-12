package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_Response_WrongName_Error(t *testing.T) {
	data := `bla: empty`
	var responses OperationResponses
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &responses)
	assert.Equal(t, err != nil, true)
	assert.Equal(t, strings.Contains(err.Error(), "bla"), true)
}

func Test_Responses_Marshal(t *testing.T) {
	expectedYaml := strings.TrimLeft(`
ok: empty # success
bad_request: empty # invalid request
`, "\n")
	var responses OperationResponses
	checkUnmarshalMarshal(t, expectedYaml, &responses)
}
