package gengo

import "github.com/specgen-io/spec"

func generateImports(api *spec.Api, imports []string) []string {
	if apiHasType(api, spec.TypeDate) {
		imports = append(imports, `"cloud.google.com/go/civil"`)
	}
	if apiHasType(api, spec.TypeUuid) {
		imports = append(imports, `"github.com/google/uuid"`)
	}
	if apiHasType(api, spec.TypeDecimal) {
		imports = append(imports, `"github.com/shopspring/decimal"`)
	}
	return imports
}

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
