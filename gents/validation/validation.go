package validation

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func RuntimeType(validation string, typ *spec.TypeDef) string {
	return RuntimeTypeFromPackage(validation, "", typ)
}

func RuntimeTypeFromPackage(validation string, customTypesPackage string, typ *spec.TypeDef) string {
	if validation == Superstruct {
		return SuperstructTypeFromPackage(customTypesPackage, typ)
	}
	if validation == IoTs {
		return IoTsTypeFromPackage(customTypesPackage, typ)
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}

func GenerateValidation(validation string, validationModule modules.Module) *sources.CodeFile {
	if validation == Superstruct {
		return generateSuperstructStaticCode(validationModule)
	}
	if validation == IoTs {
		return generateIoTsStaticCode(validationModule)
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}
