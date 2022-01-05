package gents

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func runtimeType(validation string, typ *spec.TypeDef) string {
	return runtimeTypeFromPackage(validation, "", typ)
}

func runtimeTypeFromPackage(validation string, customTypesPackage string, typ *spec.TypeDef) string {
	if validation == Superstruct {
		return SuperstructTypeFromPackage(customTypesPackage, typ)
	}
	if validation == IoTs {
		return IoTsTypeFromPackage(customTypesPackage, typ)
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
