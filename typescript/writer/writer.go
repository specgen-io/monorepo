package writer

import (
	"generator"
	"typescript/modules"
)

func TsConfig() generator.Config {
	return generator.Config{"    ", 2, nil}
}

func New(module modules.Module) generator.Writer {
	return generator.NewWriter2(module.GetPath(), TsConfig())
}
