package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
	"path/filepath"
	"sort"
)

func generateVersionImports(version *spec.Version, imports []string) []string {
	if versionHasType(version, spec.TypeDate) {
		imports = append(imports, `"cloud.google.com/go/civil"`)
	}
	if versionHasType(version, spec.TypeJson) {
		imports = append(imports, `"encoding/json"`)
	}
	if versionHasType(version, spec.TypeUuid) {
		imports = append(imports, `"github.com/google/uuid"`)
	}
	if versionHasType(version, spec.TypeDecimal) {
		imports = append(imports, `"github.com/shopspring/decimal"`)
	}
	return imports
}

func generateApiImports(api *spec.Api, imports []string) []string {
	if apiHasType(api, spec.TypeDate) {
		imports = append(imports, `"cloud.google.com/go/civil"`)
	}
	if apiHasType(api, spec.TypeJson) {
		imports = append(imports, `"encoding/json"`)
	}
	if apiHasType(api, spec.TypeUuid) {
		imports = append(imports, `"github.com/google/uuid"`)
	}
	if apiHasType(api, spec.TypeDecimal) {
		imports = append(imports, `"github.com/shopspring/decimal"`)
	}
	return imports
}

type Package struct {
	root string
	packageName string
}

func NewPackage (root string, packageName string) Package {
	return Package{root, packageName}
}

func (p Package) Import() string {
	return fmt.Sprintf(`"%s/%s"`, p.root, p.packageName)
}

func (p Package) ImportAlias(alias string) string {
	return fmt.Sprintf(`%s %s`, alias, p.Import())
}

func (p Package) Path(filename string) string {
	return filepath.Join(p.packageName, filename)
}

func sortImports(imports []string) {
	sort.Strings(imports)
}