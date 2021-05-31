package gengo

import (
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

func generateParamsParser(packageName string, path string) *gen.TextFile {
	code := `
package [[.PackageName]]

import (
	"fmt"
	"strconv"
	"strings"
)

type ParamsParser struct {
	values map[string][]string
	Errors []error
}

func NewParamsParser(values map[string][]string) *ParamsParser {
	return &ParamsParser{values, []error{}}
}

func (parser *ParamsParser) parseInt(s string) int {
	v, err := strconv.Atoi(s)
	parser.addError(err)
	return v
}

func (parser *ParamsParser) parseInt64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	parser.addError(err)
	return v
}

func (parser *ParamsParser) parseFloat32(s string) float32 {
	v, err := strconv.ParseFloat(s, 32)
	parser.addError(err)
	return float32(v)
}

func (parser *ParamsParser) parseFloat64(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	parser.addError(err)
	return v
}

func (parser *ParamsParser) parseStringEnum(s string, vs []string) string {
	for _, v := range vs {
		if s == v {
			return s
		}
	}
	parser.addError(fmt.Errorf("unexpected value %s, expected one of %s", s, strings.Join(vs, ", ")))
	return ""
}

func (parser *ParamsParser) addError(err error) {
	if err != nil {
		parser.Errors = append(parser.Errors, err)
	}
}

func (parser *ParamsParser) exactlyOneValue(name string) bool {
	pValues := parser.values[name]
	if len(pValues) != 1 {
		parser.addError(fmt.Errorf("param %s: expected exactly one value, found: %s", name, pValues))
		return false
	} else {
		return true
	}
}

func (parser *ParamsParser) notMoreThenOneValue(name string) bool {
	pValues := parser.values[name]
	if len(pValues) > 1 {
		parser.addError(fmt.Errorf("param %s: too many values, expected one or zero, found %d", name, len(pValues)))
		return false
	} else {
		return true
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
	return parser.values[name]
}

func (parser *ParamsParser) Int(name string) int {
	if !parser.exactlyOneValue(name) {
		return 0
	}
	return parser.parseInt(parser.values[name][0])
}

func (parser *ParamsParser) IntNullable(name string) *int {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseInt(pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) IntDefaulted(name string, defaultValue int) int {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseInt(*value)
	}
}

func (parser *ParamsParser) IntArray(name string) []int {
	stringValues := parser.StringArray(name)
	convertedValues := []int{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseInt(stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) Int64(name string) int64 {
	if !parser.exactlyOneValue(name) {
		return 0
	}
	return parser.parseInt64(parser.values[name][0])
}

func (parser *ParamsParser) Int64Nullable(name string) *int64 {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseInt64(pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) Int64Defaulted(name string, defaultValue int64) int64 {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseInt64(*value)
	}
}

func (parser *ParamsParser) Int64Array(name string) []int64 {
	stringValues := parser.StringArray(name)
	convertedValues := []int64{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseInt64(stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) Float32(name string) float32 {
	if !parser.exactlyOneValue(name) {
		return 0
	}
	return parser.parseFloat32(parser.values[name][0])
}

func (parser *ParamsParser) Float32Nullable(name string) *float32 {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseFloat32(pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) Float32Defaulted(name string, defaultValue float32) float32 {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseFloat32(*value)
	}
}

func (parser *ParamsParser) Float32Array(name string) []float32 {
	stringValues := parser.StringArray(name)
	convertedValues := []float32{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseFloat32(stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) Float64(name string) float64 {
	if !parser.exactlyOneValue(name) {
		return 0
	}
	return parser.parseFloat64(parser.values[name][0])
}

func (parser *ParamsParser) Float64Nullable(name string) *float64 {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseFloat64(pValues[0])
		return &convertedValue
	}
}

func (parser *ParamsParser) Float64Defaulted(name string, defaultValue float64) float64 {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseFloat64(*value)
	}
}

func (parser *ParamsParser) Float64Array(name string) []float64 {
	stringValues := parser.StringArray(name)
	convertedValues := []float64{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseFloat64(stringValue))
	}
	return convertedValues
}

func (parser *ParamsParser) StringEnum(name string, values []string) string {
	if !parser.exactlyOneValue(name) {
		return ""
	}
	return parser.parseStringEnum(parser.values[name][0], values)
}

func (parser *ParamsParser) StringEnumNullable(name string, values []string) *string {
	if !parser.notMoreThenOneValue(name) {
		return nil
	}
	pValues := parser.values[name]
	if len(pValues) == 0 {
		return nil
	} else {
		convertedValue := parser.parseStringEnum(pValues[0], values)
		return &convertedValue
	}
}

func (parser *ParamsParser) StringEnumDefaulted(name string, values []string, defaultValue string) string {
	value := parser.StringNullable(name)
	if value == nil {
		return defaultValue
	} else {
		return parser.parseStringEnum(*value, values)
	}
}

func (parser *ParamsParser) StringEnumArray(name string, values []string) []string {
	stringValues := parser.StringArray(name)
	convertedValues := []string{}
	for _, stringValue := range stringValues {
		convertedValues = append(convertedValues, parser.parseStringEnum(stringValue, values))
	}
	return convertedValues
}
`

	code, _ = gen.ExecuteTemplate(code, struct { PackageName string } {packageName })
	return &gen.TextFile{path, strings.TrimSpace(code)}
}
