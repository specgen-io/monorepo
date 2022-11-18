package writer

import (
	"fmt"
	"generator"
	"java/packages"
)

func JavaConfig() generator.Config {
	return generator.Config{"\t", 2, map[string]string{}}
}

type Writer struct {
	generator.Writer
	filename string
}

func New(thePackage packages.Package, className string) *Writer {
	config := JavaConfig()
	filename := thePackage.GetPath(fmt.Sprintf("%s.java", className))
	config.Substitutions["[[.ClassName]]"] = className
	w := generator.NewWriter(config)
	w.Line(`package %s;`, thePackage.PackageName)
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
