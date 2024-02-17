package superstruct

import (
	"fmt"
	"spec"
	"typescript/types"
)

func (g *Generator) RuntimeTypeName(typeName string) string {
	return fmt.Sprintf("T%s", typeName)
}

func (g *Generator) RuntimeTypeSamePackage(typ *spec.TypeDef) string {
	return g.runtimeType(typ, true)
}

func (g *Generator) RequestBodyJsonRuntimeType(body *spec.RequestBody) string {
	return g.RuntimeType(&body.Type.Definition)
}

func (g *Generator) ResponseBodyJsonRuntimeType(body *spec.ResponseBody) string {
	return g.RuntimeType(&body.Type.Definition)
}

func (g *Generator) RuntimeType(typ *spec.TypeDef) string {
	return g.runtimeType(typ, false)
}

func (g *Generator) runtimeType(typ *spec.TypeDef, samePackage bool) string {
	switch typ.Node {
	case spec.PlainType:
		return g.plainSuperstructType(typ, samePackage)
	case spec.NullableType:
		child := g.runtimeType(typ.Child, samePackage)
		result := "t.optional(t.nullable(" + child + "))"
		return result
	case spec.ArrayType:
		child := g.runtimeType(typ.Child, samePackage)
		result := "t.array(" + child + ")"
		return result
	case spec.MapType:
		child := g.runtimeType(typ.Child, samePackage)
		result := "t.record(t.string(), " + child + ")"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func (g *Generator) plainSuperstructType(typ *spec.TypeDef, samePackage bool) string {
	switch typ.Plain {
	case spec.TypeInt32:
		return "t.number()"
	case spec.TypeInt64:
		return "t.number()"
	case spec.TypeFloat:
		return "t.number()"
	case spec.TypeDouble:
		return "t.number()"
	case spec.TypeDecimal:
		return "t.number()"
	case spec.TypeBoolean:
		return "t.boolean()"
	case spec.TypeString:
		return "t.string()"
	case spec.TypeUuid:
		return "t.string()"
	case spec.TypeDate:
		return "t.string()"
	case spec.TypeDateTime:
		return "t.StrDateTime"
	case spec.TypeJson:
		return "t.unknown()"
	default:
		if typ.Info.Model != nil {
			if !samePackage {
				if typ.Info.Model.InVersion != nil {
					return fmt.Sprintf(`%s.%s`, types.ModelsPackage, g.RuntimeTypeName(typ.Plain))
				}
				if typ.Info.Model.InHttpErrors != nil {
					return fmt.Sprintf(`%s.%s`, types.ErrorsPackage, g.RuntimeTypeName(typ.Plain))
				}
			}
			return g.RuntimeTypeName(typ.Plain)
		} else {
			panic(fmt.Sprintf(`unknown type %s`, typ.Plain))
		}
	}
}
