package gents

import "github.com/specgen-io/specgen/gen"

func NewTsWriter() *gen.Writer {
	return gen.NewWriter("    ", 2)
}
