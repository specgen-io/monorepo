package gents

import "github.com/specgen-io/specgen/v2/sources"

func NewTsWriter() *sources.Writer {
	return sources.NewWriter("    ", 2)
}
