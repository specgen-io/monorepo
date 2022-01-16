package service

import (
	"github.com/pinzolo/casee"
	"strings"
)

var ToPascalCase = casee.ToPascalCase

func joinParams(params []string) string {
	return strings.Join(params, ", ")
}
