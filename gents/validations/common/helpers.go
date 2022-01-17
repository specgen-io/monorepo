package common

import "fmt"

func ParamsRuntimeTypeName(typeName string) string {
	return fmt.Sprintf("T%s", typeName)
}
