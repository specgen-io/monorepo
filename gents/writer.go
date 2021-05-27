package gents

import "gopkg.in/specgen-io/specgen.v2/gen"

func NewTsWriter() *gen.Writer {
	return gen.NewWriter("    ", 2)
}
