package genruby

import (
	"fmt"
	spec "github.com/specgen-io/spec.v2"
	"github.com/pinzolo/casee"
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
			return `Date.strptime('`+value+`', '%Y-%m-%d')`
		case spec.TypeDateTime:
			return `DateTime.strptime('`+value+`', '%Y-%m-%dT%H:%M:%S')`
		default:
			if typ.Info.Model != nil && typ.Info.Model.IsEnum() {
				return typ.Info.Model.Name.PascalCase() + `::` + casee.ToSnakeCase(value)
			} else {
				panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
			}
		}
	default:
		panic(fmt.Sprintf("Type: %s does not support default value", typ.Name))
	}
}