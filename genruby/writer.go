package genruby

import "github.com/specgen-io/specgen/v2/sources"

func NewRubyWriter() *sources.Writer {
	return sources.NewWriter("  ", 2)
}
