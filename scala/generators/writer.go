package generators

import (
	"generator"
)

var ScalaConfig = generator.Config{"  ", 2, nil}

func NewScalaWriter() generator.Writer {
	return generator.NewWriter(ScalaConfig)
}
