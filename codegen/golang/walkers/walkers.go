package walkers

import (
	"spec"
)

func ApiIsUsingModels(api *spec.Api) bool {
	foundModels := false
	walk := spec.NewWalker().
		OnTypeDef(func(typ *spec.TypeDef) {
			if typ.Info.Model != nil && typ.Info.Model.InVersion != nil {
				foundModels = true
			}
		})
	walk.Api(api)
	return foundModels
}

func ApiIsUsingErrorModels(api *spec.Api) bool {
	foundErrorModels := false
	walk := spec.NewWalker().
		OnTypeDef(func(typ *spec.TypeDef) {
			if typ.Info.Model != nil && typ.Info.Model.InHttpErrors != nil {
				foundErrorModels = true
			}
		})
	walk.Api(api)
	return foundErrorModels
}

func ApiHasNonSingleResponse(api *spec.Api) bool {
	hasNonSingleResponse := false
	walk := spec.NewWalker().
		OnOperation(func(operation *spec.NamedOperation) {
			if len(operation.Responses) > 1 {
				hasNonSingleResponse = true
			}
		})
	walk.Api(api)
	return hasNonSingleResponse
}

func ApiHasUrlParams(api *spec.Api) bool {
	hasUrlParams := false
	walk := spec.NewWalker().
		OnOperation(func(operation *spec.NamedOperation) {
			if operation.Endpoint.UrlParams != nil && len(operation.Endpoint.UrlParams) > 0 {
				hasUrlParams = true
			}
		})
	walk.Api(api)
	return hasUrlParams
}

func ApiHasHasHeaderParams(api *spec.Api) bool {
	hasHeaderParams := false
	walk := spec.NewWalker().
		OnOperation(func(operation *spec.NamedOperation) {
			if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
				hasHeaderParams = true
			}
		})
	walk.Api(api)
	return hasHeaderParams
}

func OperationHasHeaderParams(operation *spec.NamedOperation) bool {
	hasHeaderParams := false
	walk := spec.NewWalker().
		OnOperation(func(operation *spec.NamedOperation) {
			if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
				hasHeaderParams = true
			}
		})
	walk.Operation(operation)
	return hasHeaderParams
}

func ApiHasBodyOfKind(api *spec.Api, kinds ...spec.BodyKind) bool {
	result := false
	walk := spec.NewWalker().
		OnRequestBody(func(body *spec.RequestBody) {
			for _, kind := range kinds {
				if body.Is(kind) {
					result = true
				}
			}
		}).
		OnResponseBody(func(body *spec.ResponseBody) {
			for _, kind := range kinds {
				if body.Is(kind) {
					result = true
				}
			}
		})
	walk.Api(api)
	return result
}

func ApiHasMultiResponsesWithEmptyBody(api *spec.Api) bool {
	result := false
	walk := spec.NewWalker().
		OnOperationResponse(func(response *spec.OperationResponse) {
			if len(response.Operation.Responses) > 1 && response.Body.IsEmpty() {
				result = true
			}
		})
	walk.Api(api)
	return result
}

func ApiHasMultiSuccessResponsesWithEmptyBody(api *spec.Api) bool {
	result := false
	walk := spec.NewWalker().
		OnOperationResponse(func(response *spec.OperationResponse) {
			if len(response.Operation.Responses.Success()) > 1 && response.Body.IsEmpty() {
				result = true
			}
		})
	walk.Api(api)
	return result
}

func ApiHasType(api *spec.Api, typName string) bool {
	foundType := false
	walk := spec.NewWalker().
		OnTypeDef(func(typ *spec.TypeDef) {
			if typ.Plain == typName {
				foundType = true
			}
		})
	walk.Api(api)
	return foundType
}

func ModelHasType(model *spec.NamedModel, typName string) bool {
	foundType := false
	walk := spec.NewWalker().
		OnTypeDef(func(typ *spec.TypeDef) {
			if typ.Plain == typName {
				foundType = true
			}
		})
	walk.Model(model)
	return foundType
}
