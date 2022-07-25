package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_Models_Unmarshal_WrongNameFormat(t *testing.T) {
	data := `
model_one:
  description: some model
  object:
    prop1: string
    prop2: int32
`
	var models Models
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &models)
	assert.ErrorContains(t, err, "model_one")
}

func Test_Models_Marshal(t *testing.T) {
	expectedYaml := strings.TrimLeft(`
Model1:
  description: first model
  object:
    prop1: string
    prop2: int
Model2:
  description: second model
  enum:
    - first
    - second
    - third
Model3:
  description: third model
  enum:
    first: FIRST
    second: SECOND
    third: THIRD
Model4:
  description: forth model
  oneOf:
    one: Model1
    two: Model2
Model5:
  description: fifth model
  discriminator: type
  oneOf:
    one: Model1
    two: Model2
`, "\n")
	var models Models
	checkUnmarshalMarshal(t, expectedYaml, &models)
}
