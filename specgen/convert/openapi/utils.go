package openapi

import "github.com/specgen-io/specgen/spec/v2"

func name(source string) spec.Name {
	return spec.Name{source, nil}
}

var emptyType = spec.Type{*spec.Plain(spec.TypeEmpty), nil}
