package gengo

import (
	"fmt"
	"github.com/specgen-io/spec"
)

//TODO: Tweak function names in the code below.
// For example apiHasType receives version as a parameter therefore it should be named "versionHasType".
func apiHasType(version *spec.Version, typ string) bool {
	for _, api := range version.Http.Apis {
		if operationHasType(&api, typ) {
			return true
		}
	}
	return false
}

func operationHasType(api *spec.Api, typ string) bool {
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

//TODO: Strage method name - there's no type provided
func bodyHasType(api *spec.Api) bool {
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

//TODO: What about unions?
func versionHasType(version *spec.Version, typ string) bool {
	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			for _, field := range model.Object.Fields {
				if checkType(&field.Type.Definition, typ) {
					return true
				}
			}
		}
	}
	return false
}
