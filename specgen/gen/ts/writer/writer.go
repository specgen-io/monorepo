package writer

import (
	"github.com/specgen-io/specgen/specgen/v2/generator"
)

var TsConfig = generator.Config{"    ", 2, nil}

func NewTsWriter() *generator.Writer {
	return generator.NewWriter(TsConfig)
}
