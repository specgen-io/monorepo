package scala

import (
	"github.com/specgen-io/specgen/v2/generator"
)

func NewScalaWriter() *generator.Writer {
	return generator.NewWriter("  ", 2)
}
