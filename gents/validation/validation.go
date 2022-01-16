package validation

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

type Validation interface {
	RuntimeType(typ *spec.TypeDef) string
	RuntimeTypeFromPackage(customTypesPackage string, typ *spec.TypeDef) string
	SetupLibrary(validationModule modules.Module) *sources.CodeFile
	GenerateVersionModels(version *spec.Version, validationModule modules.Module, module modules.Module) *sources.CodeFile
	GenerateParams(w *sources.Writer, typeName string, params []spec.NamedParam)
}

func New(validation string) Validation {
	if validation == Superstruct {
		return &superstructValidation{}
	}
	if validation == IoTs {
		return &ioTsValidation{}
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}
