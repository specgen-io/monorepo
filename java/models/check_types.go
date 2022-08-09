package models

import (
	"fmt"
	"spec"
)

func equalsExpression(typ *spec.TypeDef, left string, right string) string {
	return _equalsExpression(typ, left, right, true)
}

func _equalsExpression(typ *spec.TypeDef, left string, right string, canBePrimitiveType bool) string {
	switch typ.Node {
	case spec.PlainType:
		return equalsExpressionPlainType(typ.Plain, left, right, canBePrimitiveType)
	case spec.NullableType:
		return _equalsExpression(typ.Child, left, right, false)
	case spec.ArrayType, spec.MapType:
		return fmt.Sprintf(`Objects.equals(%s, that.%s)`, left, right)
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func equalsExpressionPlainType(typ string, left string, right string, canBePrimitiveType bool) string {
	switch typ {
	case spec.TypeInt32,
		spec.TypeInt64,
		spec.TypeBoolean:
		if canBePrimitiveType {
			return fmt.Sprintf(`%s == that.%s`, left, right)
		} else {
			return fmt.Sprintf(`Objects.equals(%s, that.%s)`, left, right)
		}
	case spec.TypeFloat:
		return fmt.Sprintf(`Float.compare(that.%s, %s) == 0`, left, right)
	case spec.TypeDouble:
		return fmt.Sprintf(`Double.compare(that.%s, %s) == 0`, left, right)
	case spec.TypeDecimal,
		spec.TypeString,
		spec.TypeUuid,
		spec.TypeDate,
		spec.TypeDateTime,
		spec.TypeJson:
		return fmt.Sprintf(`Objects.equals(%s, that.%s)`, left, right)
	default:
		return fmt.Sprintf(`Objects.equals(%s, that.%s)`, left, right)
	}
}
