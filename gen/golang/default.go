package golang

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/v2/spec"
)

func DefaultValue(typ *spec.TypeDef, value string) string {
	switch typ.Node {
	case spec.ArrayType:
		if value == "[]" {
			return "[]"
		} else {
			panic(fmt.Sprintf("Array type %s default value %s is not supported", typ.Name, value))
		}
	case spec.MapType:
		if value == "{}" {
			return "{}"
		} else {
			panic(fmt.Sprintf("Map type %s default value %s is not supported", typ.Name, value))
		}
	case spec.PlainType:
		switch typ.Plain {
		case
			spec.TypeInt32,
			spec.TypeInt64,
			spec.TypeFloat,
			spec.TypeDouble:
			return value
		case
			spec.TypeDecimal:
			return `decimal.NewFromString("` + value + `")`
		case spec.TypeBoolean:
			return value
		case spec.TypeString:
			return `"` + value + `"`
		case spec.TypeUuid:
			return `uuid.Parse("` + value + `")`
		case spec.TypeDate:
			return `civil.ParseDate("` + value + `")`
		case spec.TypeDateTime:
			return `civil.ParseDateTime("` + value + `")`
		default:
			model := typ.Info.Model
			if model != nil && model.IsEnum() {
				return model.Name.PascalCase() + `.` + casee.ToSnakeCase(value)
			} else {
				panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
			}
		}
	default:
		panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
	}
}
