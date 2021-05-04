package gents

import (
	"specgen/gen"
)

func NewTsWriter() *gen.Writer {
	return gen.NewWriter("    ", 2)
}
