package generators

import "github.com/specgen-io/specgen/v2/generator"

var All = []generator.Generator{}

func Add(generators ...generator.Generator) {
	All = append(All, generators...)
}
