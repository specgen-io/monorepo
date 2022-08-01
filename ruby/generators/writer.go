package generators

import (
	"github.com/specgen-io/specgen/generator/v2"
)

var RubyConfig = generator.Config{"  ", 2, nil}

func NewRubyWriter() *generator.Writer {
	return generator.NewWriter(RubyConfig)
}
