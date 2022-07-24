package openapi

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/spec"
	"github.com/specgen-io/specgen/yamlx/v2"
)

func OpenApiType(types ...*spec.TypeDef) *yamlx.YamlMap {
	if len(types) > 1 {
		result := yamlx.Map()
		anyOfItems := yamlx.Array()
		for _, typ := range types {
			anyOfItems.Add(OpenApiType(typ))
		}
		result.Add("anyOf", anyOfItems)
		return result
	} else {
		typ := types[0]
		switch typ.Node {
		case spec.PlainType:
			result := PlainOpenApiType(typ.Info, typ.Plain)
			return result
		case spec.NullableType:
			child := OpenApiType(typ.Child)
			return child
		case spec.ArrayType:
			child := OpenApiType(typ.Child)
			result := yamlx.Map()
			result.Add("type", "array")
			result.Add("items", child)
			return result
		case spec.MapType:
			child := OpenApiType(typ.Child)
			result := yamlx.Map()
			result.Add("type", "object")
			result.Add("additionalProperties", child)
			return result
		default:
			panic(fmt.Sprintf("Unknown type: %v", typ))
		}
	}

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
		result.Add("format", "date-time")
		return result
	case spec.TypeJson:
		result := yamlx.Map()
		result.Add("type", "object")
		return result
	case spec.TypeEmpty:
		result := yamlx.Map()
		result.Add("type", "object")
		return result
	default:
		result := yamlx.Map()
		result.Add("$ref", componentSchemas(versionedModelName(typeInfo.Model.Version.Version.Source, typ)))
		return result
	}
}

func componentSchemas(name string) string {
	return "#/components/schemas/" + name
}
