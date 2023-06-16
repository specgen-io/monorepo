package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"testing"
)

func Test_Request_body_Unmarshal_Json(t *testing.T) {
	data := `Some`
	var body RequestBody
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &body)
	assert.NilError(t, err)
	assert.Equal(t, body.Type == nil, false)
	assert.Equal(t, body.Type.Definition, ParseType("Some"))
}

func Test_Request_body_Unmarshal_FormData(t *testing.T) {
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

func Test_Request_body_Unmarshal_FormUrlEncoded(t *testing.T) {
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
