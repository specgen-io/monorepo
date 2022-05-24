package imports

import (
	"github.com/specgen-io/specgen/v2/gen/golang/common"
	"github.com/specgen-io/specgen/v2/sources"
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

func (self *imports) Write(w *sources.Writer) {
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
	if common.ApiHasType(api, spec.TypeDate) {
		self.Add("cloud.google.com/go/civil")
	}
	if common.ApiHasType(api, spec.TypeJson) {
		self.Add("encoding/json")
	}
	if common.ApiHasType(api, spec.TypeUuid) {
		self.Add("github.com/google/uuid")
	}
	if common.ApiHasType(api, spec.TypeDecimal) {
		self.Add("github.com/shopspring/decimal")
	}
	return self
}

func (self *imports) AddModelsTypes(version *spec.Version) *imports {
	self.Add("errors")
	self.Add("encoding/json")
	if common.VersionModelsHasType(version, spec.TypeDate) {
		self.Add("cloud.google.com/go/civil")
	}
	if common.VersionModelsHasType(version, spec.TypeUuid) {
		self.Add("github.com/google/uuid")
	}
	if common.VersionModelsHasType(version, spec.TypeDecimal) {
		self.Add("github.com/shopspring/decimal")
	}
	return self
}
