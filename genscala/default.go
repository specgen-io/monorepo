package genscala

import (
	"fmt"
	"github.com/ModaOperandi/spec"
	"github.com/vsapronov/casee"
)

func DefaultValue(typ *spec.TypeDef, value string) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainScalaValue(typ.Plain, typ.Info.Model, value)
	default:
		panic(fmt.Sprintf("Type: %v does not support default value", typ))
	}
}

func PlainScalaValue(typ string, model *spec.NamedModel, value string) string {
	switch typ {
	case
		spec.TypeByte,
		spec.TypeInt16,
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
	case spec.TypeChar:
		return "'" + value + "'"
	case spec.TypeString:
		return `"` + value + `"`
	case spec.TypeUuid:
		return `UUID.fromString("` + value + `")`
	case spec.TypeDate:
		return `LocalDate.parse("` + value + `", DateTimeFormatter.ISO_LOCAL_DATE)`
	case spec.TypeDateTime:
		return `LocalDateTime.parse("` + value + `", DateTimeFormatter.ISO_LOCAL_DATE_TIME)`
	case spec.TypeTime:
		return `LocalTime.parse("` + value + `", DateTimeFormatter.ISO_LOCAL_TIME)`
	default:
		if model != nil && model.IsEnum() {
			return model.Name.PascalCase() + `.` + casee.ToPascalCase(value)
		} else {
			panic(fmt.Sprintf("Type: %v does not support default value", typ))
		}
	}
}
