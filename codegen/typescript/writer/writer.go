package writer

import (
	"generator"
	"strings"
	"typescript/module"
)

func TsConfig() generator.Config {
	return generator.Config{"    ", 2, nil}
}

type Writer struct {
	generator.Writer
	filename string
	module   module.Module
	Imports  *imports
}

func New(module module.Module) *Writer {
	return &Writer{
		generator.NewWriter(TsConfig()),
		module.GetPath(),
		module,
		NewImports(module),
	}
}

func (w *Writer) Indented() *Writer {
	return &Writer{w.Writer.Indented(), w.filename, w.module, w.Imports}
}

func (w *Writer) IndentedWith(size int) *Writer {
	return &Writer{w.Writer.IndentedWith(size), w.filename, w.module, w.Imports}
}

func (w *Writer) ToCodeFile() *generator.CodeFile {
	lines := []string{}
	imports := w.Imports.Lines()
	if len(imports) > 0 {
		lines = append(lines, imports...)
		lines = append(lines, "")
	}
	lines = append(lines, w.Code()...)
	code := strings.Join(w.Imports.Lines(), "\n")
	return &generator.CodeFile{w.filename, code}
}
