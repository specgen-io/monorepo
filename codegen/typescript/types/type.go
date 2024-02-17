package types

import (
	"fmt"
	"spec"
)

func ResponseBodyTsType(body *spec.ResponseBody) string {
	switch body.Kind() {
	case spec.BodyText:
		return "string"
	case spec.BodyJson:
		return TsType(&body.Type.Definition)
	default:
		panic(fmt.Sprintf("Unsupported response body kind: %v", body.Kind()))
	}
}

func RequestBodyTsType(body *spec.RequestBody) string {
	switch body.Kind() {
	case spec.BodyText:
		return "string"
	case spec.BodyJson:
		return TsType(&body.Type.Definition)
	default:
		panic(fmt.Sprintf("Unsupported request body kind: %v", body.Kind()))
	}
}

func ParamTsType(param *spec.NamedParam) string {
	return TsType(&param.Type.Definition)
}

func TsType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainTsType(typ)
	case spec.NullableType:
		child := TsType(typ.Child)
		result := child + " | undefined"
		return result
	case spec.ArrayType:
		child := TsType(typ.Child)
		result := child + "[]"
		return result
	case spec.MapType:
		child := TsType(typ.Child)
		result := "Record<string, " + child + ">"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainTsType(typ *spec.TypeDef) string {
	switch typ.Plain {
	case spec.TypeInt32:
		return "number"
	case spec.TypeInt64:
		return "number"
	case spec.TypeFloat:
		return "number"
	case spec.TypeDouble:
		return "number"
	case spec.TypeDecimal:
		return "number"
	case spec.TypeBoolean:
		return "boolean"
	case spec.TypeString:
		return "string"
	case spec.TypeUuid:
		return "string"
	case spec.TypeDate:
		return "string"
	case spec.TypeDateTime:
		return "Date"
	case spec.TypeJson:
		return "unknown"
	default:
		if typ.Info.Model != nil {
			if typ.Info.Model.InVersion != nil {
				return fmt.Sprintf("%s.%s", ModelsPackage, typ)
			}
			if typ.Info.Model.InHttpErrors != nil {
				return fmt.Sprintf("%s.%s", ErrorsPackage, typ)
			}
			panic(fmt.Sprintf(`unknown location of type %s`, typ.Plain))
		} else {
			panic(fmt.Sprintf(`unknown type %s`, typ.Plain))
		}
	}
}

var ModelsPackage = "models"
var ErrorsPackage = "errors"
