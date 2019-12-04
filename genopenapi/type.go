package genopenapi

import (
	"fmt"
	"github.com/ModaOperandi/spec"
	"strconv"
)

func OpenApiType(typ *spec.TypeDef, defaultValue *string) *YamlMap {
	switch typ.Node {
	case spec.PlainType:
		result := PlainOpenApiType(typ.Plain)
		if defaultValue != nil {
			result.Set("default", DefaultValue(typ.Plain, *defaultValue))
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

// TODO: Default values provided in spec are actually exactly what should appear in Swagger in the string form.
//       This value translation string -> go type -> string is only needed because of how OpenAPI spec if formed.
//       It's possible that after move to yaml.v3 this value translation can be avoided.
func DefaultValue(typ string, defaultValue string) interface{} {
	switch typ {
	case spec.TypeByte:
		value, err := strconv.ParseInt(defaultValue, 10, 8)
		if err != nil { failDefaultParse(typ, defaultValue, err) }
		return byte(value)
	case spec.TypeInt16:
		value, err := strconv.ParseInt(defaultValue, 10, 16)
		if err != nil { failDefaultParse(typ, defaultValue, err) }
		return int16(value)
	case spec.TypeInt32:
		value, err := strconv.ParseInt(defaultValue, 10, 32)
		if err != nil { failDefaultParse(typ, defaultValue, err) }
		return int32(value)
	case spec.TypeInt64:
		value, err := strconv.ParseInt(defaultValue, 10, 64)
		if err != nil { failDefaultParse(typ, defaultValue, err) }
		return value
	case spec.TypeFloat:
		value, err := strconv.ParseFloat(defaultValue, 32)
		if err != nil { failDefaultParse(typ, defaultValue, err) }
		return float32(value)
	case spec.TypeDouble:
		value, err := strconv.ParseFloat(defaultValue, 64)
		if err != nil { failDefaultParse(typ, defaultValue, err) }
		return value
	case spec.TypeDecimal:
		value, err := strconv.ParseFloat(defaultValue, 64)
		if err != nil { failDefaultParse(typ, defaultValue, err) }
		return value
	case spec.TypeBoolean:
		value, err := strconv.ParseBool(defaultValue)
		if err != nil { failDefaultParse(typ, defaultValue, err) }
		return value
	case
		spec.TypeChar,
		spec.TypeString,
		spec.TypeUuid,
		spec.TypeDate,
		spec.TypeDateTime,
		spec.TypeTime:
		return defaultValue
	default:
		panic(fmt.Sprintf("type: %s does not support default value", typ))
	}
}

func failDefaultParse(typ string, defaultValue string, err error) {
	panic(fmt.Sprintf("Parsing default value %s for type %s failed, message: %s", defaultValue, typ, err.Error()))
}

func PlainOpenApiType(typ string) *YamlMap {
	switch typ {
	case spec.TypeByte:
		return Map().Set("type", "integer").Set("format", "int8")
	case spec.TypeInt16:
		return Map().Set("type", "integer").Set("format", "int16")
	case spec.TypeInt32:
		return Map().Set("type", "integer").Set("format", "int32")
	case spec.TypeInt64:
		return Map().Set("type", "integer").Set("format", "int64")
	case spec.TypeFloat:
		return Map().Set("type", "number").Set("format", "float")
	case spec.TypeDouble:
		return Map().Set("type", "number").Set("format", "double")
	case spec.TypeDecimal:
		return Map().Set("type", "number").Set("format", "decimal")
	case spec.TypeBoolean:
		return Map().Set("type", "boolean")
	case spec.TypeChar:
		return Map().Set("type", "string").Set("format", "char")
	case spec.TypeString:
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
