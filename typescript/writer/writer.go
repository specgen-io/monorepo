package writer

import (
	"github.com/specgen-io/specgen/generator/v2"
)

var TsConfig = generator.Config{"    ", 2, nil}

func NewTsWriter() *generator.Writer {
	return generator.NewWriter(TsConfig)
}
