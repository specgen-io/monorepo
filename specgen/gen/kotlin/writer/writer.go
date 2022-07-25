package writer

import (
	"github.com/specgen-io/specgen/generator/v2"
)

var KotlinConfig = generator.Config{"\t", 2, nil}

func NewKotlinWriter() *generator.Writer {
	return generator.NewWriter(KotlinConfig)
}
