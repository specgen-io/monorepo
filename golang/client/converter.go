package client

import (
	"fmt"
	"generator"
	"golang/writer"
	"spec"
)

func callRawConvert(typ *spec.TypeDef, paramNameVar string) string {
	return fmt.Sprintf("convert.%s(%s)", converterMethodName(typ), paramNameVar)
}

func callConverter(typ *spec.TypeDef, paramName string, paramNameVar string) string {
	return fmt.Sprintf(`%s("%s", %s)`, converterMethodName(typ), paramName, paramNameVar)
}

func converterMethodName(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return converterMethodNamePlain(typ)
	case spec.NullableType:
		return converterMethodNamePlain(typ.Child) + "Nullable"
	case spec.ArrayType:
		return converterMethodNamePlain(typ.Child) + "Array"
	default:
		panic(fmt.Sprintf("Unsupported string param type: %v", typ.Plain))
	}
}

func converterMethodNamePlain(typ *spec.TypeDef) string {
	if typ.Info.Model != nil && typ.Info.Model.IsEnum() {
		return "StringEnum"
	}
	switch typ.Plain {
	case spec.TypeInt32:
		return "Int"
	case spec.TypeInt64:
		return "Int64"
	case spec.TypeFloat:
		return "Float32"
	case spec.TypeDouble:
		return "Float64"
	case spec.TypeDecimal:
		return "Decimal"
	case spec.TypeBoolean:
		return "Bool"
	case spec.TypeString:
		return "String"
	case spec.TypeUuid:
		return "Uuid"
	case spec.TypeDate:
		return "Date"
	case spec.TypeDateTime:
		return "DateTime"
	default:
		panic(fmt.Sprintf("Unsupported string param type: %v", typ.Plain))
	}
}

func (g *Generator) GenerateConverter() *generator.CodeFile {
	w := writer.New(g.Modules.Convert, `convert.go`)
	w.Lines(`
import (
	"cloud.google.com/go/civil"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"strconv"
)

func Int(value int) string {
	return strconv.Itoa(value)
}

func Int64(value int64) string {
	return strconv.FormatInt(value, 10)
}

func Float32(value float32) string {
	return strconv.FormatFloat(float64(value), 'f', -1, 32)
}

func Float64(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func Decimal(value decimal.Decimal) string {
	return value.String()
}

func Bool(value bool) string {
	return strconv.FormatBool(value)
}

func Uuid(value uuid.UUID) string {
	return value.String()
}

func Date(value civil.Date) string {
	return value.String()
}

func DateTime(value civil.DateTime) string {
	return value.String()
}

type ParamsSetter interface {
	Add(key string, value string)
	Set(key string, value string)
}

type ParamsConverter struct {
	parser ParamsSetter
}

func NewParamsConverter(parser ParamsSetter) *ParamsConverter {
	return &ParamsConverter{parser}
}

func (self *ParamsConverter) String(key string, value string) {
	self.parser.Set(key, value)
}

func (self *ParamsConverter) StringNullable(key string, value *string) {
	self.parser.Set(key, *value)
}

func (self *ParamsConverter) StringArray(key string, values []string) {
	for _, value := range values {
		self.parser.Add(key, value)
	}
}

func (self *ParamsConverter) Int(key string, value int) {
	self.parser.Set(key, Int(value))
}

func (self *ParamsConverter) IntNullable(key string, value *int) {
	self.parser.Set(key, Int(*value))
}

func (self *ParamsConverter) IntArray(key string, values []int) {
	for _, value := range values {
		self.parser.Add(key, Int(value))
	}
}

func (self *ParamsConverter) Int64(key string, value int64) {
	self.parser.Set(key, Int64(value))
}

func (self *ParamsConverter) Int64Nullable(key string, value *int64) {
	self.parser.Set(key, Int64(*value))
}

func (self *ParamsConverter) Int64Array(key string, values []int64) {
	for _, value := range values {
		self.parser.Add(key, Int64(value))
	}
}

func (self *ParamsConverter) Float32(key string, value float32) {
	self.parser.Set(key, Float32(value))
}

func (self *ParamsConverter) Float32Nullable(key string, value *float32) {
	self.parser.Set(key, Float32(*value))
}

func (self *ParamsConverter) Float32Array(key string, values []float32) {
	for _, value := range values {
		self.parser.Add(key, Float32(value))
	}
}

func (self *ParamsConverter) Float64(key string, value float64) {
	self.parser.Set(key, Float64(value))
}

func (self *ParamsConverter) Float64Nullable(key string, value *float64) {
	self.parser.Set(key, Float64(*value))
}

func (self *ParamsConverter) Float64Array(key string, values []float64) {
	for _, value := range values {
		self.parser.Add(key, Float64(value))
	}
}

func (self *ParamsConverter) Decimal(key string, value decimal.Decimal) {
	self.parser.Set(key, Decimal(value))
}

func (self *ParamsConverter) DecimalNullable(key string, value *decimal.Decimal) {
	self.parser.Set(key, Decimal(*value))
}

func (self *ParamsConverter) DecimalArray(key string, values []decimal.Decimal) {
	for _, value := range values {
		self.parser.Add(key, Decimal(value))
	}
}

func (self *ParamsConverter) Bool(key string, value bool) {
	self.parser.Set(key, Bool(value))
}

func (self *ParamsConverter) BoolNullable(key string, value *bool) {
	self.parser.Set(key, Bool(*value))
}

func (self *ParamsConverter) BoolArray(key string, values []bool) {
	for _, value := range values {
		self.parser.Add(key, Bool(value))
	}
}

func (self *ParamsConverter) Uuid(key string, value uuid.UUID) {
	self.parser.Set(key, Uuid(value))
}

func (self *ParamsConverter) UuidNullable(key string, value *uuid.UUID) {
	self.parser.Set(key, Uuid(*value))
}

func (self *ParamsConverter) UuidArray(key string, values []uuid.UUID) {
	for _, value := range values {
		self.parser.Add(key, Uuid(value))
	}
}

func (self *ParamsConverter) Date(key string, value civil.Date) {
	self.parser.Set(key, Date(value))
}

func (self *ParamsConverter) DateNullable(key string, value *civil.Date) {
	self.parser.Set(key, Date(*value))
}

func (self *ParamsConverter) DateArray(key string, values []civil.Date) {
	for _, value := range values {
		self.parser.Add(key, Date(value))
	}
}

func (self *ParamsConverter) DateTime(key string, value civil.DateTime) {
	self.parser.Set(key, DateTime(value))
}

func (self *ParamsConverter) DateTimeNullable(key string, value *civil.DateTime) {
	self.parser.Set(key, DateTime(*value))
}

func (self *ParamsConverter) DateTimeArray(key string, values []civil.DateTime) {
	for _, value := range values {
		self.parser.Add(key, DateTime(value))
	}
}

func (self *ParamsConverter) StringEnum(key string, value interface{}) {
	self.parser.Set(key, fmt.Sprintf("%v", value))
}

func (self *ParamsConverter) StringEnumNullable(key string, value *interface{}) {
	self.parser.Set(key, fmt.Sprintf("%v", *value))
}

func (self *ParamsConverter) StringEnumArray(key string, values []interface{}) {
	for _, value := range values {
		self.parser.Add(key, fmt.Sprintf("%v", value))
	}
}
`)
	return w.ToCodeFile()
}
