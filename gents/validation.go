package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func importEncoding(validation string) string {
	if validation == Superstruct {
		return importSuperstructEncoding
	}
	if validation == IoTs {
		return importIoTsEncoding
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}

func runtimeType(validation string, typ *spec.TypeDef) string {
	if validation == Superstruct {
		return SuperstructType(typ)
	}
	if validation == IoTs {
		return IoTsType(typ)
	}
	panic(fmt.Sprintf("Unknown validation: %s", validation))
}
