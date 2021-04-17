package genscala

import (
	"fmt"
	spec "github.com/specgen-io/spec.v1"
	"github.com/pinzolo/casee"
)

func DefaultValue(typ *spec.TypeDef, value string) string {
	switch typ.Node {
	case spec.ArrayType:
		if value == "[]" {
			return "List()"
		} else {
			panic(fmt.Sprintf("Array type %s default value %s is not supported", typ.Name, value))
		}
	case spec.MapType:
		if value == "{}" {
			return "Map()"
		} else {
			panic(fmt.Sprintf("Map type %s default value %s is not supported", typ.Name, value))
		}
	case spec.PlainType:
		switch typ.Plain {
		case
			spec.TypeInt32,
			spec.TypeInt64,
			spec.TypeDouble:
			return value
		case spec.TypeFloat:
			return value + "f"
		case spec.TypeDecimal:
			return `BigDecimal("` + value + `")`
		case spec.TypeBoolean:
			return value
		case spec.TypeString:
			return `"` + value + `"`
		case spec.TypeUuid:
			return `UUID.fromString("` + value + `")`
		case spec.TypeDate:
			return `LocalDate.parse("` + value + `", DateTimeFormatter.ISO_LOCAL_DATE)`
		case spec.TypeDateTime:
			return `LocalDateTime.parse("` + value + `", DateTimeFormatter.ISO_LOCAL_DATE_TIME)`
		default:
			if typ.Info.Model != nil && typ.Info.Model.IsEnum() {
				return typ.Info.Model.Name.PascalCase() + `.` + casee.ToPascalCase(value)
			} else {
				panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
			}
		}
	default:
		panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
	}
}