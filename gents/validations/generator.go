package validations

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/gents/validations/iots"
	"github.com/specgen-io/specgen/v2/gents/validations/superstruct"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

type Validation interface {
	RuntimeType(typ *spec.TypeDef) string
	RuntimeTypeFromPackage(customTypesPackage string, typ *spec.TypeDef) string
	SetupLibrary(validationModule modules.Module) *sources.CodeFile
	VersionModels(version *spec.Version, validationModule modules.Module, module modules.Module) *sources.CodeFile
	WriteParamsType(w *sources.Writer, typeName string, params []spec.NamedParam)
}

func New(validation string) Validation {
	if validation == superstruct.Superstruct {
		return &superstruct.Generator{}
	}
	if validation == iots.IoTs {
		return &iots.Generator{}
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}
