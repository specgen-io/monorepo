package spec

import (
	"gotest.tools/assert"
	"testing"
)

func Test_Object_FieldUniqueness(t *testing.T) {
	data := `
models:
  MyObject:
    object:
      the_field: string
      theField: string
`

	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(spec)
	assert.Equal(t, len(errors), 0)

	_, errors = validate(spec)
	assert.Equal(t, len(errors), 1)
}

func Test_Object_EmptyIsNotAllowed(t *testing.T) {
	data := `
models:
  MyObject:
    object:
      the_field: empty
`

	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(spec)
	assert.Equal(t, len(errors), 0)

	_, errors = validate(spec)
	assert.Equal(t, len(errors), 1)
}

func Test_OneOf_ItemsUniqueness(t *testing.T) {
	data := `
models:
  MyUnion:
    oneOf:
      the_item: string
      theItem: string
`

	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(spec)
	assert.Equal(t, len(errors), 0)

	_, errors = validate(spec)
	assert.Equal(t, len(errors), 1)
}

func Test_OneOf_EmptyIsNotAllowed(t *testing.T) {
	data := `
models:
  MyUnion:
    oneOf:
      the_item: empty
`

	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(spec)
	assert.Equal(t, len(errors), 0)

	_, errors = validate(spec)
	assert.Equal(t, len(errors), 1)
}

func Test_Enum_ItemsUniqueness(t *testing.T) {
	data := `
models:
  MyUnion:
    enum:
      - the_item
      - the_item
`

	spec, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(spec)
	assert.Equal(t, len(errors), 0)

	_, errors = validate(spec)
	assert.Equal(t, len(errors), 1)
}
