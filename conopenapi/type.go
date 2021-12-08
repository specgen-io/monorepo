package conopenapi

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func SpecTypeFromSchemaRef(schema *openapi3.SchemaRef) *spec.TypeDef {
	if schema.Value == nil {
		return spec.Plain(strings.TrimPrefix(schema.Ref, "#/components/schemas/"))
	} else {
		return spec.Plain("problem") //TODO: This is not correct
	}
}

func SpecTypeFromSchema(schema *openapi3.Schema) *spec.TypeDef {
	switch schema.Type {
	case TypeBoolean:
		switch schema.Format {
		case "":
			return spec.Plain(spec.TypeBoolean)
		default:
			panic(fmt.Sprintf("Unknown format of boolean: %s", schema.Format))
		}
	case TypeString:
		switch schema.Format {
		case "":
			return spec.Plain(spec.TypeString)
		case FormatUuid:
			return spec.Plain(spec.TypeUuid)
		case FormatDate:
			return spec.Plain(spec.TypeDate)
		case FormatDateTime:
			return spec.Plain(spec.TypeDateTime)
		default:
			panic(fmt.Sprintf("Unknown format of string: %s", schema.Format))
		}
	case TypeNumber:
		switch schema.Format {
		case FormatFloat:
			return spec.Plain(spec.TypeFloat)
		case FormatDouble:
			return spec.Plain(spec.TypeDouble)
		case FormatDecimal:
			return spec.Plain(spec.TypeDecimal)
		default:
			panic(fmt.Sprintf("Unknown format of number: %s", schema.Format))
		}
	case TypeInteger:
		switch schema.Format {
		case FormatInt64:
			return spec.Plain(spec.TypeInt64)
		case FormatInt32:
			return spec.Plain(spec.TypeInt32)
		default:
			panic(fmt.Sprintf("Unknown format of integer: %s", schema.Format))
		}
	default:
		panic(fmt.Sprintf("Unknown type: %s", schema.Type))
	}
	return nil
}

const (
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
	FormatDateTime string = "datetime"
)
