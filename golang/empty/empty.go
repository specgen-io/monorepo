package empty

import (
	"generator"
	"golang/module"
	"golang/writer"
)

func GenerateEmpty(emptyModule module.Module) *generator.CodeFile {
	w := writer.New(emptyModule, `empty.go`)
	w.Lines(`
type Type struct{}

var Value = Type{}
`)
	return w.ToCodeFile()
}
