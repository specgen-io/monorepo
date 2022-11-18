package writer

import (
	"fmt"
	"generator"
	"scala/packages"
)

func ScalaConfig() generator.Config {
	return generator.Config{"  ", 2, map[string]string{}}
}

type Writer struct {
	generator.Writer
	filename string
}

func New(thepackage packages.Package, className string) *Writer {
	config := ScalaConfig()
	filename := thepackage.GetPath(fmt.Sprintf("%s.scala", className))
	config.Substitutions["[[.ClassName]]"] = className
	w := generator.NewWriter(filename, config)
	w.Line(`package %s`, thepackage.PackageName)
	w.EmptyLine()
	return &Writer{w, filename}
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
