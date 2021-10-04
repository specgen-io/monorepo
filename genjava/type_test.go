package genjava

import (
	"github.com/specgen-io/spec"
	"gotest.tools/assert"
	"testing"
)

func CheckPlainType(t *testing.T, plainType string, expected string) {
	typ := spec.Plain(plainType)
	javaType := JavaType(typ)
	assert.Equal(t, javaType, expected)
}

func TestPlainTypeInt(t *testing.T) {
	CheckPlainType(t, spec.TypeInt32, "int")
}

func TestPlainTypeLong(t *testing.T) {
	CheckPlainType(t, spec.TypeInt64, "long")
}

func TestPlainTypeFloat(t *testing.T) {
	CheckPlainType(t, spec.TypeFloat, "float")
}

func TestPlainTypeDouble(t *testing.T) {
	CheckPlainType(t, spec.TypeDouble, "double")
}

func TestPlainTypeDecimal(t *testing.T) {
	CheckPlainType(t, spec.TypeDecimal, "BigDecimal")
}

func TestPlainTypeBoolean(t *testing.T) {
	CheckPlainType(t, spec.TypeBoolean, "boolean")
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
	CheckPlainType(t, spec.TypeJson, "Object")
}

func TestPlainTypeEmpty(t *testing.T) {
	CheckPlainType(t, spec.TypeEmpty, "void")
}

func TestNullableType(t *testing.T) {
	typ := spec.Nullable(spec.Plain(spec.TypeInt32))
	javaType := JavaType(typ)
	assert.Equal(t, javaType, "Integer")
}

func TestArrayType(t *testing.T) {
	typ := spec.Array(spec.Plain(spec.TypeString))
	javaType := JavaType(typ)
	assert.Equal(t, javaType, "String[]")
}

func TestMapType(t *testing.T) {
	typ := spec.Map(spec.Plain("Model"))
	javaType := JavaType(typ)
	assert.Equal(t, javaType, "Map<String, Model>")
}

func TestComplexType(t *testing.T) {
	typ := spec.Array(spec.Nullable(spec.Plain(spec.TypeBoolean)))
	javaType := JavaType(typ)
	assert.Equal(t, javaType, "Boolean[]")
}

func TestMapScalarType(t *testing.T) {
	typ := spec.Map(spec.Plain(spec.TypeInt32))
	javaType := JavaType(typ)
	assert.Equal(t, javaType, "Map<String, Integer>")
}
