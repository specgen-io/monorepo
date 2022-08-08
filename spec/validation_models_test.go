package spec

import (
	"errors"
	"testing"
)

func Test_Validation_Models(t *testing.T) {
	runReadSpecificationCases(t, validationModelsCases)
}

var validationModelsCases = []ReadSpecificationCase{
	{
		`object field non unique error`,
		`
models:
  MyObject:
    object:
      the_field: string
      theField: string
`,
		errors.New(`failed to validate specification`),
		[]Message{Error(`object model MyObject fields names are too similiar to each other: the_field, theField`).At(&Location{specificationMetaLines + 3, 5})},
		nil,
	},
	{
		`object field is empty error`,
		`
models:
  MyObject:
    object:
      the_field: empty
`,
		errors.New(`failed to validate specification`),
		[]Message{Error(`type empty can not be used in models`).At(&Location{specificationMetaLines + 4, 18})},
		nil,
	},
	{
		`oneOf items aren't unique error`,
		`
models:
  MyUnion:
    oneOf:
      the_item: string
      theItem: string
`,
		errors.New(`failed to validate specification`),
		[]Message{Error(`oneOf model MyUnion items names are too similiar to each other: the_item, theItem`).At(&Location{specificationMetaLines + 3, 5})},
		nil,
	},
	{
		`oneOf item is empty error`,
		`
models:
  MyUnion:
    oneOf:
      the_item: empty
`,
		errors.New(`failed to validate specification`),
		[]Message{Error(`type empty can not be used in models`).At(&Location{specificationMetaLines + 4, 17})},
		nil,
	},
	{
		`enum items aren't unique error`,
		`
models:
  MyUnion:
    enum:
      - the_item
      - the_item
`,
		errors.New(`failed to validate specification`),
		[]Message{Error(`enum model MyUnion items names are too similiar to each other: the_item, the_item`).At(&Location{specificationMetaLines + 3, 5})},
		nil,
	},
}
