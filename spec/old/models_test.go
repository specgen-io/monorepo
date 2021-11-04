package old

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_Models_Unmarshal(t *testing.T) {
	data := `
Model1:
  description: first model
  fields:
    prop1: string
    prop2: int32
Model2:      # second model
  prop1: string
  prop2: int32
Model3:
  enum:
  - first
  - second
  - third
Model4:
  oneOf:
    one: Model1
    two: Model2
`
	var models Models
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &models)
	assert.Equal(t, err, nil)

	assert.Equal(t, len(models), 4)
	model1 := models[0]
	model2 := models[1]
	model3 := models[2]
	model4 := models[3]

	assert.Equal(t, model1.Name.Source, "Model1")
	assert.Equal(t, model1.IsObject(), true)
	assert.Equal(t, *model1.Object.Description, "first model")
	assert.Equal(t, model2.Name.Source, "Model2")
	assert.Equal(t, model2.IsObject(), true)
	assert.Equal(t, *model2.Object.Description, "second model")
	assert.Equal(t, model3.Name.Source, "Model3")
	assert.Equal(t, model3.IsEnum(), true)
	assert.Equal(t, model4.Name.Source, "Model4")
	assert.Equal(t, model4.IsOneOf(), true)
}

func Test_Models_Unmarshal_WrongNameFormat(t *testing.T) {
	data := `
model_one:
  description: some model
  fields:
    prop1: string
    prop2: int32
`
	var models Models
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &models)
	assert.ErrorContains(t, err, "model_one")
}

func Test_Models_Marshal(t *testing.T) {
	expectedYaml := strings.TrimLeft(`
Model1: # first model
  object:
    prop1: string
    prop2: int
Model2: # second model
  enum:
    - first
    - second
    - third
Model3: # third model
  enum:
    first: FIRST
    second: SECOND
    third: THIRD
Model4: # forth model
  oneOf:
    one: Model1
    two: Model2
Model5:
  discriminator: type
  oneOf:
    one: Model1
    two: Model2
`, "\n")
	var models Models
	checkUnmarshalMarshal(t, expectedYaml, &models)
}
