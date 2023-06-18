package spec

import (
	"bytes"
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"testing"
)

func checkUnmarshalMarshal(t *testing.T, expectedYaml string, value interface{}) {
	err := yaml.UnmarshalWith(decodeStrict, []byte(expectedYaml), value)
	assert.Equal(t, err, nil)

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	err = encoder.Encode(value)
	actualYaml := buf.String()

	assert.Equal(t, expectedYaml, actualYaml)
}

func StrPtr(str string) *string {
	return &str
}
