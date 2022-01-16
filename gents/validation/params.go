package validation

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func GenerateParams(w *sources.Writer, typeName string, params []spec.NamedParam, validation string) {
	if validation == Superstruct {
		generateSuperstructParams(w, typeName, params)
		return
	}
	if validation == IoTs {
		generateIoTsParams(w, typeName, params)
		return
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}
