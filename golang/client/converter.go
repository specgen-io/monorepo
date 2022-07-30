package client

import (
	"fmt"
	"strings"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/golang/v2/module"
	"github.com/specgen-io/specgen/spec/v2"
)

func callRawConvert(typ *spec.TypeDef, paramNameVar string) string {
	return fmt.Sprintf("convert%s(%s)", converterMethodName(typ), paramNameVar)
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

func generateConverter(module module.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

import (
	"cloud.google.com/go/civil"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"strconv"
	"test-client/spec/models"
)

func convertInt(value int) string {
	return strconv.Itoa(value)
}

func convertInt64(value int64) string {
	return strconv.FormatInt(value, 10)
}

func convertFloat32(value float32) string {
	return strconv.FormatFloat(float64(value), 'f', -1, 32)
}

func convertFloat64(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func convertDecimal(value decimal.Decimal) string {
	return value.String()
}

func convertBool(value bool) string {
	return strconv.FormatBool(value)
}

func convertUuid(value uuid.UUID) string {
	return value.String()
}

func convertDate(value civil.Date) string {
	return value.String()
}

func convertDateTime(value civil.DateTime) string {
	return value.String()
}

func convertStringEnum(value models.Choice) string {
	return string(value)
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
	self.parser.Set(key, convertInt(value))
}

func (self *ParamsConverter) IntNullable(key string, value *int) {
	self.parser.Set(key, convertInt(*value))
}

func (self *ParamsConverter) IntArray(key string, values []int) {
	for _, value := range values {
		self.parser.Add(key, convertInt(value))
	}
}

func (self *ParamsConverter) Int64(key string, value int64) {
	self.parser.Set(key, convertInt64(value))
}

func (self *ParamsConverter) Int64Nullable(key string, value *int64) {
	self.parser.Set(key, convertInt64(*value))
}

func (self *ParamsConverter) Int64Array(key string, values []int64) {
	for _, value := range values {
		self.parser.Add(key, convertInt64(value))
	}
}

func (self *ParamsConverter) Float32(key string, value float32) {
	self.parser.Set(key, convertFloat32(value))
}

func (self *ParamsConverter) Float32Nullable(key string, value *float32) {
	self.parser.Set(key, convertFloat32(*value))
}

func (self *ParamsConverter) Float32Array(key string, values []float32) {
	for _, value := range values {
		self.parser.Add(key, convertFloat32(value))
	}
}

func (self *ParamsConverter) Float64(key string, value float64) {
	self.parser.Set(key, convertFloat64(value))
}

func (self *ParamsConverter) Float64Nullable(key string, value *float64) {
	self.parser.Set(key, convertFloat64(*value))
}

func (self *ParamsConverter) Float64Array(key string, values []float64) {
	for _, value := range values {
		self.parser.Add(key, convertFloat64(value))
	}
}

func (self *ParamsConverter) Decimal(key string, value decimal.Decimal) {
	self.parser.Set(key, convertDecimal(value))
}

func (self *ParamsConverter) DecimalNullable(key string, value *decimal.Decimal) {
	self.parser.Set(key, convertDecimal(*value))
}

func (self *ParamsConverter) DecimalArray(key string, values []decimal.Decimal) {
	for _, value := range values {
		self.parser.Add(key, convertDecimal(value))
	}
}

func (self *ParamsConverter) Bool(key string, value bool) {
	self.parser.Set(key, convertBool(value))
}

func (self *ParamsConverter) BoolNullable(key string, value *bool) {
	self.parser.Set(key, convertBool(*value))
}

func (self *ParamsConverter) BoolArray(key string, values []bool) {
	for _, value := range values {
		self.parser.Add(key, convertBool(value))
	}
}

func (self *ParamsConverter) Uuid(key string, value uuid.UUID) {
	self.parser.Set(key, convertUuid(value))
}

func (self *ParamsConverter) UuidNullable(key string, value *uuid.UUID) {
	self.parser.Set(key, convertUuid(*value))
}

func (self *ParamsConverter) UuidArray(key string, values []uuid.UUID) {
	for _, value := range values {
		self.parser.Add(key, convertUuid(value))
	}
}

func (self *ParamsConverter) Date(key string, value civil.Date) {
	self.parser.Set(key, convertDate(value))
}

func (self *ParamsConverter) DateNullable(key string, value *civil.Date) {
	self.parser.Set(key, convertDate(*value))
}

func (self *ParamsConverter) DateArray(key string, values []civil.Date) {
	for _, value := range values {
		self.parser.Add(key, convertDate(value))
	}
}

func (self *ParamsConverter) DateTime(key string, value civil.DateTime) {
	self.parser.Set(key, convertDateTime(value))
}

func (self *ParamsConverter) DateTimeNullable(key string, value *civil.DateTime) {
	self.parser.Set(key, convertDateTime(*value))
}

func (self *ParamsConverter) DateTimeArray(key string, values []civil.DateTime) {
	for _, value := range values {
		self.parser.Add(key, convertDateTime(value))
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
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{module.Name})
	return &generator.CodeFile{module.GetPath("converter.go"), strings.TrimSpace(code)}
}
