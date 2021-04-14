package genopenapi

import (
	"github.com/vsapronov/yaml"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func TestYamlArray(t *testing.T) {
	array := Array().Add("one").Add("two")
	yamlData, _ := yaml.Marshal(array)
	expectedYaml := `
- one
- two
`
	assert.Equal(t, strings.TrimSpace(expectedYaml), strings.TrimSpace(string(yamlData)))
}

func TestYamlMap(t *testing.T) {
	theMap := Map().Set("key1", "one").Set("key2", "two")
	yamlData, _ := yaml.Marshal(theMap)
	expectedYaml := `
key1: one
key2: two
`
	assert.Equal(t, strings.TrimSpace(expectedYaml), strings.TrimSpace(string(yamlData)))
}

func TestYamlMapNestedValues(t *testing.T) {
	array := Array().Add("one").Add("two")
	theMap := Map().Set("string", "the string").Set("array", array)
	yamlData, _ := yaml.Marshal(theMap)
	expectedYaml := `
string: the string
array:
    - one
    - two
`
	assert.Equal(t, strings.TrimSpace(expectedYaml), strings.TrimSpace(string(yamlData)))
}