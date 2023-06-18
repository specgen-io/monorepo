package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

func Test_ResponseBody_Unmarshal(t *testing.T) {
	data := "MyType    # some description"
	var definition Definition
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &definition)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(definition.Type.Definition, ParseType("MyType")), true)
	assert.Equal(t, *definition.Description, "some description")
}

func Test_ResponseBody_String_Unmarshal(t *testing.T) {
	data := "string    # some description"
	var definition Definition
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &definition)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(definition.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *definition.Description, "some description")
}

func Test_ResponseBody_Empty_Unmarshal(t *testing.T) {
	data := "empty    # some description"
	var definition Definition
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &definition)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(definition.Type.Definition, ParseType("empty")), true)
	assert.Equal(t, *definition.Description, "some description")
}

func Test_ResponseBody_Marshal(t *testing.T) {
	expectedYaml := "MyType # something here\n"
	var definition Definition
	checkUnmarshalMarshal(t, expectedYaml, &definition)
}

func Test_ResponseBody_String_Marshal(t *testing.T) {
	expectedYaml := "string # something here\n"
	var definition Definition
	checkUnmarshalMarshal(t, expectedYaml, &definition)
}

func Test_ResponseBody_Empty_Marshal(t *testing.T) {
	expectedYaml := "empty # something here\n"
	var definition Definition
	checkUnmarshalMarshal(t, expectedYaml, &definition)
}
