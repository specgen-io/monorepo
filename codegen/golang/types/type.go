package types

import (
	"fmt"
	"spec"
)

var VersionModelsPackage = "models"
var ErrorsModelsPackage = "errmodels"

type Types struct{}

func NewTypes() *Types {
	return &Types{}
}

func (types *Types) ResponseBodyType(body *spec.ResponseBody) string {
	if body.IsEmpty() {
		return EmptyType
	} else {
		return types.GoType(&body.Type.Definition)
	}
}

func (types *Types) GoType(typ *spec.TypeDef) string {
	return types.goType(typ, false)
}

func (types *Types) GoTypeSamePackage(typ *spec.TypeDef) string {
	return types.goType(typ, true)
}

func (types *Types) goType(typ *spec.TypeDef, samePackage bool) string {
	switch typ.Node {
	case spec.PlainType:
		return types.plainType(typ, samePackage)
	case spec.NullableType:
		child := types.goType(typ.Child, samePackage)
		if typ.Child.Node == spec.PlainType {
			return "*" + child
		}
		return child
	case spec.ArrayType:
		child := types.goType(typ.Child, samePackage)
		result := "[]" + child
		return result
	case spec.MapType:
		child := types.goType(typ.Child, samePackage)
		result := "map[string]" + child
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func (types *Types) plainType(typ *spec.TypeDef, samePackage bool) string {
	switch typ.Plain {
	case spec.TypeInt32:
		return "int"
	case spec.TypeInt64:
		return "int64"
	case spec.TypeFloat:
		return "float32"
	case spec.TypeDouble:
		return "float64"
	case spec.TypeDecimal:
		return "decimal.Decimal"
	case spec.TypeBoolean:
		return "bool"
	case spec.TypeString:
		return "string"
	case spec.TypeUuid:
		return "uuid.UUID"
	case spec.TypeDate:
		return "civil.Date"
	case spec.TypeDateTime:
		return "civil.DateTime"
	case spec.TypeJson:
		return "json.RawMessage"
	case spec.TypeEmpty:
		return EmptyType
	default:
		if typ.Info.Model != nil {
			if !samePackage {
				if typ.Info.Model.InVersion != nil {
					return fmt.Sprintf("%s.%s", VersionModelsPackage, typ)
				}
				if typ.Info.Model.InHttpErrors != nil {
					return fmt.Sprintf("%s.%s", ErrorsModelsPackage, typ)
				}
			}
			return typ.Plain
		} else {
			panic(fmt.Sprintf(`unknown type %s`, typ.Plain))
		}
	}
}

const EmptyType = `empty.Type`
