package service

import (
	"fmt"

	"generator"
	"spec"
	"typescript/module"
	"typescript/validations"
)

type ServiceGenerator interface {
	VersionRouting(version *spec.Version, targetModule module.Module, modelsModule, validationModule, paramsModule, errorsModule, responsesModule module.Module) *generator.CodeFile
	SpecRouter(specification *spec.Spec, rootModule module.Module, specRouterModule module.Module) *generator.CodeFile
	Responses(targetModule, validationModule, errorsModule module.Module) *generator.CodeFile
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
