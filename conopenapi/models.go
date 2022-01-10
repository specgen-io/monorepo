package conopenapi

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/specgen-io/specgen/v2/spec"
)

func convertModels(schemas map[string]*openapi3.SchemaRef) []spec.NamedModel {
	models := []spec.NamedModel{}
	if schemas == nil {
		return models
	}
	for schemaName, schema := range schemas {
		models = append(models, spec.NamedModel{name(schemaName), convertSchema(schema), nil})
	}
	return models
}

func contains(array []string, what string) bool {
	for _, item := range array {
		if item == what {
			return true
		}
	}
	return false
}

func convertSchema(schema *openapi3.SchemaRef) spec.Model {
	switch schema.Value.Type {
	case TypeObject:
		fields := []spec.NamedDefinition{}
		for propName, propSchema := range schema.Value.Properties {
			fieldRequired := contains(schema.Value.Required, propName)
			fields = append(fields, spec.NamedDefinition{name(propName), *convertProperty(propSchema, fieldRequired)})
		}
		return spec.Model{&spec.Object{fields}, nil, nil, nil, nil}
	default:
		panic(fmt.Sprintf(`schema is not supported`))
	}
}

func convertProperty(schema *openapi3.SchemaRef, required bool) *spec.Definition {
	var description *string = nil
	if schema.Value.Description != "" {
		description = &schema.Value.Description
	}
	return &spec.Definition{*specType(schema, required), description, nil}
}
