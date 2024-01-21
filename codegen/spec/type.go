package spec

import (
	"errors"
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
	"strings"
)

type TypeNode int

const (
	PlainType    TypeNode = 0
	NullableType TypeNode = 1
	ArrayType    TypeNode = 2
	MapType      TypeNode = 3
)

type TypeDef struct {
	Name  string
	Node  TypeNode
	Child *TypeDef
	Plain string
	Info  *TypeInfo
}

func Plain(typ string) *TypeDef {
	return &TypeDef{Name: typ, Node: PlainType, Plain: typ}
}

func Array(typ *TypeDef) *TypeDef {
	return &TypeDef{Name: typ.Name + "[]", Node: ArrayType, Child: typ}
}

func Map(typ *TypeDef) *TypeDef {
	return &TypeDef{Name: typ.Name + "{}", Node: MapType, Child: typ}
}

func Nullable(typ *TypeDef) *TypeDef {
	return &TypeDef{Name: typ.Name + "?", Node: NullableType, Child: typ}
}

func (self *TypeDef) IsEmpty() bool {
	return self.Node == PlainType && self.Plain == TypeEmpty
}

func (self *TypeDef) IsNullable() bool {
	return self.Node == NullableType
}

func (self *TypeDef) IsModel() bool {
	return self.Info.Model != nil
}

func (self *TypeDef) IsEnum() bool {
	return self.Info.Model != nil && self.Info.Model.IsEnum()
}

func (self *TypeDef) IsErrorModel() bool {
	return self.Info.Model != nil && self.Info.Model.InHttpErrors != nil
}

func (self *TypeDef) BaseType() *TypeDef {
	if self.IsNullable() {
		return self.Child
	}
	return self
}

var plainTypeFormat = FormatOr(PascalCase, LowerCase)

func parseType(value string) (*TypeDef, error) {
	if strings.HasSuffix(value, "?") {
		child, err := parseType(value[:len(value)-1])
		if err != nil {
			return nil, err
		}
		return &TypeDef{Name: value, Node: NullableType, Child: child}, nil
	} else if strings.HasSuffix(value, "[]") {
		child, err := parseType(value[:len(value)-2])
		if err != nil {
			return nil, err
		}
		return &TypeDef{Name: value, Node: ArrayType, Child: child}, nil
	} else if strings.HasSuffix(value, "{}") {
		child, err := parseType(value[:len(value)-2])
		if err != nil {
			return nil, err
		}
		return &TypeDef{Name: value, Node: MapType, Child: child}, nil
	} else {
		err := plainTypeFormat.Check(value)
		if err != nil {
			return nil, errors.New("type " + err.Error())
		}
		return &TypeDef{Name: value, Node: PlainType, Plain: mapTypeAlias(value)}, nil
	}
}

type Type struct {
	Definition TypeDef
	Location   *yaml.Node
}

func NewType(str string) *Type {
	typ, err := parseType(str)
	if err != nil {
		panic(err)
	}
	return &Type{*typ, nil}
}

func (typ Type) String() string {
	return typ.Definition.String()
}

func (value *Type) UnmarshalYAML(node *yaml.Node) error {
	str := ""
	err := node.DecodeWith(decodeStrict, &str)
	if err != nil {
		return err
	}
	typ, err := parseType(str)
	if err != nil {
		return yamlError(node, err.Error())
	}
	*value = Type{*typ, node}
	return nil
}

const (
	TypeInt32    string = "int32"
	TypeInt64    string = "int64"
	TypeFloat    string = "float"
	TypeDouble   string = "double"
	TypeDecimal  string = "decimal"
	TypeBoolean  string = "boolean"
	TypeString   string = "string"
	TypeUuid     string = "uuid"
	TypeDate     string = "date"
	TypeDateTime string = "datetime"
	TypeJson     string = "json"
	TypeFile     string = "file"
	TypeEmpty    string = "empty"
)

const (
	TypeAliasInt  string = "int"
	TypeAliasLong string = "long"
	TypeAliasBool string = "bool"
	TypeAliasStr  string = "str"
)

var TypesAliases = map[string]string{
	TypeAliasInt:  TypeInt32,
	TypeAliasLong: TypeInt64,
	TypeAliasBool: TypeBoolean,
	TypeAliasStr:  TypeString,
}

func mapTypeAlias(value string) string {
	if mapped, ok := TypesAliases[value]; ok {
		return mapped
	}
	return value
}

type TypeStructure int

const (
	StructureNone   TypeStructure = 0
	StructureScalar TypeStructure = 1
	StructureArray  TypeStructure = 2
	StructureObject TypeStructure = 3
)

type TypeInfo struct {
	Structure   TypeStructure
	Defaultable bool
	Model       *NamedModel
}

var Types = map[string]TypeInfo{
	TypeInt32:    {StructureScalar, true, nil},
	TypeInt64:    {StructureScalar, true, nil},
	TypeFloat:    {StructureScalar, true, nil},
	TypeDouble:   {StructureScalar, true, nil},
	TypeDecimal:  {StructureScalar, true, nil},
	TypeBoolean:  {StructureScalar, true, nil},
	TypeString:   {StructureScalar, true, nil},
	TypeUuid:     {StructureScalar, true, nil},
	TypeDate:     {StructureScalar, true, nil},
	TypeDateTime: {StructureScalar, true, nil},
	TypeJson:     {StructureObject, false, nil},
	TypeFile:     {StructureScalar, false, nil},
	TypeEmpty:    {StructureNone, false, nil},
}

func ModelTypeInfo(model *NamedModel) *TypeInfo {
	if model.IsObject() || model.IsOneOf() {
		return &TypeInfo{StructureObject, false, model}
	}
	if model.IsEnum() {
		return &TypeInfo{StructureScalar, true, model}
	}
	panic(fmt.Sprintf("Unknown model kind: %v", model))
}

func NullableTypeInfo(childInfo *TypeInfo) *TypeInfo {
	if childInfo != nil {
		return &TypeInfo{childInfo.Structure, false, nil}
	}
	return nil
}

func ArrayTypeInfo() *TypeInfo {
	return &TypeInfo{StructureArray, true, nil}
}

func MapTypeInfo() *TypeInfo {
	return &TypeInfo{StructureObject, true, nil}
}

func ParseType(value string) TypeDef {
	typ, err := parseType(value)
	if err != nil {
		panic(err.Error())
	}
	return *typ
}

func ParseEndpoint(endpointStr string) *Endpoint {
	endpoint, err := parseEndpoint(endpointStr, nil)
	if err != nil {
		panic(err.Error())
	}
	return endpoint
}

func (typ TypeDef) String() string {
	switch typ.Node {
	case PlainType:
		switch typ.Plain {
		case TypeInt32:
			return TypeAliasInt
		case TypeInt64:
			return TypeAliasLong
		default:
			return typ.Plain
		}
		return typ.Plain
	case NullableType:
		child := typ.Child.String()
		result := child + "?"
		return result
	case ArrayType:
		child := typ.Child.String()
		result := child + "[]"
		return result
	case MapType:
		child := typ.Child.String()
		result := child + "{}"
		return result
	default:
		panic(fmt.Sprintf("Unknown type node: %d", typ.Node))
	}
}
