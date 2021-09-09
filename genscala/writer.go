package genscala

import "specgen/gen"

func NewScalaWriter() *gen.Writer {
	return gen.NewWriter("  ", 2)
}