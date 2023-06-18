package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"reflect"
	"strings"
	"testing"
)

func Test_DefinitionDefault_Unmarshal_Short(t *testing.T) {
	data := "string = the value  # something here"
	var definition DefinitionDefault
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &definition)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(definition.Type.Definition, ParseType("string")), true)
	assert.Equal(t, definition.Default != nil, true)
	assert.Equal(t, *definition.Default, "the value")
	assert.Equal(t, definition.Description != nil, true)
	assert.Equal(t, *definition.Description, "something here")
}

func Test_Definition_Unmarshal_Short(t *testing.T) {
	data := "MyType    # some description"
	var definition Definition
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &definition)
	assert.Equal(t, err, nil)
	assert.Equal(t, reflect.DeepEqual(definition.Type.Definition, ParseType("MyType")), true)
	assert.Equal(t, *definition.Description, "some description")
}

func Test_DefinitionDefault_Marshal_Short(t *testing.T) {
	definition := &DefinitionDefault{
		Type:        Type{Definition: *Plain("string")},
		Default:     StrPtr("the value"),
		Description: StrPtr("something here"),
	}
	data, err := yaml.Marshal(definition)
	expected := "string = the value # something here"
	assert.Equal(t, err, nil)
	actual := strings.TrimRight(string(data), "\n")
	assert.Equal(t, actual, expected)
}

func Test_Definition_Marshal_Short(t *testing.T) {
	expectedYaml := "string # something here\n"
	var definition Definition
	checkUnmarshalMarshal(t, expectedYaml, &definition)
}
