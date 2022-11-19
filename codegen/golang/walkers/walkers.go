package walkers

import (
	"spec"
)

func ApiHasNonEmptyBody(api *spec.Api) bool {
	hasNonEmptyBody := false
	walk := spec.NewWalker().
		OnOperation(func(operation *spec.NamedOperation) {
			if operation.BodyIs(spec.BodyJson) || operation.BodyIs(spec.BodyString) {
				hasNonEmptyBody = true
			}
		})
	walk.Api(api)
	return hasNonEmptyBody
}

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

func ApiHasType(api *spec.Api, typName string) bool {
	foundType := false
	walk := spec.NewWalker().
		OnTypeDef(func(typ *spec.TypeDef) {
			if typ.Plain != typName {
				foundType = true
			}
		})
	walk.Api(api)
	return foundType
}
