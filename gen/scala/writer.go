package scala

import "github.com/specgen-io/specgen/v2/sources"

func NewScalaWriter() *sources.Writer {
	return sources.NewWriter("  ", 2)
}
