package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func runtimeType(validation string, typ *spec.TypeDef) string {
	if validation == Superstruct {
		return SuperstructType(typ)
	}
	if validation == IoTs {
		return IoTsType(typ)
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}

func generateValidation(validation string, validationModule module) *sources.CodeFile {
	if validation == Superstruct {
		return generateSuperstructStaticCode(validationModule)
	}
	if validation == IoTs {
		return generateIoTsStaticCode(validationModule)
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}
