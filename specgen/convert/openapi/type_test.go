package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/specgen-io/specgen/spec/v2"
	"gotest.tools/assert"
	"reflect"
	"testing"
)

type SpecTypeTest struct {
	OpenapiSchema string
	Required      bool
	SpecType      string
}

var specTypeTests = []SpecTypeTest{
	{
		`{"type": "integer", "format": "int32"}`,
		true,
		`int32`,
	},
	{
		`{"type": "number", "format": "float"}`,
		true,
		`float`,
	},
	{
		`{"type": "string"}`,
		true,
		`string`,
	},
	{
		`{"type": "string"}`,
		true,
		`string`,
	},
	{
		`{"type": "string", "format": "uuid"}`,
		true,
		`uuid`,
	},
	{
		`{"type": "string", "format": "date-time"}`,
		true,
		`datetime`,
	},
	{
		`{"type": "array", "items": {"type": "string"}}`,
		true,
		`string[]`,
	},
	{
		`{"type": "object", "additionalProperties": {"type": "string"}}`,
		true,
		`string{}`,
	},
}

func TestSpecType(t *testing.T) {
	for _, test := range specTypeTests {
		schema := &openapi3.SchemaRef{}
		err := schema.UnmarshalJSON([]byte(test.OpenapiSchema))
		assert.Equal(t, err, nil)
		assert.Equal(t, schema != nil, true)

		actual := specType(schema, test.Required)

		if test.SpecType != "" {
			expected := spec.NewType(test.SpecType)
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf(`openapi schema: %s, required: %t, expected spec type: %s: actual spec type: %s`, test.OpenapiSchema, test.Required, expected, actual)
			}
		}
	}
}
