package genopenapi

import (
	"github.com/specgen-io/spec"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func checkType(t *testing.T, typ *spec.TypeDef, expected string) {
	openApiType := OpenApiType(typ, nil)
	openApiTypeYaml, err := ToYamlString(openApiType)
	assert.NilError(t, err)
	assert.Equal(t, strings.TrimSpace(openApiTypeYaml), strings.TrimSpace(expected))
}

func TestPlainTypeInt(t *testing.T) {
	expected := `
type: integer
format: int32
`
	checkType(t, spec.Plain(spec.TypeInt32), expected)
}

func TestPlainTypeLong(t *testing.T) {
	expected := `
type: integer
format: int64
`
	checkType(t, spec.Plain(spec.TypeInt64), expected)
}

func TestPlainTypeFloat(t *testing.T) {
	expected := `
type: number
format: float
`
	checkType(t, spec.Plain(spec.TypeFloat), expected)
}

func TestPlainTypeDouble(t *testing.T) {
	expected := `
type: number
format: double
`
	checkType(t, spec.Plain(spec.TypeDouble), expected)
}

func TestPlainTypeDecimal(t *testing.T) {
	expected := `
type: number
format: decimal
`
	checkType(t, spec.Plain(spec.TypeDecimal), expected)
}

func TestPlainTypeBoolean(t *testing.T) {
	expected := `type: boolean`
	checkType(t, spec.Plain(spec.TypeBoolean), expected)
}

func TestPlainTypeDate(t *testing.T) {
	expected := `
type: string
format: date
`
	checkType(t, spec.Plain(spec.TypeDate), expected)
}

func TestPlainTypeDateTime(t *testing.T) {
	expected := `
type: string
format: datetime
`
	checkType(t, spec.Plain(spec.TypeDateTime), expected)
}

func TestPlainTypeJson(t *testing.T) {
	expected := `type: object`
	checkType(t, spec.Plain(spec.TypeJson), expected)
}

func TestPlainTypeString(t *testing.T) {
	expected := `type: string`
	checkType(t, spec.Plain(spec.TypeString), expected)
}

func TestArrayType(t *testing.T) {
	expected := `
type: array
items:
  type: string
`
	typ := spec.Array(spec.Plain(spec.TypeString))
	checkType(t, typ, expected)
}

func TestMapType(t *testing.T) {
	expected := `
type: object
additionalProperties:
  type: string
`
	typ := spec.Map(spec.Plain(spec.TypeString))
	checkType(t, typ, expected)
}
