package client

import (
	"fmt"
	"generator"
	"golang/writer"
	"spec"
)

func callRawParamsConvert(typ *spec.TypeDef, paramNameVar string) string {
	return fmt.Sprintf("params.%s(%s)", converterMethodName(typ), paramNameVar)
}

func callTypesConverter(typ *spec.TypeDef, paramName string, paramNameVar string) string {
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
	case spec.TypeFile:
		return "CreateFile"
	default:
		panic(fmt.Sprintf("Unsupported string param type: %v", typ.Plain))
	}
}

func (g *Generator) TypeConverter() *generator.CodeFile {
	w := writer.New(g.Modules.Params, `convert_types.go`)
	w.Lines(`
import (
	"cloud.google.com/go/civil"
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
`)
	return w.ToCodeFile()
}

func (g *Generator) Params() *generator.CodeFile {
	w := writer.New(g.Modules.Params, `params.go`)
	w.Lines(`
import (
	"cloud.google.com/go/civil"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ParamsSetter interface {
	Add(key, value string)
}

type ParamsWriter struct {
	setter ParamsSetter
}

func NewParamsWriter(setter ParamsSetter) *ParamsWriter {
	return &ParamsWriter{setter}
}

func (self *ParamsWriter) String(key string, value string) {
	self.setter.Add(key, value)
}

func (self *ParamsWriter) StringNullable(key string, value *string) {
	self.setter.Add(key, *value)
}

func (self *ParamsWriter) StringArray(key string, values []string) {
	for _, value := range values {
		self.setter.Add(key, value)
	}
}

func (self *ParamsWriter) Int(key string, value int) {
	self.setter.Add(key, Int(value))
}

func (self *ParamsWriter) IntNullable(key string, value *int) {
	self.setter.Add(key, Int(*value))
}

func (self *ParamsWriter) IntArray(key string, values []int) {
	for _, value := range values {
		self.setter.Add(key, Int(value))
	}
}

func (self *ParamsWriter) Int64(key string, value int64) {
	self.setter.Add(key, Int64(value))
}

func (self *ParamsWriter) Int64Nullable(key string, value *int64) {
	self.setter.Add(key, Int64(*value))
}

func (self *ParamsWriter) Int64Array(key string, values []int64) {
	for _, value := range values {
		self.setter.Add(key, Int64(value))
	}
}

func (self *ParamsWriter) Float32(key string, value float32) {
	self.setter.Add(key, Float32(value))
}

func (self *ParamsWriter) Float32Nullable(key string, value *float32) {
	self.setter.Add(key, Float32(*value))
}

func (self *ParamsWriter) Float32Array(key string, values []float32) {
	for _, value := range values {
		self.setter.Add(key, Float32(value))
	}
}

func (self *ParamsWriter) Float64(key string, value float64) {
	self.setter.Add(key, Float64(value))
}

func (self *ParamsWriter) Float64Nullable(key string, value *float64) {
	self.setter.Add(key, Float64(*value))
}

func (self *ParamsWriter) Float64Array(key string, values []float64) {
	for _, value := range values {
		self.setter.Add(key, Float64(value))
	}
}

func (self *ParamsWriter) Decimal(key string, value decimal.Decimal) {
	self.setter.Add(key, Decimal(value))
}

func (self *ParamsWriter) DecimalNullable(key string, value *decimal.Decimal) {
	self.setter.Add(key, Decimal(*value))
}

func (self *ParamsWriter) DecimalArray(key string, values []decimal.Decimal) {
	for _, value := range values {
		self.setter.Add(key, Decimal(value))
	}
}

func (self *ParamsWriter) Bool(key string, value bool) {
	self.setter.Add(key, Bool(value))
}

func (self *ParamsWriter) BoolNullable(key string, value *bool) {
	self.setter.Add(key, Bool(*value))
}

func (self *ParamsWriter) BoolArray(key string, values []bool) {
	for _, value := range values {
		self.setter.Add(key, Bool(value))
	}
}

func (self *ParamsWriter) Uuid(key string, value uuid.UUID) {
	self.setter.Add(key, Uuid(value))
}

func (self *ParamsWriter) UuidNullable(key string, value *uuid.UUID) {
	self.setter.Add(key, Uuid(*value))
}

func (self *ParamsWriter) UuidArray(key string, values []uuid.UUID) {
	for _, value := range values {
		self.setter.Add(key, Uuid(value))
	}
}

func (self *ParamsWriter) Date(key string, value civil.Date) {
	self.setter.Add(key, Date(value))
}

func (self *ParamsWriter) DateNullable(key string, value *civil.Date) {
	self.setter.Add(key, Date(*value))
}

func (self *ParamsWriter) DateArray(key string, values []civil.Date) {
	for _, value := range values {
		self.setter.Add(key, Date(value))
	}
}

func (self *ParamsWriter) DateTime(key string, value civil.DateTime) {
	self.setter.Add(key, DateTime(value))
}

func (self *ParamsWriter) DateTimeNullable(key string, value *civil.DateTime) {
	self.setter.Add(key, DateTime(*value))
}

func (self *ParamsWriter) DateTimeArray(key string, values []civil.DateTime) {
	for _, value := range values {
		self.setter.Add(key, DateTime(value))
	}
}

func (self *ParamsWriter) StringEnum(key string, value interface{}) {
	self.setter.Add(key, fmt.Sprintf("%v", value))
}

func (self *ParamsWriter) StringEnumNullable(key string, value *interface{}) {
	self.setter.Add(key, fmt.Sprintf("%v", *value))
}

func (self *ParamsWriter) StringEnumArray(key string, values []interface{}) {
	for _, value := range values {
		self.setter.Add(key, fmt.Sprintf("%v", value))
	}
}
`)
	return w.ToCodeFile()
}

func (g *Generator) FormDataParams() *generator.CodeFile {
	w := writer.New(g.Modules.Params, `form_data_params.go`)

	w.Template(
		map[string]string{
			`HttpFilePackage`: g.Modules.HttpFile.Package,
		}, `
import (
	"errors"
	"fmt"
	"[[.HttpFilePackage]]"
	"io"
	"path/filepath"
)

type FormDataParamsSetter interface {
	WriteField(key, value string) error
	CreateFormFile(fieldname, filename string) (io.Writer, error)
	Close() error
}

type FormDataParamsWriter struct {
	ParamsWriter
	bridgeSetter
}

func NewFormDataParamsWriter(setter FormDataParamsSetter) *FormDataParamsWriter {
	bridgeSetter := bridgeSetter{setter, []error{}}
	return &FormDataParamsWriter{*NewParamsWriter(&bridgeSetter), bridgeSetter}
}

type bridgeSetter struct {
	FormDataParamsSetter
	errors []error
}

func (self *bridgeSetter) Add(key, value string) {
	err := self.WriteField(key, value)
	if err != nil {
		self.errors = append(self.errors, errors.New(fmt.Sprintf("failed to set parameter %s: %s ", key, err.Error())))
	}
}

func (self *FormDataParamsWriter) CreateFile(fieldName string, file *httpfile.File) {
	dst, err := self.CreateFormFile(fieldName, filepath.Base(file.Name))
	_, err = io.Copy(dst, file.Content)
	if err != nil {
		self.errors = append(self.errors, errors.New(fmt.Sprintf("failed to write file: %s", err.Error())))
	}
}

func (self *FormDataParamsWriter) CloseWriter() error {
	err := self.Close()
	if err != nil {
		self.errors = append(self.errors, errors.New(fmt.Sprintf("failed to write form parameters: %s", err.Error())))
	}
	if len(self.errors) > 0 {
		message := "errors: "
		for _, err := range self.errors {
			message += err.Error() + ", "
		}
		return errors.New(message)
	}
	return nil
}
`)
	return w.ToCodeFile()
}
