package genscala

import "github.com/specgen-io/specgen/v2/gen"

func NewScalaWriter() *gen.Writer {
	return gen.NewWriter("  ", 2)
}