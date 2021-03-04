package genopenapi

import (
	"fmt"
	"github.com/specgen-io/spec"
	"strconv"
)

func OpenApiType(typ *spec.TypeDef, defaultValue *string) *YamlMap {
	switch typ.Node {
	case spec.PlainType:
		result := PlainOpenApiType(typ.Plain)
		if defaultValue != nil {
			result.Set("default", DefaultValue(typ, *defaultValue))
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
func DefaultValue(typ *spec.TypeDef, valueStr string) interface{} {
	switch typ.Node {
	case spec.ArrayType:
		if valueStr == "[]" {
			return []interface{}{}
		} else {
			panic(fmt.Sprintf("Array type %s default value %s is not supported", typ.Name, valueStr))
		}
	case spec.MapType:
		if valueStr == "{}" {
			return map[string]interface{}{}
		} else {
			panic(fmt.Sprintf("Map type %s default value %s is not supported", typ.Name, valueStr))
		}
	case spec.PlainType:
		switch typ.Plain {
		case spec.TypeInt32:
			value, err := strconv.ParseInt(valueStr, 10, 32)
			if err != nil { failDefaultParse(typ, valueStr, err) }
			return int32(value)
		case spec.TypeInt64:
			value, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil { failDefaultParse(typ, valueStr, err) }
			return value
		case spec.TypeFloat:
			value, err := strconv.ParseFloat(valueStr, 32)
			if err != nil { failDefaultParse(typ, valueStr, err) }
			return float32(value)
		case spec.TypeDouble:
			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil { failDefaultParse(typ, valueStr, err) }
			return value
		case spec.TypeDecimal:
			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil { failDefaultParse(typ, valueStr, err) }
			return value
		case spec.TypeBoolean:
			value, err := strconv.ParseBool(valueStr)
			if err != nil { failDefaultParse(typ, valueStr, err) }
			return value
		case
			spec.TypeString,
			spec.TypeUuid,
			spec.TypeDate,
			spec.TypeDateTime:
			return valueStr
		default:
			if typ.Info.Model != nil && typ.Info.Model.IsEnum() {
				return valueStr
			} else {
				panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
			}
		}
	default:
		panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
	}
}

func failDefaultParse(typ *spec.TypeDef, defaultValue string, err error) {
	panic(fmt.Sprintf("Parsing default value %s for type %s failed, message: %s", defaultValue, typ.Name, err.Error()))
}

func PlainOpenApiType(typ string) *YamlMap {
	switch typ {
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
	case spec.TypeString:
		return Map().Set("type", "string")
	case spec.TypeUuid:
		return Map().Set("type", "string").Set("format", "uuid")
	case spec.TypeDate:
		return Map().Set("type", "string").Set("format", "date")
	case spec.TypeDateTime:
		return Map().Set("type", "string").Set("format", "datetime")
	case spec.TypeJson:
		return Map().Set("type", "object")
	default:
		return Map().Set("$ref", "#/components/schemas/"+typ)
	}
}
