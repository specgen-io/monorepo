package genjava

import "specgen/gen"

func NewJavaWriter() *gen.Writer {
	return gen.NewWriter("\t", 2)
}
