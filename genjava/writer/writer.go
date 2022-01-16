package writer

import "github.com/specgen-io/specgen/v2/sources"

func NewJavaWriter() *sources.Writer {
	return sources.NewWriter("\t", 2)
}
