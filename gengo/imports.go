package gengo

import (
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
	"sort"
)

type imports struct {
	imports map[string]string
}

func Imports() *imports {
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

func (self *imports) Write(w *gen.Writer) {
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
	if apiHasType(api, spec.TypeDate) {
		self.Add("cloud.google.com/go/civil")
	}
	if apiHasType(api, spec.TypeJson) {
		self.Add("encoding/json")
	}
	if apiHasType(api, spec.TypeUuid) {
		self.Add("github.com/google/uuid")
	}
	if apiHasType(api, spec.TypeDecimal) {
		self.Add("github.com/shopspring/decimal")
	}
	return self
}

func (self *imports) AddModelsTypes(version *spec.Version) *imports {
	if versionModelsHasType(version, spec.TypeDate) {
		self.Add("cloud.google.com/go/civil")
	}
	if versionModelsHasType(version, spec.TypeJson) {
		self.Add("encoding/json")
	}
	if versionModelsHasType(version, spec.TypeUuid) {
		self.Add("github.com/google/uuid")
	}
	if versionModelsHasType(version, spec.TypeDecimal) {
		self.Add("github.com/shopspring/decimal")
	}
	if isOneOfModel(version) {
		self.Add("errors")
		self.Add("fmt")
	}
	return self
}
