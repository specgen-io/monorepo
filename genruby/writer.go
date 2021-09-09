package genruby

import "github.com/specgen-io/specgen/v2/gen"

func NewRubyWriter() *gen.Writer {
	return gen.NewWriter("  ", 2)
}