package yamlx

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func TestYamlArray(t *testing.T) {
	array := Array()
	array.Add("one", "two")
	yamlData, _ := yaml.Marshal(array)
	expectedYaml := `
- one
- two
`
	assert.Equal(t, strings.TrimSpace(expectedYaml), strings.TrimSpace(string(yamlData)))
}
