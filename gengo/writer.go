package gengo

import "github.com/specgen-io/specgen/v2/sources"

func NewGoWriter() *sources.Writer {
	return sources.NewWriter("\t", 2)
}
