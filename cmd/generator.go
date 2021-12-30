package cmd

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
)

type GeneratorArg struct {
	Arg
	Values   []string
	Required bool
	Default  string
}

type Generator struct {
	Name      string
	Usage     string
	Args      []GeneratorArg
	Generator GeneratorFunc
}

type GeneratorFunc func(specification *spec.Spec, params map[string]string) *gen.Sources
