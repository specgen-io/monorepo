package types

import (
	"fmt"
	"spec"
)

func ModelsHasEnum(models []*spec.NamedModel) bool {
	for _, model := range models {
		if model.IsEnum() {
			return true
		}
	}
	return false
}

func ApiHasType(api *spec.Api, typ string) bool {
	for _, operation := range api.Operations {
		if OperationHasType(&operation, typ) {
			return true
		}
	}
	return false
}

func OperationHasType(operation *spec.NamedOperation, typ string) bool {
	if paramHasType(operation.QueryParams, typ) {
		return true
	}
	if paramHasType(operation.HeaderParams, typ) {
		return true
	}
	if paramHasType(operation.Endpoint.UrlParams, typ) {
		return true
	}
	if operation.Body != nil {
		if checkType(&operation.Body.Type.Definition, typ) {
			return true
		}
	}
	for _, response := range operation.Responses {
		if checkType(&response.Type.Definition, typ) {
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

func VersionModelsHasType(models []*spec.NamedModel, typ string) bool {
	for _, model := range models {
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

func checkType(typedef *spec.TypeDef, typ string) bool {
	if typedef.Plain != typ {
		switch typedef.Node {
		case spec.PlainType:
			return false
		case spec.NullableType:
		case spec.ArrayType:
		case spec.MapType:
			return checkType(typedef.Child, typ)
		default:
			panic(fmt.Sprintf("Unknown type: %v", typ))
		}
	}
	return true
}
