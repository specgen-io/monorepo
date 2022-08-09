package openapi

import "spec"

func name(source string) spec.Name {
	return spec.Name{source, nil}
}

var emptyType = spec.Type{*spec.Plain(spec.TypeEmpty), nil}
