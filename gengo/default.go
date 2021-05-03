package gengo

import (
	"fmt"
	"github.com/pinzolo/casee"
	spec "github.com/specgen-io/spec.v2"
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
			return `guuid.New(` + value + `)`
		case spec.TypeDate:
			return `time.Parse("2006-01-02", "` + value + `")`
		case spec.TypeDateTime:
			return `time.Parse("2006-01-02T15:04:05", "` + value + `")`
		default:
			modelInfo := typ.Info.ModelInfo
			if modelInfo != nil && modelInfo.Model.IsEnum() {
				return modelInfo.Model.Name.PascalCase() + `.` + casee.ToSnakeCase(value)
			} else {
				panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
			}
		}
	default:
		panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
	}
}
