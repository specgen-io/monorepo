package genjava

import (
	"github.com/pinzolo/casee"
	"strings"
)

var ToPascalCase = casee.ToPascalCase

func JoinDelimParams(params []string) string {
	return strings.Join(params, ", ")
}

func JoinParams(params []string) string {
	return strings.Join(params, "")
}
