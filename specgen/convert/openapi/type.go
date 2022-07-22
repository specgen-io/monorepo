package openapi

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/specgen-io/specgen/spec/v2"
	"strings"
)

func specType(schema *openapi3.SchemaRef, required bool) *spec.Type {
	return &spec.Type{*specTypeDef(schema, required), nil}
}

func specTypeDef(schema *openapi3.SchemaRef, required bool) *spec.TypeDef {
	if !required {
		return spec.Nullable(specTypeDef(schema, true))
	}
	if schema.Ref != "" {
		return spec.Plain(strings.TrimPrefix(schema.Ref, "#/components/schemas/"))
	} else {
		switch schema.Value.Type {
		case TypeArray:
			return spec.Array(specTypeDef(schema.Value.Items, true))
		case TypeObject:
			if schema.Value.AdditionalProperties != nil {
				return spec.Map(specTypeDef(schema.Value.AdditionalProperties, true))
			} else {
				panic(fmt.Sprintf(`object schema is not supported yet`))
			}
			return spec.Plain("object")
		default:
			return specPlainType(schema.Value.Type, schema.Value.Format)
		}
	}
}

func specPlainType(typ string, format string) *spec.TypeDef {
	switch typ {
	case TypeBoolean:
		switch format {
		case "":
			return spec.Plain(spec.TypeBoolean)
		default:
			panic(fmt.Sprintf("Unknown format of boolean: %s", format))
		}
	case TypeString:
		switch format {
		case "":
			return spec.Plain(spec.TypeString)
		case FormatUuid:
			return spec.Plain(spec.TypeUuid)
		case FormatDate:
			return spec.Plain(spec.TypeDate)
		case FormatDateTime:
			return spec.Plain(spec.TypeDateTime)
		default:
			panic(fmt.Sprintf("Unknown format of string: %s", format))
		}
	case TypeNumber:
		switch format {
		case FormatFloat:
			return spec.Plain(spec.TypeFloat)
		case FormatDouble:
			return spec.Plain(spec.TypeDouble)
		case FormatDecimal:
			return spec.Plain(spec.TypeDecimal)
		default:
			panic(fmt.Sprintf("Unknown format of number: %s", format))
		}
	case TypeInteger:
		switch format {
		case FormatInt64:
			return spec.Plain(spec.TypeInt64)
		case FormatInt32:
			return spec.Plain(spec.TypeInt32)
		case "":
			return spec.Plain(spec.TypeInt32)
		default:
			panic(fmt.Sprintf("Unknown format of integer: %s", format))
		}
	default:
		panic(fmt.Sprintf("Unknown type: %s", typ))
	}
	return nil
}

const (
	TypeArray   string = "array"
	TypeObject  string = "object"
	TypeInteger string = "integer"
	TypeNumber  string = "number"
	TypeString  string = "string"
	TypeBoolean string = "boolean"
)

const (
	FormatInt64    string = "int64"
	FormatInt32    string = "int32"
	FormatFloat    string = "float"
	FormatDouble   string = "double"
	FormatDecimal  string = "decimal"
	FormatUuid     string = "uuid"
	FormatDate     string = "date"
	FormatDateTime string = "date-time"
)
