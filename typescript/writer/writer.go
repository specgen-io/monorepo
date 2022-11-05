package writer

import (
	"generator"
	"typescript/module"
)

func TsConfig() generator.Config {
	return generator.Config{"    ", 2, nil}
}

type TsWriter struct {
	generator.Writer
	module module.Module
}

func New(module module.Module) *TsWriter {
	return &TsWriter{
		generator.NewWriter2(module.GetPath(), TsConfig()),
		module,
	}
}

func (w *TsWriter) Imports() *imports {
	return NewImports(w.module)
}
