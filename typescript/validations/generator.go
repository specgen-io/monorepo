package validations

import (
	"fmt"

	"generator"
	"spec"
	"typescript/modules"
	iots2 "typescript/validations/iots"
	superstruct2 "typescript/validations/superstruct"
)

type Validation interface {
	RuntimeTypeSamePackage(typ *spec.TypeDef) string
	RuntimeType(typ *spec.TypeDef) string
	SetupLibrary(validationModule modules.Module) *generator.CodeFile
	Models(models []*spec.NamedModel, validationModule modules.Module, module modules.Module) *generator.CodeFile
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
