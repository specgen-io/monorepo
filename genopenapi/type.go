package genopenapi

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
	"github.com/specgen-io/specgen/v2/yamlx"
	"strconv"
)

func OpenApiType(typ *spec.TypeDef, defaultValue *string) *yamlx.YamlMap {
	switch typ.Node {
	case spec.PlainType:
		result := PlainOpenApiType(typ.Info, typ.Plain)
		if defaultValue != nil {
			result.Add("default", DefaultValue(typ, *defaultValue))
		}
		return result
	case spec.NullableType:
		child := OpenApiType(typ.Child, defaultValue)
		return child
	case spec.ArrayType:
		child := OpenApiType(typ.Child, nil)
		result := yamlx.Map()
		result.Add("type", "array")
		result.Add("items", child)
		return result
	case spec.MapType:
		child := OpenApiType(typ.Child, nil)
		result := yamlx.Map()
		result.Add("type", "object")
		result.Add("additionalProperties", child)
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
			if err != nil {
				failDefaultParse(typ, valueStr, err)
			}
			return int32(value)
		case spec.TypeInt64:
			value, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				failDefaultParse(typ, valueStr, err)
			}
			return value
		case spec.TypeFloat:
			value, err := strconv.ParseFloat(valueStr, 32)
			if err != nil {
				failDefaultParse(typ, valueStr, err)
			}
			return float32(value)
		case spec.TypeDouble:
			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				failDefaultParse(typ, valueStr, err)
			}
			return value
		case spec.TypeDecimal:
			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				failDefaultParse(typ, valueStr, err)
			}
			return value
		case spec.TypeBoolean:
			value, err := strconv.ParseBool(valueStr)
			if err != nil {
				failDefaultParse(typ, valueStr, err)
			}
			return value
		case
			spec.TypeString,
			spec.TypeUuid,
			spec.TypeDate,
			spec.TypeDateTime:
			return valueStr
		default:
			model := typ.Info.Model
			if model != nil && model.IsEnum() {
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

func PlainOpenApiType(typeInfo *spec.TypeInfo, typ string) *yamlx.YamlMap {
	switch typ {
	case spec.TypeInt32:
		result := yamlx.Map()
		result.Add("type", "integer")
		result.Add("format", "int32")
		return result
	case spec.TypeInt64:
		result := yamlx.Map()
		result.Add("type", "integer")
		result.Add("format", "int64")
		return result
	case spec.TypeFloat:
		result := yamlx.Map()
		result.Add("type", "number")
		result.Add("format", "float")
		return result
	case spec.TypeDouble:
		result := yamlx.Map()
		result.Add("type", "number")
		result.Add("format", "double")
		return result
	case spec.TypeDecimal:
		result := yamlx.Map()
		result.Add("type", "number")
		result.Add("format", "decimal")
		return result
	case spec.TypeBoolean:
		result := yamlx.Map()
		result.Add("type", "boolean")
		return result
	case spec.TypeString:
		result := yamlx.Map()
		result.Add("type", "string")
		return result
	case spec.TypeUuid:
		result := yamlx.Map()
		result.Add("type", "string")
		result.Add("format", "uuid")
		return result
	case spec.TypeDate:
		result := yamlx.Map()
		result.Add("type", "string")
		result.Add("format", "date")
		return result
	case spec.TypeDateTime:
		result := yamlx.Map()
		result.Add("type", "string")
		result.Add("format", "datetime")
		return result
	case spec.TypeJson:
		result := yamlx.Map()
		result.Add("type", "object")
		return result
	default:
		result := yamlx.Map()
		result.Add("$ref", "#/components/schemas/"+versionedModelName(typeInfo.Model.Version.Version.Source, typ))
		return result
	}
}
