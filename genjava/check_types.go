package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func checkType(fieldType *spec.TypeDef, typ string) bool {
	switch fieldType.Node {
	case spec.PlainType:
		if fieldType.Plain != typ {
			return false
		}
	case spec.NullableType:
		return checkType(fieldType.Child, typ)
	case spec.ArrayType:
		return checkType(fieldType.Child, typ)
	case spec.MapType:
		return checkType(fieldType.Child, typ)
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
	return true
}

func isJavaArrayType(typ *spec.TypeDef) bool {
	switch typ.Node {
	case spec.ArrayType:
		return true
	case spec.NullableType:
		return isJavaArrayType(typ.Child)
	default:
		return false
	}
}

func equalsExpression(typ *spec.TypeDef, left string, right string) string {
	return _equalsExpression(typ, left, right, true)
}

func _equalsExpression(typ *spec.TypeDef, left string, right string, canBePrimitiveType bool) string {
	switch typ.Node {
	case spec.PlainType:
		return equalsExpressionPlainType(typ.Plain, left, right, canBePrimitiveType)
	case spec.NullableType:
		return _equalsExpression(typ.Child, left, right, false)
	case spec.ArrayType:
		return fmt.Sprintf(`Arrays.equals(%s, that.%s)`, left, right)
	case spec.MapType:
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

func toStringExpression(typ *spec.TypeDef, value string) string {
	if isJavaArrayType(typ) {
		return fmt.Sprintf(`Arrays.toString(%s)`, value)
	} else {
		return fmt.Sprintf(`%s`, value)
	}
}
