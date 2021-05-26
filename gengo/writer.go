package gengo

import "github.com/specgen-io/specgen/gen"

func NewGoWriter() *gen.Writer {
	return gen.NewWriter("\t", 2)
}
