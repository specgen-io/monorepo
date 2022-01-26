package service

import (
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/v2/genjava/imports"
	"strings"
)

var ToPascalCase = casee.ToPascalCase

var GenerateImports = imports.GenerateImports
var GeneralImports = imports.GeneralImports

func joinParams(params []string) string {
	return strings.Join(params, ", ")
}
