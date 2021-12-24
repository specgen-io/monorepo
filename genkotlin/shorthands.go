package genkotlin

import (
	"github.com/pinzolo/casee"
	"strings"
)

var ToPascalCase = casee.ToPascalCase

func JoinDelimParams(params []string) string {
	return strings.Join(params, ", ")
}

func TrimSlash(param string) string {
	return strings.Trim(param, "/")
}
