package gengo

import "gopkg.in/specgen-io/specgen.v2/gen"

func NewGoWriter() *gen.Writer {
	return gen.NewWriter("\t", 2)
}
