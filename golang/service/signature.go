package service

import (
	"fmt"
	"golang/common"
	"strings"

	"golang/responses"
	"golang/types"
	"spec"
)

func OperationSignature(operation *spec.NamedOperation, apiPackage *string) string {
	return fmt.Sprintf(`%s(%s) %s`,
		operation.Name.PascalCase(),
		strings.Join(common.OperationParams(operation), ", "),
		operationReturn(operation, apiPackage),
	)
}

func operationReturn(operation *spec.NamedOperation, responsePackageName *string) string {
	if len(operation.Responses) == 1 {
		response := operation.Responses[0]
		if response.Type.Definition.IsEmpty() {
			return `error`
		}
		return fmt.Sprintf(`(*%s, error)`, types.GoType(&response.Type.Definition))
	}
	responseType := responses.ResponseTypeName(operation)
	if responsePackageName != nil {
		responseType = *responsePackageName + "." + responseType
	}
	return fmt.Sprintf(`(*%s, error)`, responseType)
}