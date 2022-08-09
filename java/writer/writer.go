package writer

import (
	"generator"
)

var JavaConfig = generator.Config{"\t", 2, nil}

func NewJavaWriter() *generator.Writer {
	return generator.NewWriter(JavaConfig)
}
