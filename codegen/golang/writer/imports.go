package writer

import (
	"fmt"
	"golang/module"
	"sort"
)

type imports struct {
	imports map[string]string
}

func NewImports() *imports {
	return &imports{imports: make(map[string]string)}
}

func (self *imports) Module(module module.Module) *imports {
	self.Add(module.Package)
	return self
}

func (self *imports) ModuleAliased(module module.Module) *imports {
	if module.Alias == "" {
		panic(fmt.Sprintf(`module %s does not have alias and can't imported as aliased'`, module.Package))
	}
	self.AddAliased(module.Package, module.Alias)
	return self
}

func (self *imports) Add(theImport string) *imports {
	self.imports[theImport] = ""
	return self
}

func (self *imports) AddAliased(theImport string, alias string) *imports {
	self.imports[theImport] = alias
	return self
}

func (self *imports) Lines() []string {
	if len(self.imports) > 0 {
		imports := make([]string, 0, len(self.imports))
		for theImport := range self.imports {
			imports = append(imports, theImport)
		}
		sort.Strings(imports)

		lines := []string{}
		lines = append(lines, `import (`)
		for _, theImport := range imports {
			alias := self.imports[theImport]
			if alias != "" {
				lines = append(lines, fmt.Sprintf(`	%s "%s"`, alias, theImport))
			} else {
				lines = append(lines, fmt.Sprintf(`	"%s"`, theImport))
			}
		}
		lines = append(lines, `)`)
		return lines
	}
	return []string{}
}
