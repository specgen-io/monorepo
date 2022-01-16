package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/gents/validation"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

type serviceGenerator interface {
	generateVersionRouting(version *spec.Version, validationModule, paramsModule, module modules.Module) *sources.CodeFile
	generateSpecRouter(specification *spec.Spec, rootModule modules.Module, module modules.Module) *sources.CodeFile
}

func newServiceGenerator(server string, validation validation.Validation) serviceGenerator {
	if server == "express" {
		return &expressGenerator{validation}
	}
	if server == "koa" {
		return &koaGenerator{validation}
	}
	panic(fmt.Sprintf("Unknown server: %s", server))
}
