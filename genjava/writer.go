package genjava

import "github.com/specgen-io/specgen/v2/gen"

func NewJavaWriter() *gen.Writer {
	return gen.NewWriter("\t", 2)
}
