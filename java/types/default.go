package types

import (
	"fmt"
	"github.com/pinzolo/casee"
	"spec"
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
			return "Map<>"
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
			return `LocalDate.parse("` + value + `")`
		case spec.TypeDateTime:
			return `LocalDateTime.parse("` + value + `")`
		default:
			if typ.Info.Model != nil && typ.Info.Model.IsEnum() {
				return typ.Info.Model.Name.PascalCase() + `.` + casee.ToUpperCase(value)
			} else {
				panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
			}
		}
	default:
		panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
	}
}
