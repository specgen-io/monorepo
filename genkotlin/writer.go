package genkotlin

import "github.com/specgen-io/specgen/v2/sources"

func NewKotlinWriter() *sources.Writer {
	return sources.NewWriter("\t", 2)
}
