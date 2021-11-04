package old

import (
	"gopkg.in/specgen-io/yaml.v3"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_Union_Unmarshal(t *testing.T) {
	data := `
description: OneOf description
oneOf:
  first: TheFirst
  second: TheSecond
  third: TheThird
`
	var oneOf = OneOf{}
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &oneOf)
	assert.Equal(t, err, nil)
	assert.Equal(t, *oneOf.Description, "OneOf description")
	assert.Equal(t, oneOf.Discriminator == nil, true)
	assert.Equal(t, len(oneOf.Items), 3)
	item1 := oneOf.Items[0]
	item2 := oneOf.Items[1]
	item3 := oneOf.Items[2]
	assert.Equal(t, item1.Name.Source, "first")
	assert.Equal(t, item1.Type.Definition, ParseType("TheFirst"))
	assert.Equal(t, item2.Name.Source, "second")
	assert.Equal(t, item2.Type.Definition, ParseType("TheSecond"))
	assert.Equal(t, item3.Name.Source, "third")
	assert.Equal(t, item3.Type.Definition, ParseType("TheThird"))
}

func Test_Union_Unmarshal_With_Discriminator(t *testing.T) {
	data := `
description: OneOf description
discriminator: type
oneOf:
  first: TheFirst
  second: TheSecond
  third: TheThird
`
	var oneOf = OneOf{}
	err := yaml.UnmarshalWith(decodeStrict, []byte(data), &oneOf)
	assert.Equal(t, err, nil)
	assert.Equal(t, *oneOf.Description, "OneOf description")
	assert.Equal(t, *oneOf.Discriminator, "type")
	assert.Equal(t, len(oneOf.Items), 3)
	item1 := oneOf.Items[0]
	item2 := oneOf.Items[1]
	item3 := oneOf.Items[2]
	assert.Equal(t, item1.Name.Source, "first")
	assert.Equal(t, item1.Type.Definition, ParseType("TheFirst"))
	assert.Equal(t, item2.Name.Source, "second")
	assert.Equal(t, item2.Type.Definition, ParseType("TheSecond"))
	assert.Equal(t, item3.Name.Source, "third")
	assert.Equal(t, item3.Type.Definition, ParseType("TheThird"))
}

func Test_OneOf_Discriminator_Marshal(t *testing.T) {
	expectedYaml := strings.TrimLeft(`
discriminator: type
oneOf:
  first: TheFirst
  second: TheSecond
  third: TheThird
`, "\n")
	var oneOf OneOf
	checkUnmarshalMarshal(t, expectedYaml, &oneOf)
}
