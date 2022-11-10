package writer

import (
	"fmt"
	"golang/module"
	"golang/types"
	"sort"
	"spec"
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
				lines = append(lines, fmt.Sprintf(`  %s "%s"`, alias, theImport))
			} else {
				lines = append(lines, fmt.Sprintf(`  "%s"`, theImport))
			}
		}
		lines = append(lines, `)`)
		return lines
	}
	return []string{}
}

func (self *imports) AddApiTypes(api *spec.Api) *imports {
	if types.ApiHasType(api, spec.TypeDate) {
		self.Add("cloud.google.com/go/civil")
	}
	if types.ApiHasType(api, spec.TypeJson) {
		self.Add("encoding/json")
	}
	if types.ApiHasType(api, spec.TypeUuid) {
		self.Add("github.com/google/uuid")
	}
	if types.ApiHasType(api, spec.TypeDecimal) {
		self.Add("github.com/shopspring/decimal")
	}
	return self
}

func (self *imports) AddModelsTypes(models []*spec.NamedModel) *imports {
	self.Add("errors")
	self.Add("encoding/json")
	if types.VersionModelsHasType(models, spec.TypeDate) {
		self.Add("cloud.google.com/go/civil")
	}
	if types.VersionModelsHasType(models, spec.TypeUuid) {
		self.Add("github.com/google/uuid")
	}
	if types.VersionModelsHasType(models, spec.TypeDecimal) {
		self.Add("github.com/shopspring/decimal")
	}
	return self
}
