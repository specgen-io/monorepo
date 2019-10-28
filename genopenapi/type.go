package genopenapi

import (
	"fmt"
	"github.com/ModaOperandi/spec"
)

func OpenApiType(typ *spec.Type, defaultValue *string) *YamlMap {
	switch typ.Node {
	case spec.PlainType:
		result := PlainOpenApiType(typ.PlainType)
		if defaultValue != nil {
			result.Set("default", defaultValue)
		}
		return result
	case spec.NullableType:
		child := OpenApiType(typ.Child, defaultValue)
		return child
	case spec.ArrayType:
		child := OpenApiType(typ.Child, nil)
		result := Map().Set("type", "array").Set("items", child.Yaml)
		return result
	case spec.MapType:
		child := OpenApiType(typ.Child, nil)
		result := Map().Set("type", "object").Set("additionalProperties", child.Yaml)
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainOpenApiType(typ string) *YamlMap {
	switch typ {
	case spec.TypeByte:
		return Map().Set("type", "integer").Set("format", "int8")
	case spec.TypeShort, spec.TypeInt16:
		return Map().Set("type", "integer").Set("format", "int16")
	case spec.TypeInt, spec.TypeInt32:
		return Map().Set("type", "integer").Set("format", "int32")
	case spec.TypeLong, spec.TypeInt64:
		return Map().Set("type", "integer").Set("format", "int64")
	case spec.TypeFloat:
		return Map().Set("type", "number").Set("format", "float")
	case spec.TypeDouble:
		return Map().Set("type", "number").Set("format", "double")
	case spec.TypeDecimal:
		return Map().Set("type", "number").Set("format", "decimal")
	case spec.TypeBool, spec.TypeBoolean:
		return Map().Set("type", "boolean")
	case spec.TypeChar:
		return Map().Set("type", "string").Set("format", "char")
	case spec.TypeString, spec.TypeStr:
		return Map().Set("type", "string")
	case spec.TypeUuid:
		return Map().Set("type", "string").Set("format", "uuid")
	case spec.TypeDate:
		return Map().Set("type", "string").Set("format", "date")
	case spec.TypeDateTime:
		return Map().Set("type", "string").Set("format", "datetime")
	case spec.TypeTime:
		return Map().Set("type", "string").Set("format", "time")
	case spec.TypeJson:
		return Map().Set("type", "object")
	default:
		return Map().Set("$ref", "#/components/schemas/"+typ)
	}
}
