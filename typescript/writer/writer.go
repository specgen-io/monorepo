package writer

import (
	"generator"
)

var TsConfig = generator.Config{"    ", 2, nil}

func NewTsWriter() generator.Writer {
	return generator.NewWriter(TsConfig)
}
