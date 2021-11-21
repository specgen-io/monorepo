package genkotlin

import (
	"github.com/specgen-io/specgen/v2/spec"
	"gotest.tools/assert"
	"testing"
)

func CheckPlainType(t *testing.T, plainType string, expected string) {
	typ := spec.Plain(plainType)
	kotlinType := KotlinType(typ)
	assert.Equal(t, kotlinType, expected)
}

func TestPlainTypeInt(t *testing.T) {
	CheckPlainType(t, spec.TypeInt32, "Int")
}

func TestPlainTypeLong(t *testing.T) {
	CheckPlainType(t, spec.TypeInt64, "Long")
}

func TestPlainTypeFloat(t *testing.T) {
	CheckPlainType(t, spec.TypeFloat, "Float")
}

func TestPlainTypeDouble(t *testing.T) {
	CheckPlainType(t, spec.TypeDouble, "Double")
}

func TestPlainTypeDecimal(t *testing.T) {
	CheckPlainType(t, spec.TypeDecimal, "BigDecimal")
}

func TestPlainTypeBoolean(t *testing.T) {
	CheckPlainType(t, spec.TypeBoolean, "Boolean")
}

func TestPlainTypeString(t *testing.T) {
	CheckPlainType(t, spec.TypeString, "String")
}

func TestPlainTypeUuid(t *testing.T) {
	CheckPlainType(t, spec.TypeUuid, "UUID")
}

func TestPlainTypeDate(t *testing.T) {
	CheckPlainType(t, spec.TypeDate, "LocalDate")
}

func TestPlainTypeDateTime(t *testing.T) {
	CheckPlainType(t, spec.TypeDateTime, "LocalDateTime")
}

func TestPlainTypeJson(t *testing.T) {
	CheckPlainType(t, spec.TypeJson, "JsonNode")
}

func TestPlainTypeEmpty(t *testing.T) {
	CheckPlainType(t, spec.TypeEmpty, "Void")
}

func TestNullableType(t *testing.T) {
	typ := spec.Nullable(spec.Plain(spec.TypeInt32))
	javaType := KotlinType(typ)
	assert.Equal(t, javaType, "Int?")
}

func TestArrayType(t *testing.T) {
	typ := spec.Array(spec.Plain(spec.TypeString))
	javaType := KotlinType(typ)
	assert.Equal(t, javaType, "List<String>")
}

func TestMapType(t *testing.T) {
	typ := spec.Map(spec.Plain("Model"))
	javaType := KotlinType(typ)
	assert.Equal(t, javaType, "Map<String, Model>")
}

func TestComplexType(t *testing.T) {
	typ := spec.Array(spec.Nullable(spec.Plain(spec.TypeBoolean)))
	javaType := KotlinType(typ)
	assert.Equal(t, javaType, "List<Boolean?>")
}

func TestMapScalarType(t *testing.T) {
	typ := spec.Map(spec.Plain(spec.TypeInt32))
	javaType := KotlinType(typ)
	assert.Equal(t, javaType, "Map<String, Int>")
}
