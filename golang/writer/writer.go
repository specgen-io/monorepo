package writer

import (
	"generator"
	"golang/module"
)

func GoConfig() generator.Config {
	return generator.Config{"\t", 2, map[string]string{"PERCENT_": "%"}}
}

type Writer struct {
	generator.Writer
	filename string
	module   module.Module
}

func New(module module.Module, filename string) *Writer {
	config := GoConfig()
	w := generator.NewWriter(module.GetPath(filename), config)
	w.Line("package %s", module.Name)
	w.EmptyLine()
	return &Writer{w, module.GetPath(filename), module}
}

func (w *Writer) Indented() *Writer {
	return &Writer{w.Writer.Indented(), w.filename, w.module}
}

func (w *Writer) IndentedWith(size int) *Writer {
	return &Writer{w.Writer.IndentedWith(size), w.filename, w.module}
}

func (w *Writer) ToCodeFile() *generator.CodeFile {
	return &generator.CodeFile{w.filename, w.String()}
}
