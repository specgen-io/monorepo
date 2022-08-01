package generators

import (
	"fmt"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/spec/v2"
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
			spec.TypeInt64:
			return value
		case
			spec.TypeDouble,
			spec.TypeFloat,
			spec.TypeDecimal:
			return value + ".to_f"
		case spec.TypeBoolean:
			return value
		case
			spec.TypeString,
			spec.TypeUuid:
			return `'` + value + `'`
		case spec.TypeDate:
			return `Date.strptime('` + value + `', '%Y-%m-%d')`
		case spec.TypeDateTime:
			return `DateTime.strptime('` + value + `', '%Y-%m-%dT%H:%M:%S')`
		default:
			model := typ.Info.Model
			if model != nil && model.IsEnum() {
				return model.Name.PascalCase() + `::` + casee.ToSnakeCase(value)
			} else {
				panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
			}
		}
	default:
		panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
	}
}
