package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func generateIoTsParams(w *gen.Writer, typeName string, isHeader bool, params []spec.NamedParam) {
	if len(params) > 0 {
		w.EmptyLine()
		w.Line("const %s = t.type({", paramsRuntimeTypeName(typeName))
		for _, param := range params {
			paramName := param.Name.Source
			if isHeader {
				paramName = strings.ToLower(param.Name.Source)
			}
			paramName = tsIdentifier(paramName)
			w.Line("  %s: %s,", paramName, ParamIoTsTypeDefaulted(&param))
		}
		w.Line("})")
		w.EmptyLine()
		w.Line("type %s = t.TypeOf<typeof %s>", typeName, paramsRuntimeTypeName(typeName))
	}
}

func ParamIoTsTypeDefaulted(param *spec.NamedParam) string {
	theType := ParamIoTsType(&param.Type.Definition)
	if param.Default != nil {
		theType = fmt.Sprintf("t.withDefault(%s, %s)", theType, DefaultValue(&param.Type.Definition, *param.Default))
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
		return "t.DateFromISOString"
	default:
		return fmt.Sprintf("%s.T"+typ, modelsPackage)
	}
}
