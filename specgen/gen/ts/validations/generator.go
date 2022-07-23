package validations

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/ts/modules"
	iots2 "github.com/specgen-io/specgen/v2/gen/ts/validations/iots"
	superstruct2 "github.com/specgen-io/specgen/v2/gen/ts/validations/superstruct"
	"github.com/specgen-io/specgen/v2/generator"
	"github.com/specgen-io/specgen/v2/spec"
)

type Validation interface {
	RuntimeType(typ *spec.TypeDef) string
	RuntimeTypeFromPackage(customTypesPackage string, typ *spec.TypeDef) string
	SetupLibrary(validationModule modules.Module) *generator.CodeFile
	VersionModels(version *spec.Version, validationModule modules.Module, module modules.Module) *generator.CodeFile
	WriteParamsType(w *generator.Writer, typeName string, params []spec.NamedParam)
}

func New(validation string) Validation {
	if validation == superstruct2.Superstruct {
		return &superstruct2.Generator{}
	}
	if validation == iots2.IoTs {
		return &iots2.Generator{}
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}
