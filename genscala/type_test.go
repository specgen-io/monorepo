package genscala

import (
	spec "github.com/specgen-io/spec.v2"
	"gotest.tools/assert"
	"testing"
)

func CheckPlainType(t *testing.T, plainType string, expected string) {
	typ := spec.Plain(plainType)
	scalaType := ScalaType(typ)
	assert.Equal(t, scalaType, expected)
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

func TestPlainTypeDate(t *testing.T) {
	CheckPlainType(t, spec.TypeDate, "java.time.LocalDate")
}

func TestPlainTypeDateTime(t *testing.T) {
	CheckPlainType(t, spec.TypeDateTime, "java.time.LocalDateTime")
}

func TestPlainTypeJson(t *testing.T) {
	CheckPlainType(t, spec.TypeJson, "io.circe.Json")
}

func TestPlainTypeString(t *testing.T) {
	CheckPlainType(t, spec.TypeString, "String")
}

func TestNullableType(t *testing.T) {
	typ := spec.Nullable(spec.Plain(spec.TypeString))
	scalaType := ScalaType(typ)
	assert.Equal(t, scalaType, "Option[String]")
}

func TestArrayType(t *testing.T) {
	typ := spec.Array(spec.Plain(spec.TypeString))
	scalaType := ScalaType(typ)
	assert.Equal(t, scalaType, "List[String]")
}

func TestMapType(t *testing.T) {
	typ := spec.Map(spec.Plain("Model"))
	scalaType := ScalaType(typ)
	assert.Equal(t, scalaType, "Map[String, Model]")
}
