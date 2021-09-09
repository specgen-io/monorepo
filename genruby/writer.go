package genruby

import "specgen/gen"

func NewRubyWriter() *gen.Writer {
	return gen.NewWriter("  ", 2)
}