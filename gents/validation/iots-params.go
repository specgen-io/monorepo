package validation

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gents/common"
	"github.com/specgen-io/specgen/v2/gents/types"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func generateIoTsParams(w *sources.Writer, typeName string, params []spec.NamedParam) {
	if len(params) > 0 {
		w.EmptyLine()
		w.Line("const %s = t.type({", ParamsRuntimeTypeName(typeName))
		for _, param := range params {
			w.Line("  %s: %s,", common.TSIdentifier(param.Name.Source), ParamIoTsTypeDefaulted(&param))
		}
		w.Line("})")
		w.EmptyLine()
		w.Line("type %s = t.TypeOf<typeof %s>", typeName, ParamsRuntimeTypeName(typeName))
	}
}

func ParamIoTsTypeDefaulted(param *spec.NamedParam) string {
	theType := ParamIoTsType(&param.Type.Definition)
	if param.Default != nil {
		theType = fmt.Sprintf("t.withDefault(%s, %s)", theType, types.DefaultValue(&param.Type.Definition, *param.Default))
	}
	return theType
}

func ParamIoTsType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return ParamPlainIoTsType(typ.Plain)
	case spec.NullableType:
		if typ.Child.Node != spec.PlainType {
			panic(fmt.Sprintf("Unsupported string param type: %v", typ))
		}
		return ParamPlainIoTsType(typ.Child.Plain)
	case spec.ArrayType:
		if typ.Child.Node != spec.PlainType {
			panic(fmt.Sprintf("Unsupported string param type: %v", typ))
		}
		child := ParamPlainIoTsType(typ.Child.Plain)
		result := "t.array(" + child + ")"
		return result
	default:
		panic(fmt.Sprintf("Unsupported string param type: %v", typ))
	}
}

func ParamPlainIoTsType(typ string) string {
	switch typ {
	case spec.TypeInt32:
		return "t.IntFromString"
	case spec.TypeInt64:
		return "t.IntFromString"
	case spec.TypeFloat:
		return "t.NumberFromString"
	case spec.TypeDouble:
		return "t.NumberFromString"
	case spec.TypeDecimal:
		return "t.NumberFromString"
	case spec.TypeBoolean:
		return "t.BooleanFromString"
	case spec.TypeString:
		return "t.string"
	case spec.TypeUuid:
		return "t.string"
	case spec.TypeDate:
		return "t.string"
	case spec.TypeDateTime:
		return "t.DateISOStringNoTimezone"
	default:
		return fmt.Sprintf("%s.T"+typ, types.ModelsPackage)
	}
}
