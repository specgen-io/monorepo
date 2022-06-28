package generator

import (
	"github.com/specgen-io/specgen/v2/spec"
)

type Arg struct {
	Name        string
	Title       string
	Description string
}

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

type GeneratorArgsValues map[Arg]string

type GeneratorFunc func(specification *spec.Spec, params GeneratorArgsValues) *Sources
