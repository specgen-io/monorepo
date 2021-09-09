package gengo

import "specgen/gen"

func NewGoWriter() *gen.Writer {
	return gen.NewWriter("\t", 2)
}
