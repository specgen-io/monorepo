package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateSuperstructParams(w *sources.Writer, typeName string, params []spec.NamedParam) {
	if len(params) > 0 {
		w.EmptyLine()
		w.Line("const %s = t.type({", paramsRuntimeTypeName(typeName))
		for _, param := range params {
			w.Line("  %s: %s,", tsIdentifier(param.Name.Source), ParamSuperstructTypeDefaulted(&param))
		}
		w.Line("})")
		w.EmptyLine()
		w.Line("type %s = t.Infer<typeof %s>", typeName, paramsRuntimeTypeName(typeName))
	}
}

func ParamSuperstructTypeDefaulted(param *spec.NamedParam) string {
	theType := ParamSuperstructType(&param.Type.Definition)
	if param.Default != nil {
		theType = fmt.Sprintf("t.defaulted(%s, %s)", theType, DefaultValue(&param.Type.Definition, *param.Default))
	}
	return theType
}

func ParamSuperstructType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return ParamPlainSuperstructType(typ.Plain)
	case spec.NullableType:
		if typ.Child.Node != spec.PlainType {
			panic(fmt.Sprintf("Unsupported string type: %v", typ))
		}
		child := ParamPlainSuperstructType(typ.Child.Plain)
		result := "t.optional(" + child + ")"
		return result
	case spec.ArrayType:
		if typ.Child.Node != spec.PlainType {
			panic(fmt.Sprintf("Unsupported string type: %v", typ))
		}
		child := ParamPlainSuperstructType(typ.Child.Plain)
		result := "t.array(" + child + ")"
		return result
	default:
		panic(fmt.Sprintf("Unsupported string type: %v", typ))
	}
}

func ParamPlainSuperstructType(typ string) string {
	switch typ {
	case spec.TypeInt32:
		return "t.StrInteger"
	case spec.TypeInt64:
		return "t.StrInteger"
	case spec.TypeFloat:
		return "t.StrFloat"
	case spec.TypeDouble:
		return "t.StrFloat"
	case spec.TypeDecimal:
		return "t.StrFloat"
	case spec.TypeBoolean:
		return "t.StrBoolean"
	case spec.TypeString:
		return "t.string()"
	case spec.TypeUuid:
		return "t.string()"
	case spec.TypeDate:
		return "t.string()"
	case spec.TypeDateTime:
		return "t.StrDateTime"
	default:
		return fmt.Sprintf("%s.T"+typ, modelsPackage)
	}
}
