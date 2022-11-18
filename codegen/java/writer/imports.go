package writer

import (
	"fmt"
	"java/packages"
)

type imports struct {
	lines []string
}

func NewImports() *imports {
	return &imports{lines: []string{}}
}

func (imports *imports) add(line string) {
	imports.lines = append(imports.lines, line)
}

func (imports *imports) Add(classes ...string) {
	for _, class := range classes {
		imports.add(fmt.Sprintf("import %s;", class))
	}
}

func (imports *imports) Star(thePackage packages.Package) {
	imports.Add(thePackage.PackageStar)
}

func (imports *imports) AddStatic(members ...string) {
	for _, member := range members {
		imports.add(fmt.Sprintf("import static %s;", member))
	}
}

func (imports *imports) StaticStar(thePackage packages.Package) {
	imports.AddStatic(thePackage.PackageStar)
}

func (imports *imports) Lines() []string {
	return imports.lines
}
