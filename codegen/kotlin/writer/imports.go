package writer

import (
	"fmt"
	"kotlin/packages"
)

type imports struct {
	lines []string
}

func NewImports() *imports {
	return &imports{lines: []string{}}
}

func (imports *imports) Add(items ...string) {
	for _, item := range items {
		imports.lines = append(imports.lines, fmt.Sprintf("import %s", item))
	}
}

func (imports *imports) Package(thePackage packages.Package) {
	imports.Add(thePackage.PackageName)
}

func (imports *imports) PackageStar(thePackage packages.Package) {
	imports.Add(thePackage.PackageStar)
}

func (imports *imports) Lines() []string {
	return imports.lines
}
