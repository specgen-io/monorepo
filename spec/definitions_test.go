package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"reflect"
	"strings"
	"testing"
)

func Test_NamedDefinitions_Unmarshal(t *testing.T) {
	data := `
prop1: string  # some field
prop2:
  type: string
  description: another field
`
	var definitions NamedDefinitions
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &definitions)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(definitions), 2)
	prop1 := definitions[0]
	prop2 := definitions[1]

	assert.Equal(t, prop1.Name.Source, "prop1")
	assert.Equal(t, reflect.DeepEqual(prop1.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *prop1.Description, "some field")

	assert.Equal(t, prop2.Name.Source, "prop2")
	assert.Equal(t, reflect.DeepEqual(prop2.Type.Definition, ParseType("string")), true)
	assert.Equal(t, *prop2.Description, "another field")
}

func Test_NamedDefinitions_Marshal(t *testing.T) {
	expectedYaml := strings.TrimLeft(`
prop1: string # some field
prop2: int?
`, "\n")
	var definitions NamedDefinitions
	checkUnmarshalMarshal(t, expectedYaml, &definitions)
}
