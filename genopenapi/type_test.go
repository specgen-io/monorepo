package genopenapi

import (
	"github.com/ModaOperandi/spec"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func checkType(t *testing.T, typ *spec.Type, expected string) {
	openApiType := OpenApiType(typ, nil)
	assert.Equal(t, strings.TrimSpace(openApiType.String()), strings.TrimSpace(expected))
}

func TestPlainTypeByte(t *testing.T) {
	expected := `
type: integer
format: int8
`
	checkType(t, spec.Plain(spec.TypeByte), expected)
}

func TestPlainTypeShort(t *testing.T) {
	expected := `
type: integer
format: int16
`
	checkType(t, spec.Plain(spec.TypeShort), expected)
	checkType(t, spec.Plain(spec.TypeInt16), expected)
}

func TestPlainTypeInt(t *testing.T) {
	expected := `
type: integer
format: int32
`
	checkType(t, spec.Plain(spec.TypeInt), expected)
	checkType(t, spec.Plain(spec.TypeInt32), expected)
}

func TestPlainTypeLong(t *testing.T) {
	expected := `
type: integer
format: int64
`
	checkType(t, spec.Plain(spec.TypeLong), expected)
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
	checkType(t, spec.Plain(spec.TypeBool), expected)
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

func TestPlainTypeTime(t *testing.T) {
	expected := `
type: string
format: time
`
	checkType(t, spec.Plain(spec.TypeTime), expected)
}

func TestPlainTypeJson(t *testing.T) {
	expected := `type: object`
	checkType(t, spec.Plain(spec.TypeJson), expected)
}

func TestPlainTypeChar(t *testing.T) {
	expected := `
type: string
format: char
`
	checkType(t, spec.Plain(spec.TypeChar), expected)
}

func TestPlainTypeString(t *testing.T) {
	expected := `type: string`
	checkType(t, spec.Plain(spec.TypeString), expected)
	checkType(t, spec.Plain(spec.TypeStr), expected)
}

func TestCustomType(t *testing.T) {
	expected := `$ref: '#/components/schemas/MyCustomType'`
	checkType(t, spec.Plain("MyCustomType"), expected)
}

func TestArrayType(t *testing.T) {
	expected := `
type: array
items:
  type: string
`
	typ := spec.Type{Node: spec.ArrayType, Child: &spec.Type{Node: spec.PlainType, PlainType: spec.TypeString}}
	checkType(t, &typ, expected)
}

func TestMapType(t *testing.T) {
	expected := `
type: object
additionalProperties:
  $ref: '#/components/schemas/Model'
`
	typ := spec.Type{Node: spec.MapType, Child: &spec.Type{Node: spec.PlainType, PlainType: "Model"}}
	checkType(t, &typ, expected)
}
