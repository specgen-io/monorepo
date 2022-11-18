package writer

import (
	"fmt"
	"generator"
	"golang/module"
	"strings"
)

func GoConfig() generator.Config {
	return generator.Config{"\t", 2, map[string]string{"PERCENT_": "%"}}
}

type Writer struct {
	generator.Writer
	filename string
	module   module.Module
	Imports  *imports
}

func New(module module.Module, filename string) *Writer {
	config := GoConfig()
	w := generator.NewWriter(config)
	return &Writer{w, module.GetPath(filename), module, NewImports()}
}

func (w *Writer) Indented() *Writer {
	return &Writer{w.Writer.Indented(), w.filename, w.module, w.Imports}
}

func (w *Writer) IndentedWith(size int) *Writer {
	return &Writer{w.Writer.IndentedWith(size), w.filename, w.module, w.Imports}
}

func (w *Writer) ToCodeFile() *generator.CodeFile {
	lines := []string{fmt.Sprintf("package %s", w.module.Name), ``}
	lines = append(lines, w.Imports.Lines()...)
	code := strings.Join(lines, "\n") + w.String()
	return &generator.CodeFile{w.filename, code}
}
