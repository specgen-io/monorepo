package genscala

import (
	"github.com/ModaOperandi/spec"
	"gotest.tools/assert"
	"testing"
)

func CheckPlainType(t *testing.T, plainType string, expected string) {
	typ := spec.Type{Node: spec.PlainType, PlainType: plainType}
	scalaType := ScalaType(&typ)
	assert.Equal(t, scalaType, expected)
}

func TestPlainTypeByte(t *testing.T) {
	CheckPlainType(t, spec.TypeByte, "Byte")
}

func TestPlainTypeShort(t *testing.T) {
	CheckPlainType(t, spec.TypeShort, "Short")
	CheckPlainType(t, spec.TypeInt16, "Short")
}

func TestPlainTypeInt(t *testing.T) {
	CheckPlainType(t, spec.TypeInt, "Int")
	CheckPlainType(t, spec.TypeInt32, "Int")
}

func TestPlainTypeLong(t *testing.T) {
	CheckPlainType(t, spec.TypeLong, "Long")
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
	CheckPlainType(t, spec.TypeBool, "Boolean")
}

func TestPlainTypeDate(t *testing.T) {
	CheckPlainType(t, spec.TypeDate, "java.time.LocalDate")
}

func TestPlainTypeDateTime(t *testing.T) {
	CheckPlainType(t, spec.TypeDateTime, "java.time.LocalDateTime")
}

func TestPlainTypeTime(t *testing.T) {
	CheckPlainType(t, spec.TypeTime, "java.time.LocalTime")
}

func TestPlainTypeJson(t *testing.T) {
	CheckPlainType(t, spec.TypeJson, "JsonNode")
}

func TestPlainTypeChar(t *testing.T) {
	CheckPlainType(t, spec.TypeChar, "Char")
}

func TestPlainTypeString(t *testing.T) {
	CheckPlainType(t, spec.TypeString, "String")
	CheckPlainType(t, spec.TypeStr, "String")
}

func TestNullableType(t *testing.T) {
	typ := spec.Type{Node: spec.NullableType, Child: &spec.Type{Node: spec.PlainType, PlainType: spec.TypeString}}
	scalaType := ScalaType(&typ)
	assert.Equal(t, scalaType, "Option[String]")
}

func TestArrayType(t *testing.T) {
	typ := spec.Type{Node: spec.ArrayType, Child: &spec.Type{Node: spec.PlainType, PlainType: spec.TypeString}}
	scalaType := ScalaType(&typ)
	assert.Equal(t, scalaType, "List[String]")
}

func TestMapType(t *testing.T) {
	typ := spec.Type{Node: spec.MapType, Child: &spec.Type{Node: spec.PlainType, PlainType: "Model"}}
	scalaType := ScalaType(&typ)
	assert.Equal(t, scalaType, "Map[String, Model]")
}
