package gen

import (
	"github.com/specgen-io/specgen/v2/sources"
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
	Title     string
	Usage     string
	Args      []GeneratorArg
	Generator GeneratorFunc
}

type GeneratorFunc func(specification *spec.Spec, params map[Arg]string) *sources.Sources
