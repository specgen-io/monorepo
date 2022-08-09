package imports

import (
	"sort"

	"generator"
	"golang/types"
	"spec"
)

type imports struct {
	imports map[string]string
}

func New() *imports {
	return &imports{imports: make(map[string]string)}
}

func (self *imports) Add(theImport string) *imports {
	self.imports[theImport] = ""
	return self
}

func (self *imports) AddAlias(theImport string, alias string) *imports {
	self.imports[theImport] = alias
	return self
}

func (self *imports) Write(w *generator.Writer) {
	if len(self.imports) > 0 {
		imports := make([]string, 0, len(self.imports))
		for theImport := range self.imports {
			imports = append(imports, theImport)
		}
		sort.Strings(imports)

		w.EmptyLine()
		w.Line(`import (`)
		for _, theImport := range imports {
			alias := self.imports[theImport]
			if alias != "" {
				w.Line(`  %s "%s"`, alias, theImport)
			} else {
				w.Line(`  "%s"`, theImport)
			}
		}
		w.Line(`)`)
	}
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

func (self *imports) AddModelsTypes(version *spec.Version) *imports {
	self.Add("errors")
	self.Add("encoding/json")
	if types.VersionModelsHasType(version, spec.TypeDate) {
		self.Add("cloud.google.com/go/civil")
	}
	if types.VersionModelsHasType(version, spec.TypeUuid) {
		self.Add("github.com/google/uuid")
	}
	if types.VersionModelsHasType(version, spec.TypeDecimal) {
		self.Add("github.com/shopspring/decimal")
	}
	return self
}
