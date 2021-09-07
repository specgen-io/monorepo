package gengo

import (
	"github.com/pinzolo/casee"
	"strings"
)

var ToPascalCase = casee.ToPascalCase
var ToUpperCase = casee.ToUpperCase

func JoinDelimParams(params []string) string {
	return strings.Join(params, ", ")
}

func JoinParams(params []string) string {
	return strings.Join(params, "")
}
