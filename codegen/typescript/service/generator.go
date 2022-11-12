package service

import (
	"fmt"

	"generator"
	"spec"
	"typescript/validations"
)

type ServiceGenerator interface {
	VersionRouting(version *spec.Version) *generator.CodeFile
	SpecRouter(specification *spec.Spec) *generator.CodeFile
	Responses() *generator.CodeFile
}

type Generator struct {
	validations.Validation
	ServiceGenerator
	Modules *Modules
}

func NewServiceGenerator(server, validationName string, modules *Modules) *Generator {
	validation := validations.New(validationName, &(modules.Modules))
	var serviceGenerator ServiceGenerator = nil
	switch server {
	case Express:
		serviceGenerator = &expressGenerator{modules, validation}
	case Koa:
		serviceGenerator = &koaGenerator{modules, validation}
	default:
		panic(fmt.Sprintf("Unknown server: %s", server))
	}
	return &Generator{validation, serviceGenerator, modules}
}

var Express = "express"
var Koa = "koa"

var Servers = []string{Express, Koa}
