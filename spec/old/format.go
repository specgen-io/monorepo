package old

import (
	"fmt"
	"regexp"
)

type Format struct {
	Name    string
	Regex   string
	Example string
}

type FormatError struct {
	Value  string
	Format Format
}

func (err *FormatError) Error() string {
	message := fmt.Sprintf("format error: '%s' is in wrong format, should be %s; examples: %s", err.Value, err.Format.Name, err.Format.Example)
	return message
}

func (format *Format) Check(value string) error {
	isMatching, err := regexp.MatchString(format.Regex, value)
	if err != nil {
		return err
	}
	if !isMatching {
		return &FormatError{Value: value, Format: *format}
	}
	return nil
}

var PascalCase = Format{Name: "pascal case", Regex: "^[A-Z][a-z0-9]+([A-Z][a-z0-9]+)*$", Example: "ThisIsPascalCase"}

var SnakeCase = Format{Name: "snake case", Regex: "^[a-z][a-z0-9]*(_[a-z][a-z0-9]*)*$", Example: "this_is_snake_case"}

var LowerCase = Format{Name: "lower case", Regex: "^[a-z][a-z]*[0-9]*$", Example: "thisislowercase"}

var HttpMethod = Format{Name: "HTTP method", Regex: "^GET|POST|PUT|DELETE$", Example: "GET, POST, PUT, DELETE"}

var Integer = Format{Name: "integer", Regex: "^[-+]?\\d+$", Example: "123"}

var Float = Format{Name: "float", Regex: "^[-+]?\\d+\\.?\\d*$", Example: "123.4"}

var Boolean = Format{Name: "boolean", Regex: "^true$|^false$", Example: "true or false"}

var UUID = Format{Name: "uuid", Regex: "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$", Example: "fbd3036f-0f1c-4e98-b71c-d4cd61213f90"}

var Date = Format{Name: "date", Regex: "^\\d{4}-\\d{2}-\\d{2}$", Example: "2019-12-31"}

var DateTime = Format{Name: "datetime", Regex: "^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}$", Example: "2019-12-31T15:53:45"}

var HttpParams = Format{Name: "symbols a-z, A-Z, 0-9, -, _ starting from letter", Regex: "^[a-zA-Z]([a-zA-Z0-9_-])*$", Example: "some123, Some-123, some_123"}

var JsonField = Format{Name: "symbols a-z, A-Z, 0-9, _ starting from letter or _", Regex: "^[a-zA-Z_]([a-zA-Z_0-9])*$", Example: "some123, Some123, some_123, _some123"}

var VersionFormat = Format{Name: "version", Regex: "^v\\d+$", Example: "v1"}

func FormatOr(f1 Format, f2 Format) Format {
	return Format{
		Name:    fmt.Sprintf("%s or %s", f1.Name, f2.Name),
		Regex:   fmt.Sprintf("%s|%s", f1.Regex, f2.Regex),
		Example: fmt.Sprintf("%s or %s", f1.Example, f2.Example),
	}
}
