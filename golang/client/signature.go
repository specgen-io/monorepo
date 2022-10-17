package client

import (
	"fmt"
	"golang/common"
	"golang/types"
	"strings"

	"golang/responses"
	"spec"
)

func OperationSignature(operation *spec.NamedOperation, types *types.Types, apiPackage *string) string {
	return fmt.Sprintf(`%s(%s) %s`,
		operation.Name.PascalCase(),
		strings.Join(common.OperationParams(types, operation), ", "),
		operationReturn(operation, types, apiPackage),
	)
}

func operationReturn(operation *spec.NamedOperation, types *types.Types, responsePackageName *string) string {
	if common.ResponsesNumber(operation) == 1 {
		response := operation.Responses[0]
		return fmt.Sprintf(`(*%s, error)`, types.GoType(&response.Type.Definition))
	}
	responseType := responses.ResponseTypeName(operation)
	if responsePackageName != nil {
		responseType = *responsePackageName + "." + responseType
	}
	return fmt.Sprintf(`(*%s, error)`, responseType)
}
