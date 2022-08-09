package iots

import (
	"fmt"
	"spec"
)

func (g *Generator) RuntimeType(typ *spec.TypeDef) string {
	return g.RuntimeTypeFromPackage("", typ)
}

func (g *Generator) RuntimeTypeFromPackage(customTypesPackage string, typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return g.plainIoTsType(customTypesPackage, typ.Plain)
	case spec.NullableType:
		child := g.RuntimeTypeFromPackage(customTypesPackage, typ.Child)
		result := "t.union([" + child + ", t.null])"
		return result
	case spec.ArrayType:
		child := g.RuntimeTypeFromPackage(customTypesPackage, typ.Child)
		result := "t.array(" + child + ")"
		return result
	case spec.MapType:
		child := g.RuntimeTypeFromPackage(customTypesPackage, typ.Child)
		result := "t.record(t.string, " + child + ")"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func (g *Generator) plainIoTsType(customTypesPackage string, typ string) string {
	switch typ {
	case spec.TypeInt32:
		return "t.number"
	case spec.TypeInt64:
		return "t.number"
	case spec.TypeFloat:
		return "t.number"
	case spec.TypeDouble:
		return "t.number"
	case spec.TypeDecimal:
		return "t.number"
	case spec.TypeBoolean:
		return "t.boolean"
	case spec.TypeString:
		return "t.string"
	case spec.TypeUuid:
		return "t.string"
	case spec.TypeDate:
		return "t.string"
	case spec.TypeDateTime:
		return "t.DateISOStringNoTimezone"
	case spec.TypeJson:
		return "t.unknown"
	default:
		if customTypesPackage == "" {
			return "T" + typ
		} else {
			return customTypesPackage + ".T" + typ
		}
	}
}
