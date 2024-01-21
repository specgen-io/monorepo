package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_ResponseBody_Json_Unmarshal(t *testing.T) {
	data := `MyType`
	var body ResponseBody
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &body)
	assert.NilError(t, err)
	assert.Equal(t, body.IsJson(), true)
	assert.Equal(t, body.Type == nil, false)
	assert.Equal(t, body.Type.Definition, ParseType("MyType"))
}

func Test_ResponseBody_Text_Unmarshal(t *testing.T) {
	data := `string`
	var body ResponseBody
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &body)
	assert.NilError(t, err)
	assert.Equal(t, body.IsText(), true)
}

func Test_ResponseBody_Binary_Unmarshal(t *testing.T) {
	data := `binary`
	var body ResponseBody
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &body)
	assert.NilError(t, err)
	assert.Equal(t, body.IsBinary(), true)
}

func Test_ResponseBody_Text_Marshal(t *testing.T) {
	data := strings.TrimLeft(`
string
`, "\n")
	var body ResponseBody
	checkUnmarshalMarshal(t, data, &body)
}

func Test_ResponseBody_Binary_Marshal(t *testing.T) {
	data := strings.TrimLeft(`
binary
`, "\n")
	var body ResponseBody
	checkUnmarshalMarshal(t, data, &body)
}

func Test_ResponseBody_Json_Marshal(t *testing.T) {
	data := strings.TrimLeft(`
MyType
`, "\n")
	var body ResponseBody
	checkUnmarshalMarshal(t, data, &body)
}
