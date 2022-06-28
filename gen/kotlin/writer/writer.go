package writer

import (
	"github.com/specgen-io/specgen/v2/generator"
)

func NewKotlinWriter() *generator.Writer {
	return generator.NewWriter("\t", 2)
}
