package openapi

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/specgen-io/specgen/spec/v2"
)

func (c *Converter) models(schemas map[string]*openapi3.SchemaRef) []spec.NamedModel {
	models := []spec.NamedModel{}
	if schemas == nil {
		return models
	}
	for schemaName, schema := range schemas {
		specModel := c.schema(schema)
		if specModel != nil {
			models = append(models, spec.NamedModel{name(schemaName), *specModel, nil})
		}
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

func (c *Converter) schema(schema *openapi3.SchemaRef) *spec.Model {
	switch schema.Value.Type {
	case TypeObject:
		fields := []spec.NamedDefinition{}
		for propName, propSchema := range schema.Value.Properties {
			fieldRequired := contains(schema.Value.Required, propName)
			fields = append(fields, spec.NamedDefinition{name(propName), *c.property(propSchema, fieldRequired)})
		}
		return &spec.Model{&spec.Object{fields}, nil, nil, nil, nil}
	default:
		panic(fmt.Sprintf(`schema is not supported yet`))
	}
}

func (c *Converter) property(schema *openapi3.SchemaRef, required bool) *spec.Definition {
	var description *string = nil
	if schema.Value.Description != "" {
		description = &schema.Value.Description
	}
	return &spec.Definition{*specType(schema, required), description, nil}
}
