package yamlx

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func TestYamlMap(t *testing.T) {
	theMap := Map()
	theMap.Add("key1", "one")
	theMap.Add("key2", "two")
	yamlData, _ := yaml.Marshal(theMap)
	expectedYaml := `
key1: one
key2: two
`
	assert.Equal(t, strings.TrimSpace(expectedYaml), strings.TrimSpace(string(yamlData)))
}

func TestYamlMapNestedValues(t *testing.T) {
	array := Array()
	array.Add("one", "two")
	theMap := Map()
	theMap.Add("string", "the string")
	theMap.Add("array", array)
	yamlData, _ := yaml.Marshal(theMap)
	expectedYaml := `
string: the string
array:
    - one
    - two
`
	assert.Equal(t, strings.TrimSpace(expectedYaml), strings.TrimSpace(string(yamlData)))
}
