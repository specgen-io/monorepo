package writer

import (
	"generator"
	"typescript/module"
)

func TsConfig() generator.Config {
	return generator.Config{"    ", 2, nil}
}

func New(module module.Module) generator.Writer {
	return generator.NewWriter2(module.GetPath(), TsConfig())
}
