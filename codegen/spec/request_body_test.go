package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_RequestBody_Json_Unmarshal(t *testing.T) {
	data := `Some`
	var body RequestBody
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &body)
	assert.NilError(t, err)
	assert.Equal(t, body.Type == nil, false)
	assert.Equal(t, body.Type.Definition, ParseType("Some"))
}

func Test_RequestBody_FormData_Unmarshal(t *testing.T) {
	data := `
form-data:
  int_param: int
  str_param: string
`
	var body RequestBody
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &body)
	assert.NilError(t, err)
	assert.Equal(t, body.FormData == nil, false)
}

func Test_RequestBody_FormUrlEncoded_Unmarshal(t *testing.T) {
	data := `
form-urlencoded:
  int_param: int
  str_param: string
`
	var body RequestBody
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &body)
	assert.NilError(t, err)
	assert.Equal(t, body.FormUrlEncoded == nil, false)
}

func Test_RequestBody_FormData_Marshal(t *testing.T) {
	data := strings.TrimLeft(`
form-data:
  int_param: int
  str_param: string
`, "\n")
	var body RequestBody
	checkUnmarshalMarshal(t, data, &body)
}

func Test_RequestBody_FormUrlEncoded_Marshal(t *testing.T) {
	data := strings.TrimLeft(`
form-urlencoded:
  int_param: int
  str_param: string
`, "\n")
	var body RequestBody
	checkUnmarshalMarshal(t, data, &body)
}
