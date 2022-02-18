package iots

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/ts/common"
	types2 "github.com/specgen-io/specgen/v2/gen/ts/types"
	common2 "github.com/specgen-io/specgen/v2/gen/ts/validations/common"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func (g *Generator) WriteParamsType(w *sources.Writer, typeName string, params []spec.NamedParam) {
	if len(params) > 0 {
		w.EmptyLine()
		w.Line("const %s = t.type({", common2.ParamsRuntimeTypeName(typeName))
		for _, param := range params {
			w.Line("  %s: %s,", common.TSIdentifier(param.Name.Source), paramIoTsTypeDefaulted(&param))
		}
		w.Line("})")
		w.EmptyLine()
		w.Line("type %s = t.TypeOf<typeof %s>", typeName, common2.ParamsRuntimeTypeName(typeName))
	}
}

func paramIoTsTypeDefaulted(param *spec.NamedParam) string {
	theType := paramIoTsType(&param.Type.Definition)
	if param.Default != nil {
		theType = fmt.Sprintf("t.withDefault(%s, %s)", theType, types2.DefaultValue(&param.Type.Definition, *param.Default))
	}
	return theType
}

func paramIoTsType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return paramPlainIoTsType(typ.Plain)
	case spec.NullableType:
		if typ.Child.Node != spec.PlainType {
			panic(fmt.Sprintf("Unsupported string param type: %v", typ))
		}
		return paramPlainIoTsType(typ.Child.Plain)
	case spec.ArrayType:
		if typ.Child.Node != spec.PlainType {
			panic(fmt.Sprintf("Unsupported string param type: %v", typ))
		}
		child := paramPlainIoTsType(typ.Child.Plain)
		result := "t.array(" + child + ")"
		return result
	default:
		panic(fmt.Sprintf("Unsupported string param type: %v", typ))
	}
}

func paramPlainIoTsType(typ string) string {
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
		return fmt.Sprintf("%s.T"+typ, types2.ModelsPackage)
	}
}
