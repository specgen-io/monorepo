package service

import (
	"fmt"

	"github.com/specgen-io/specgen/generator/v2"
	"github.com/specgen-io/specgen/spec/v2"
	"github.com/specgen-io/specgen/typescript/v2/modules"
	"github.com/specgen-io/specgen/typescript/v2/validations"
)

type ServiceGenerator interface {
	VersionRouting(version *spec.Version, validationModule, paramsModule, module modules.Module) *generator.CodeFile
	SpecRouter(specification *spec.Spec, rootModule modules.Module, module modules.Module) *generator.CodeFile
}

func NewServiceGenerator(server string, validation validations.Validation) ServiceGenerator {
	if server == Express {
		return &expressGenerator{validation}
	}
	if server == Koa {
		return &koaGenerator{validation}
	}
	panic(fmt.Sprintf("Unknown server: %s", server))
}

var Express = "express"
var Koa = "koa"

var Servers = []string{Express, Koa}
