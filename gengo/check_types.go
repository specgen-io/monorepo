package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func versionHasType(version *spec.Version, typ string) bool {
	for _, api := range version.Http.Apis {
		if apiHasType(&api, typ) {
			return true
		}
	}
	return false
}

func apiHasType(api *spec.Api, typ string) bool {
	for _, operation := range api.Operations {
		if paramHasType(operation.QueryParams, typ) {
			return true
		}
		if paramHasType(operation.HeaderParams, typ) {
			return true
		}
		if paramHasType(operation.Endpoint.UrlParams, typ) {
			return true
		}
	}
	return false
}

func apiHasBody(api *spec.Api) bool {
	for _, operation := range api.Operations {
		if operation.Body != nil {
			return true
		}
	}
	return false
}

func paramHasType(namedParams []spec.NamedParam, typ string) bool {
	if namedParams != nil && len(namedParams) > 0 {
		for _, param := range namedParams {
			if checkType(&param.Type.Definition, typ) {
				return true
			}
		}
	}
	return false
}

func checkType(fieldType *spec.TypeDef, typ string) bool {
	if fieldType.Plain != typ {
		switch fieldType.Node {
		case spec.PlainType:
			return false
		case spec.NullableType:
			return checkType(fieldType.Child, typ)
		case spec.ArrayType:
			return checkType(fieldType.Child, typ)
		case spec.MapType:
			return checkType(fieldType.Child, typ)
		default:
			panic(fmt.Sprintf("Unknown type: %v", typ))
		}
	}
	return true
}

func versionModelsHasType(version *spec.Version, typ string) bool {
	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			for _, field := range model.Object.Fields {
				if checkType(&field.Type.Definition, typ) {
					return true
				}
			}
		}
		if model.IsOneOf() {
			for _, item := range model.OneOf.Items {
				if checkType(&item.Type.Definition, typ) {
					return true
				}
			}
		}
	}
	return false
}
