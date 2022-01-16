package validation

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gents/common"
	"github.com/specgen-io/specgen/v2/gents/types"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func (v *superstructValidation) GenerateParams(w *sources.Writer, typeName string, params []spec.NamedParam) {
	if len(params) > 0 {
		w.EmptyLine()
		w.Line("const %s = t.type({", ParamsRuntimeTypeName(typeName))
		for _, param := range params {
			w.Line("  %s: %s,", common.TSIdentifier(param.Name.Source), paramSuperstructTypeDefaulted(&param))
		}
		w.Line("})")
		w.EmptyLine()
		w.Line("type %s = t.Infer<typeof %s>", typeName, ParamsRuntimeTypeName(typeName))
	}
}

func paramSuperstructTypeDefaulted(param *spec.NamedParam) string {
	theType := paramSuperstructType(&param.Type.Definition)
	if param.Default != nil {
		theType = fmt.Sprintf("t.defaulted(%s, %s)", theType, types.DefaultValue(&param.Type.Definition, *param.Default))
	}
	return theType
}

func paramSuperstructType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return paramPlainSuperstructType(typ.Plain)
	case spec.NullableType:
		if typ.Child.Node != spec.PlainType {
			panic(fmt.Sprintf("Unsupported string type: %v", typ))
		}
		child := paramPlainSuperstructType(typ.Child.Plain)
		result := "t.optional(" + child + ")"
		return result
	case spec.ArrayType:
		if typ.Child.Node != spec.PlainType {
			panic(fmt.Sprintf("Unsupported string type: %v", typ))
		}
		child := paramPlainSuperstructType(typ.Child.Plain)
		result := "t.array(" + child + ")"
		return result
	default:
		panic(fmt.Sprintf("Unsupported string type: %v", typ))
	}
}

func paramPlainSuperstructType(typ string) string {
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
		return fmt.Sprintf("%s.T"+typ, types.ModelsPackage)
	}
}
