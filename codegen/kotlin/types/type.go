package types

import (
	"fmt"
	"spec"
)

const TextType = `String`
const EmptyType = `Unit`

type Types struct {
	RawJsonType string
	BinaryType  BinaryType
	FileType    FileType
}

type BinaryType struct {
	RequestType  string
	ResponseType string
}

type FileType struct {
	RequestType  string
	ResponseType string
}

func NewTypes(rawJsonType, requestBinaryType, responseBinaryType, requestFileType, responseFileType string) *Types {
	return &Types{
		RawJsonType: rawJsonType,
		BinaryType:  BinaryType{RequestType: requestBinaryType, ResponseType: responseBinaryType},
		FileType:    FileType{RequestType: requestFileType, ResponseType: responseFileType},
	}
}

func (t *Types) RequestBodyKotlinType(body *spec.RequestBody) string {
	switch body.Kind() {
	case spec.BodyText:
		return TextType
	case spec.BodyEmpty:
		return EmptyType
	case spec.BodyBinary:
		return t.BinaryType.RequestType
	case spec.BodyJson:
		return t.Kotlin(&body.Type.Definition)
	default:
		panic(fmt.Sprintf("Unknown response body kind: %v", body.Kind()))
	}
}

func (t *Types) ResponseBodyKotlinType(body *spec.ResponseBody) string {
	switch body.Kind() {
	case spec.BodyText:
		return TextType
	case spec.BodyEmpty:
		return EmptyType
	case spec.BodyBinary:
		return t.BinaryType.ResponseType
	case spec.BodyFile:
		return t.FileType.ResponseType
	case spec.BodyJson:
		return t.Kotlin(&body.Type.Definition)
	default:
		panic(fmt.Sprintf("Unknown response body kind: %v", body.Kind()))
	}
}

func (t *Types) ParamKotlinType(param *spec.NamedParam) string {
	if param.Type.Definition.String() == spec.TypeFile {
		return t.FileType.RequestType
	} else {
		return t.Kotlin(&param.Type.Definition)
	}
}

func (t *Types) Kotlin(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return t.plainKotlinType(typ.Plain)
	case spec.NullableType:
		return t.Kotlin(typ.Child) + "?"
	case spec.ArrayType:
		child := t.Kotlin(typ.Child)
		return "List<" + child + ">"
	case spec.MapType:
		child := t.Kotlin(typ.Child)
		result := "Map<String, " + child + ">"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func (t *Types) plainKotlinType(typ string) string {
	switch typ {
	case spec.TypeInt32:
		return "Int"
	case spec.TypeInt64:
		return "Long"
	case spec.TypeFloat:
		return "Float"
	case spec.TypeDouble:
		return "Double"
	case spec.TypeDecimal:
		return "BigDecimal"
	case spec.TypeBoolean:
		return "Boolean"
	case spec.TypeString:
		return "String"
	case spec.TypeUuid:
		return "UUID"
	case spec.TypeDate:
		return "LocalDate"
	case spec.TypeDateTime:
		return "LocalDateTime"
	case spec.TypeJson:
		return t.RawJsonType
	case spec.TypeEmpty:
		return EmptyType
	default:
		return typ
	}
}

func dateFormatAnnotation(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return dateFormatAnnotationPlain(typ.Plain)
	case spec.NullableType:
		return dateFormatAnnotation(typ.Child)
	case spec.ArrayType:
		return dateFormatAnnotation(typ.Child)
	case spec.MapType:
		return dateFormatAnnotation(typ.Child)
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func dateFormatAnnotationPlain(typ string) string {
	switch typ {
	case spec.TypeDate:
		return "@DateTimeFormat(iso = DateTimeFormat.ISO.DATE)"
	case spec.TypeDateTime:
		return "@DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME)"
	default:
		return ""
	}
}

func (t *Types) Imports() []string {
	return []string{
		`java.math.BigDecimal`,
		`java.time.*`,
		`java.util.*`,
		`java.io.*`,
	}
}
