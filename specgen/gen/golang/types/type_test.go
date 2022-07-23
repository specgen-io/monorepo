package types

import (
	"github.com/specgen-io/specgen/v2/spec"
	"gotest.tools/assert"
	"testing"
)

func CheckPlainType(t *testing.T, plainType string, expected string) {
	typ := spec.Plain(plainType)
	goType := GoType(typ)
	assert.Equal(t, goType, expected)
}

func TestPlainTypeInt(t *testing.T) {
	CheckPlainType(t, spec.TypeInt32, "int")
}

func TestPlainTypeLong(t *testing.T) {
	CheckPlainType(t, spec.TypeInt64, "int64")
}

func TestPlainTypeFloat(t *testing.T) {
	CheckPlainType(t, spec.TypeFloat, "float32")
}

func TestPlainTypeDouble(t *testing.T) {
	CheckPlainType(t, spec.TypeDouble, "float64")
}

func TestPlainTypeDecimal(t *testing.T) {
	CheckPlainType(t, spec.TypeDecimal, "decimal.Decimal")
}

func TestPlainTypeBoolean(t *testing.T) {
	CheckPlainType(t, spec.TypeBoolean, "bool")
}

func TestPlainTypeDate(t *testing.T) {
	CheckPlainType(t, spec.TypeDate, "civil.Date")
}

func TestPlainTypeDateTime(t *testing.T) {
	CheckPlainType(t, spec.TypeDateTime, "civil.DateTime")
}

func TestPlainTypeJson(t *testing.T) {
	CheckPlainType(t, spec.TypeJson, "json.RawMessage")
}

func TestPlainTypeString(t *testing.T) {
	CheckPlainType(t, spec.TypeString, "string")
}

func TestPlainTypeEmpty(t *testing.T) {
	CheckPlainType(t, spec.TypeEmpty, "empty.Type")
}

func TestNullableType(t *testing.T) {
	typ := spec.Nullable(spec.Plain(spec.TypeString))
	goType := GoType(typ)
	assert.Equal(t, goType, "*string")
}

func TestArrayType(t *testing.T) {
	typ := spec.Array(spec.Plain(spec.TypeString))
	goType := GoType(typ)
	assert.Equal(t, goType, "[]string")
}

func TestMapType(t *testing.T) {
	typ := spec.Map(spec.Plain("Model"))
	goType := GoType(typ)
	assert.Equal(t, goType, "map[string]models.Model")
}
