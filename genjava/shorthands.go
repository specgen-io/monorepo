package genjava

import (
	"github.com/pinzolo/casee"
	"strings"
)

var ToPascalCase = casee.ToPascalCase

func JoinParams(params []string) string {
	return strings.Join(params, ", ")
}
