package writer

import (
	"fmt"
	"generator"
	"java/packages"
	"strings"
)

func JavaConfig() generator.Config {
	return generator.Config{"\t", 2, map[string]string{}}
}

type Writer struct {
	generator.Writer
	thePackage packages.Package
	className  string
	Imports    *imports
}

func New(thePackage packages.Package, className string) *Writer {
	config := JavaConfig()
	config.Substitutions["[[.ClassName]]"] = className
	return &Writer{generator.NewWriter(config), thePackage, className, NewImports()}
}

func (w *Writer) Indented() *Writer {
	return &Writer{w.Writer.Indented(), w.thePackage, w.className, w.Imports}
}

func (w *Writer) IndentedWith(size int) *Writer {
	return &Writer{w.Writer.IndentedWith(size), w.thePackage, w.className, w.Imports}
}

func (w *Writer) ToCodeFile() *generator.CodeFile {
	lines := []string{
		fmt.Sprintf(`package %s;`, w.thePackage.PackageName),
		"",
	}
	lines = append(lines, w.Imports.Lines()...)
	code := strings.Join(lines, "\n") + w.String()
	return &generator.CodeFile{w.thePackage.GetPath(fmt.Sprintf("%s.java", w.className)), code}
}
