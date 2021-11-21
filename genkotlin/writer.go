package genkotlin

import "github.com/specgen-io/specgen/v2/gen"

func NewKotlinWriter() *gen.Writer {
	return gen.NewWriter("\t", 2)
}
