package service

import (
	"fmt"
	"strings"

	"generator"
	"golang/module"
	"golang/types"
	"spec"
)

func parserDefaultName(param *spec.NamedParam) (string, *string) {
	methodName := parserMethodName(&param.Type.Definition)
	if param.Default != nil {
		defaultValue := types.DefaultValue(&param.Type.Definition, *param.Default)
		return methodName + `Defaulted`, &defaultValue
	} else {
		return methodName, nil
	}
}

func parserMethodName(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return parserMethodNamePlain(typ)
	case spec.NullableType:
		return parserMethodNamePlain(typ.Child) + "Nullable"
	case spec.ArrayType:
		return parserMethodNamePlain(typ.Child) + "Array"
	default:
		panic(fmt.Sprintf("Unsupported string param type: %v", typ.Plain))
	}
}

func parserMethodNamePlain(typ *spec.TypeDef) string {
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

func generateParamsParser(module module.Module) *generator.CodeFile {
	data := struct {
		PackageName string
	}{
		module.Name,
	}
	code := `
package [[.PackageName]]

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ParamsParser struct {
	values                   map[string][]string
	parseCommaSeparatedArray bool
	Errors                   []ParsingError
}

type ParsingError struct {
	Path    string
	Code    string
	Message *string
}

func New(values map[string][]string, parseCommaSeparatedArray bool) *ParamsParser {
	return &ParamsParser{values, parseCommaSeparatedArray, []ParsingError{}}
}

func (parser *ParamsParser) parseInt(name string, s string) int {
	v, err := strconv.Atoi(s)
	parser.addParsingError(name, "int", err)
	return v
}

func (parser *ParamsParser) parseInt64(name string, s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	parser.addParsingError(name, "int64", err)
	return v
}

func (parser *ParamsParser) parseFloat32(name string, s string) float32 {
	v, err := strconv.ParseFloat(s, 32)
	parser.addParsingError(name, "float32", err)
	return float32(v)
}

func (parser *ParamsParser) parseFloat64(name string, s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	parser.addParsingError(name, "float64", err)
	return v
}

func (parser *ParamsParser) parseDecimal(name string, s string) decimal.Decimal {
	v, err := decimal.NewFromString(s)
	parser.addParsingError(name, "decimal", err)
	return v
}

func (parser *ParamsParser) parseBool(name string, s string) bool {
	v, err := strconv.ParseBool(s)
	parser.addParsingError(name, "bool", err)
	return v
}

func (parser *ParamsParser) parseUuid(name string, s string) uuid.UUID {
	v, err := uuid.Parse(s)
	parser.addParsingError(name, "uuid", err)
	return v
}

func (parser *ParamsParser) parseDate(name string, s string) civil.Date {
	v, err := civil.ParseDate(s)
	parser.addParsingError(name, "date", err)
	return v
}

func (parser *ParamsParser) parseDateTime(name string, s string) civil.DateTime {
	t, err := time.Parse("2006-01-02T15:04:05.999Z", s)
	if err == nil {
		return civil.DateTimeOf(t)
	} else {
		v, err := civil.ParseDateTime(s)
		parser.addParsingError(name, "datetime", err)
		return v
	}
}

func (parser *ParamsParser) parseStringEnum(name string, s string, values []string) string {
	for _, value := range values {
		if s == value {
			return s
		}
	}
	parser.addValidationError(name, "unexpected_value", fmt.Sprintf("unexpected value %s, expected one of %s", s, strings.Join(values, ", ")))
	return ""
}

func (parser *ParamsParser) addParsingError(name string, format string, err error) {
	if err != nil {
		parser.addValidationError(name, "parsing_failed", fmt.Sprintf("failed to parse %s: %s ", format, err.Error()))
	}
}

func (parser *ParamsParser) addValidationError(name string, code, message string) {
	parser.Errors = append(parser.Errors, ParsingError{name, code, &message})
}

func (parser *ParamsParser) exactlyOneValue(name string) bool {
	pValues := parser.values[name]
	if len(pValues) != 1 {
		if len(pValues) == 0 {
			parser.addValidationError(name, "missing", "Parameter is missing")
		} else {
			parser.addValidationError(name, "too_many_values", fmt.Sprintf("expected exactly one value, found: %s", strings.Join(pValues, ",")))
		}
		return false
	} else {
		return true
	}
}

func (parser *ParamsParser) notMoreThenOneValue(name string) bool {
	pValues := parser.values[name]
	if len(pValues) > 1 {
		parser.addValidationError(name, "too_many_values", fmt.Sprintf("too many values provided, expected one or zero, found %s", strings.Join(pValues, ",")))
		return false
	} else {
		return true
	}
}

func (parser *ParamsParser) multipleValues(name string) []string {
	nameValues := parser.values[name]
	if parser.parseCommaSeparatedArray && len(nameValues) == 1 {
		return strings.Split(nameValues[0], ",")
	} else {
		return nameValues
	}
}

func (parser *ParamsParser) String(name string) string {
	if !parser.exactlyOneValue(name) {
		return ""
	}
	return parser.values[name][0]
}

func (parser *ParamsParser) StringNullable(name string) *string {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		return &pValues[0]
	}
}

func (parser *ParamsParser) StringDefaulted(name string, defaultValue string) string {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return *value
	}
}

func (parser *ParamsParser) StringArray(name string) []string {
	return parser.multipleValues(name)
}

func (parser *ParamsParser) Int(name string) int {
	if !parser.exactlyOneValue(name) {
		return 0
	}
	return parser.parseInt(name, parser.values[name][0])
}

func (parser *ParamsParser) IntNullable(name string) *int {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseInt(name, pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) IntDefaulted(name string, defaultValue int) int {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseInt(name, *value)
	}
}

func (parser *ParamsParser) IntArray(name string) []int {
	stringValues := parser.StringArray(name)
	convertedValues := []int{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseInt(name, stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) Int64(name string) int64 {
	if !parser.exactlyOneValue(name) {
		return 0
	}
	return parser.parseInt64(name, parser.values[name][0])
}

func (parser *ParamsParser) Int64Nullable(name string) *int64 {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseInt64(name, pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) Int64Defaulted(name string, defaultValue int64) int64 {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseInt64(name, *value)
	}
}

func (parser *ParamsParser) Int64Array(name string) []int64 {
	stringValues := parser.StringArray(name)
	convertedValues := []int64{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseInt64(name, stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) Float32(name string) float32 {
	if !parser.exactlyOneValue(name) {
		return 0
	}
	return parser.parseFloat32(name, parser.values[name][0])
}

func (parser *ParamsParser) Float32Nullable(name string) *float32 {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseFloat32(name, pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) Float32Defaulted(name string, defaultValue float32) float32 {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseFloat32(name, *value)
	}
}

func (parser *ParamsParser) Float32Array(name string) []float32 {
	stringValues := parser.StringArray(name)
	convertedValues := []float32{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseFloat32(name, stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) Float64(name string) float64 {
	if !parser.exactlyOneValue(name) {
		return 0
	}
	return parser.parseFloat64(name, parser.values[name][0])
}

func (parser *ParamsParser) Float64Nullable(name string) *float64 {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseFloat64(name, pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) Float64Defaulted(name string, defaultValue float64) float64 {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseFloat64(name, *value)
	}
}

func (parser *ParamsParser) Float64Array(name string) []float64 {
	stringValues := parser.StringArray(name)
	convertedValues := []float64{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseFloat64(name, stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) Decimal(name string) decimal.Decimal {
	if !parser.exactlyOneValue(name) {
		return decimal.Decimal{}
	}
	return parser.parseDecimal(name, parser.values[name][0])
}

func (parser *ParamsParser) DecimalNullable(name string) *decimal.Decimal {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseDecimal(name, pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) DecimalDefaulted(name string, defaultValue decimal.Decimal) decimal.Decimal {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseDecimal(name, *value)
	}
}

func (parser *ParamsParser) DecimalArray(name string) []decimal.Decimal {
	stringValues := parser.StringArray(name)
	convertedValues := []decimal.Decimal{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseDecimal(name, stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) Bool(name string) bool {
	if !parser.exactlyOneValue(name) {
		return false
	}
	return parser.parseBool(name, parser.values[name][0])
}

func (parser *ParamsParser) BoolNullable(name string) *bool {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseBool(name, pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) BoolDefaulted(name string, defaultValue bool) bool {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseBool(name, *value)
	}
}

func (parser *ParamsParser) BoolArray(name string) []bool {
	stringValues := parser.StringArray(name)
	convertedValues := []bool{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseBool(name, stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) Uuid(name string) uuid.UUID {
	if !parser.exactlyOneValue(name) {
		return uuid.Nil
	}
	return parser.parseUuid(name, parser.values[name][0])
}

func (parser *ParamsParser) UuidNullable(name string) *uuid.UUID {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseUuid(name, pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) UuidDefaulted(name string, defaultValue uuid.UUID) uuid.UUID {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseUuid(name, *value)
	}
}

func (parser *ParamsParser) UuidArray(name string) []uuid.UUID {
	stringValues := parser.StringArray(name)
	convertedValues := []uuid.UUID{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseUuid(name, stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) Date(name string) civil.Date {
	if !parser.exactlyOneValue(name) {
		return civil.Date{}
	}
	return parser.parseDate(name, parser.values[name][0])
}

func (parser *ParamsParser) DateNullable(name string) *civil.Date {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseDate(name, pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) DateDefaulted(name string, defaultValue civil.Date) civil.Date {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseDate(name, *value)
	}
}

func (parser *ParamsParser) DateArray(name string) []civil.Date {
	stringValues := parser.StringArray(name)
	convertedValues := []civil.Date{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseDate(name, stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) DateTime(name string) civil.DateTime {
	if !parser.exactlyOneValue(name) {
		return civil.DateTime{}
	}
	return parser.parseDateTime(name, parser.values[name][0])
}

func (parser *ParamsParser) DateTimeNullable(name string) *civil.DateTime {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseDateTime(name, pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) DateTimeDefaulted(name string, defaultValue civil.DateTime) civil.DateTime {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseDateTime(name, *value)
	}
}

func (parser *ParamsParser) DateTimeArray(name string) []civil.DateTime {
	stringValues := parser.StringArray(name)
	convertedValues := []civil.DateTime{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseDateTime(name, stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) StringEnum(name string, values []string) string {
	if !parser.exactlyOneValue(name) {
		return ""
	}
	return parser.parseStringEnum(name, parser.values[name][0], values)
}

func (parser *ParamsParser) StringEnumNullable(name string, values []string) *string {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseStringEnum(name, pValues[0], values)
		return &convertedValue
	}
}

func (parser *ParamsParser) StringEnumDefaulted(name string, values []string, defaultValue string) string {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseStringEnum(name, *value, values)
	}
}

func (parser *ParamsParser) StringEnumArray(name string, values []string) []string {
	stringValues := parser.StringArray(name)
	convertedValues := []string{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseStringEnum(name, stringValue, values))
	}
	return convertedValues
}
`

	code, _ = generator.ExecuteTemplate(code, data)
	return &generator.CodeFile{module.GetPath("parser.go"), strings.TrimSpace(code)}
}
