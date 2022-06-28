package writer

import (
	"github.com/specgen-io/specgen/v2/generator"
)

func NewTsWriter() *generator.Writer {
	return generator.NewWriter("    ", 2)
}
