package genjava

import (
	"github.com/pinzolo/casee"
	"strings"
)

var ToPascalCase = casee.ToPascalCase
var ToUpperCase = casee.ToUpperCase

func JoinParams(params []string) string {
	return strings.Join(params, ", ")
}
