package writer

import (
	"generator"
)

var RubyConfig = generator.Config{"  ", 2, nil}

type Writer struct {
	generator.Writer
	filename string
}

func New(filename string) *Writer {
	return &Writer{generator.NewWriter(RubyConfig), filename}
}

func (w *Writer) Indented() *Writer {
	return &Writer{w.Writer.Indented(), w.filename}
}

func (w *Writer) IndentedWith(size int) *Writer {
	return &Writer{w.Writer.IndentedWith(size), w.filename}
}

func (w *Writer) ToCodeFile() *generator.CodeFile {
	return &generator.CodeFile{w.filename, w.String()}
}
