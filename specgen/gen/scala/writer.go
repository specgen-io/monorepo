package scala

import (
	"github.com/specgen-io/specgen/generator/v2"
)

var ScalaConfig = generator.Config{"  ", 2, nil}

func NewScalaWriter() *generator.Writer {
	return generator.NewWriter(ScalaConfig)
}
