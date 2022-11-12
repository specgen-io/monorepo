package types

import (
	"gotest.tools/assert"
	"spec"
	"testing"
)

func CheckJacksonType(t *testing.T, typ *spec.TypeDef, expected string) {
	types := Types{"JsonNode"}
	javaType := types.Kotlin(typ)
	assert.Equal(t, javaType, expected)
}

func TestPlainTypeInt(t *testing.T) {
	CheckJacksonType(t, spec.Plain(spec.TypeInt32), "Int")
}

func TestPlainTypeLong(t *testing.T) {
	CheckJacksonType(t, spec.Plain(spec.TypeInt64), "Long")
}

func TestPlainTypeFloat(t *testing.T) {
	CheckJacksonType(t, spec.Plain(spec.TypeFloat), "Float")
}

func TestPlainTypeDouble(t *testing.T) {
	CheckJacksonType(t, spec.Plain(spec.TypeDouble), "Double")
}

func TestPlainTypeDecimal(t *testing.T) {
	CheckJacksonType(t, spec.Plain(spec.TypeDecimal), "BigDecimal")
}

func TestPlainTypeBoolean(t *testing.T) {
	CheckJacksonType(t, spec.Plain(spec.TypeBoolean), "Boolean")
}

func TestPlainTypeString(t *testing.T) {
	CheckJacksonType(t, spec.Plain(spec.TypeString), "String")
}

func TestPlainTypeUuid(t *testing.T) {
	CheckJacksonType(t, spec.Plain(spec.TypeUuid), "UUID")
}

func TestPlainTypeDate(t *testing.T) {
	CheckJacksonType(t, spec.Plain(spec.TypeDate), "LocalDate")
}

func TestPlainTypeDateTime(t *testing.T) {
	CheckJacksonType(t, spec.Plain(spec.TypeDateTime), "LocalDateTime")
}

func TestPlainTypeJson(t *testing.T) {
	CheckJacksonType(t, spec.Plain(spec.TypeJson), "JsonNode")
}

func TestPlainTypeEmpty(t *testing.T) {
	CheckJacksonType(t, spec.Plain(spec.TypeEmpty), "Void")
}

func TestNullableType(t *testing.T) {
	CheckJacksonType(t, spec.Nullable(spec.Plain(spec.TypeInt32)), "Int?")
}

func TestArrayType(t *testing.T) {
	CheckJacksonType(t, spec.Array(spec.Plain(spec.TypeString)), "List<String>")
}

func TestMapType(t *testing.T) {
	CheckJacksonType(t, spec.Map(spec.Plain("Model")), "Map<String, Model>")
}

func TestComplexType(t *testing.T) {
	CheckJacksonType(t, spec.Array(spec.Nullable(spec.Plain(spec.TypeBoolean))), "List<Boolean?>")
}

func TestMapScalarType(t *testing.T) {
	CheckJacksonType(t, spec.Map(spec.Plain(spec.TypeInt32)), "Map<String, Int>")
}
