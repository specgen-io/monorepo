package validations

import (
	"fmt"
	"typescript/validations/modules"
	"typescript/writer"

	"generator"
	"spec"
	"typescript/validations/iots"
	"typescript/validations/superstruct"
)

type Validation interface {
	RuntimeTypeName(typeName string) string
	RequestBodyJsonRuntimeType(body *spec.RequestBody) string
	ResponseBodyJsonRuntimeType(body *spec.ResponseBody) string
	SetupLibrary() *generator.CodeFile
	Models(version *spec.Version) *generator.CodeFile
	ErrorModels(httpErrors *spec.HttpErrors) *generator.CodeFile
	WriteParamsType(w *writer.Writer, typeName string, params []spec.NamedParam)
}

func New(validation string, modules *modules.Modules) Validation {
	if validation == superstruct.Superstruct {
		return &superstruct.Generator{modules}
	}
	if validation == iots.IoTs {
		return &iots.Generator{modules}
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}
