package ruby

import (
	"github.com/specgen-io/specgen/v2/generator"
)

func NewRubyWriter() *generator.Writer {
	return generator.NewWriter("  ", 2)
}
