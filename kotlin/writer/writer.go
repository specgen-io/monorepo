package writer

import (
	"generator"
)

var KotlinConfig = generator.Config{"\t", 2, nil}

func NewKotlinWriter() *generator.Writer {
	return generator.NewWriter(KotlinConfig)
}
