package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/gents/validations"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

type ServiceGenerator interface {
	VersionRouting(version *spec.Version, validationModule, paramsModule, module modules.Module) *sources.CodeFile
	SpecRouter(specification *spec.Spec, rootModule modules.Module, module modules.Module) *sources.CodeFile
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
