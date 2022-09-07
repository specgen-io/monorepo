package openapi

import (
	"strings"

	"generator"
	"github.com/pinzolo/casee"
	"spec"
	"yamlx"
)

func GenerateOpenapi(spec *spec.Spec, outFile string) *generator.CodeFile {
	openapi := generateSpecification(spec)
	data, _ := yamlx.ToYamlString(openapi)
	openapiFile := &generator.CodeFile{outFile, data}
	return openapiFile
}

func generateSpecification(spec *spec.Spec) *yamlx.YamlMap {
	info := yamlx.Map()
	title := spec.Name.Source
	if spec.Title != nil {
		title = *spec.Title
	}
	info.Add("title", title)
	if spec.Description != nil {
		info.Add("description", spec.Description)
	}
	info.Add("version", spec.Version)

	schemas := yamlx.Map()
	for _, version := range spec.Versions {
		for _, model := range version.Models {
			schemas.Merge(generateModel(&model).Node)
		}
	}

	components := yamlx.Map(yamlx.Pair{"schemas", schemas})

	openapi := yamlx.Map(
		yamlx.Pair{"openapi", "3.0.0"},
		yamlx.Pair{"info", info},
		yamlx.Pair{"paths", generateApis(spec)},
		yamlx.Pair{"components", components},
	)

	return openapi
}

func versionedModelName(version *spec.Version, modelName string) string {
	if version != nil && version.Name.Source != "" {
		return version.Name.Source + "." + modelName
	}
	return modelName
}

func generateApis(spec *spec.Spec) *yamlx.YamlMap {
	paths := yamlx.Map()
	groups := OperationsByUrl(spec)
	for _, group := range groups {
		path := yamlx.Map()
		for _, o := range group.Operations {
			path.Add(strings.ToLower(o.Operation.Endpoint.Method), generateOperation(o))
		}
		paths.Add(group.Url, path)
	}
	return paths
}

func generateOperation(o *spec.NamedOperation) *yamlx.YamlMap {
	version := o.InApi.InHttp.InVersion.Name
	operationId := casee.ToCamelCase(version.PascalCase() + o.InApi.Name.PascalCase() + o.Name.PascalCase())
	operation := yamlx.Map()
	operation.Add("operationId", operationId)
	operation.Add("tags", yamlx.Array(o.InApi.Name.Source))

	if o.Operation.Description != nil {
		operation.Add("description", o.Operation.Description)
	}
	if o.Operation.Body != nil {
		body := o.Operation.Body
		request := yamlx.Map()
		if body.Description != nil {
			request.Add("description", body.Description)
		}
		request.Add("required", !body.Type.Definition.IsNullable())
		request.Add("content", generateJsonContent(&body.Type.Definition))
		operation.Add("requestBody", request)
	}

	parameters := yamlx.Array()

	addParameters(parameters, "path", o.Operation.Endpoint.UrlParams)
	addParameters(parameters, "header", o.Operation.HeaderParams)
	addParameters(parameters, "query", o.Operation.QueryParams)

	if parameters.Length() > 0 {
		operation.Add("parameters", parameters)
	}

	operation.Add("responses", generateResponses(o))
	return operation
}

func addParameters(parameters *yamlx.YamlArray, in string, params []spec.NamedParam) {
	for _, p := range params {
		param := yamlx.Map()
		param.Add("in", in)
		param.Add("name", p.Name.Source)
		param.Add("required", !p.Type.Definition.IsNullable())
		schema := OpenApiType(&p.Type.Definition)
		if p.Default != nil {
			schema.AddRaw("default", *p.Default)
		}
		param.Add("schema", schema)
		if p.Description != nil {
			param.Add("description", *p.Description)
		}
		parameters.Add(param)
	}
}

func generateJsonContent(types ...*spec.TypeDef) *yamlx.YamlMap {
	schema := yamlx.Map(yamlx.Pair{"schema", OpenApiType(types...)})
	content := yamlx.Map(yamlx.Pair{"application/json", schema})
	return content
}

func generateResponses(operation *spec.NamedOperation) *yamlx.YamlMap {
	statusCodes := operation.Responses.HttpStatusCodes()
	if operation.InApi.InHttp.InVersion.InSpec.HttpErrors != nil {
		statusCodes = spec.MergeHttpStatusCodes(statusCodes, operation.InApi.InHttp.InVersion.InSpec.HttpErrors.Responses.HttpStatusCodes())
	}

	result := yamlx.Map()
	for _, statusCode := range statusCodes {
		response := operation.Responses.GetByStatusCode(statusCode)
		errorResponse := operation.InApi.InHttp.InVersion.InSpec.HttpErrors.Responses.GetByStatusCode(statusCode)
		var responseDefinition *spec.Definition = nil
		var alternateDefinition *spec.Definition = nil

		if response != nil {
			responseDefinition = &response.Definition
			if errorResponse != nil && response.Definition.Type.String() != errorResponse.Definition.Type.String() {
				alternateDefinition = &errorResponse.Definition
			}
		} else {
			responseDefinition = &errorResponse.Definition
		}

		result.Add(statusCode, generateResponse(responseDefinition, alternateDefinition))
	}

	return result
}

func generateResponse(response *spec.Definition, alternate *spec.Definition) *yamlx.YamlMap {
	result := yamlx.Map()
	description := ""
	if response.Description != nil {
		description = *response.Description
	}
	result.Add("description", description)

	types := []*spec.TypeDef{}
	if !response.Type.Definition.IsEmpty() {
		types = append(types, &response.Type.Definition)
	}
	if alternate != nil && !alternate.Type.Definition.IsEmpty() {
		types = append(types, &alternate.Type.Definition)
	}

	if len(types) > 0 {
		result.Add("content", generateJsonContent(types...))
	}
	return result
}

func generateModel(model *spec.NamedModel) *yamlx.YamlMap {
	result := yamlx.Map()
	if model.IsObject() {
		return generateObjectModel(model)
	} else if model.IsEnum() {
		return generateEnumModel(model)
	} else if model.IsOneOf() {
		if model.OneOf.Discriminator != nil {
			return generateOneOfDiscriminatorModel(model)
		} else {
			return generateOneOfWrapperModel(model)
		}
	}
	return result
}

func generateOneOfWrapperModel(model *spec.NamedModel) *yamlx.YamlMap {
	oneOfItems := yamlx.Array()

	for _, item := range model.OneOf.Items {
		itemSchema := yamlx.Map()
		itemSchema.Add("type", "object")
		if item.Description != nil {
			itemSchema.Add("description", item.Description)
		}
		itemSchema.Add("required", yamlx.Array(item.Name.Source))
		itemSchema.Add("properties", yamlx.Map(yamlx.Pair{item.Name.Source, OpenApiType(&item.Type.Definition)}))
		oneOfItems.Add(itemSchema)
	}

	schema := yamlx.Map()
	if model.Description != nil {
		schema.Add("description", model.Description)
	}
	schema.Add("oneOf", oneOfItems)
	result := yamlx.Map()
	result.Add(versionedModelName(model.InVersion, model.Name.Source), schema)
	return result
}

func itemTypeName(model *spec.NamedModel, item *spec.NamedDefinition) string {
	return versionedModelName(model.InVersion, model.Name.Source+item.Name.PascalCase())
}

func generateOneOfDiscriminatorModel(model *spec.NamedModel) *yamlx.YamlMap {
	schema := yamlx.Map()

	if model.Description != nil {
		schema.Add("description", model.Description)
	}

	anyOfItems := yamlx.Array()
	for _, item := range model.OneOf.Items {
		anyOfItems.Add(yamlx.Map(yamlx.Pair{"$ref", componentSchemas(itemTypeName(model, &item))}))
	}
	schema.Add("anyOf", anyOfItems)

	discriminator := yamlx.Map()
	discriminator.Add("propertyName", *model.OneOf.Discriminator)

	mapping := yamlx.Map()
	for _, item := range model.OneOf.Items {
		mapping.Add(item.Name.Source, componentSchemas(itemTypeName(model, &item)))
	}
	discriminator.Add("mapping", mapping)

	schema.Add("discriminator", discriminator)

	result := yamlx.Map()
	result.Add(versionedModelName(model.InVersion, model.Name.Source), schema)
	for _, item := range model.OneOf.Items {
		result.Add(itemTypeName(model, &item), generateOneOfDiscriminatorItemModel(model, &item))
	}
	return result
}

func generateOneOfDiscriminatorItemModel(model *spec.NamedModel, item *spec.NamedDefinition) *yamlx.YamlMap {
	allOfItems := yamlx.Array()
	allOfItems.Add(OpenApiType(&item.Type.Definition))

	discriminatorField := yamlx.Map()
	discriminatorField.Add("type", "object")
	discriminatorField.Add("required", yamlx.Array(model.OneOf.Discriminator))
	discriminatorField.Add("properties", yamlx.Map(yamlx.Pair{
		model.OneOf.Discriminator, yamlx.Map(yamlx.Pair{
			"type", "string",
		}),
	}))
	allOfItems.Add(discriminatorField)

	schema := yamlx.Map()
	if item.Description != nil {
		schema.Add("description", item.Description)
	}
	schema.Add("allOf", allOfItems)
	return schema
}

func generateObjectModel(model *spec.NamedModel) *yamlx.YamlMap {
	schema := yamlx.Map()
	schema.Add("type", "object")

	if model.Description != nil {
		schema.Add("description", model.Description)
	}

	required := yamlx.Array()
	for _, field := range model.Object.Fields {
		if !field.Type.Definition.IsNullable() {
			required.Add(field.Name.Source)
		}
	}

	if required.Length() > 0 {
		schema.Add("required", required)
	}

	properties := yamlx.Map()
	for _, field := range model.Object.Fields {
		property := OpenApiType(&field.Type.Definition)
		if field.Description != nil {
			property.Add("description", field.Description)
		}
		properties.Add(field.Name.Source, property)
	}
	schema.Add("properties", properties)

	result := yamlx.Map()
	result.Add(versionedModelName(model.InVersion, model.Name.Source), schema)
	return result
}

func generateEnumModel(model *spec.NamedModel) *yamlx.YamlMap {
	schema := yamlx.Map()
	schema.Add("type", "string")

	if model.Description != nil {
		schema.Add("description", model.Description)
	}

	openApiItems := yamlx.Array()
	for _, item := range model.Enum.Items {
		openApiItems.Add(item.Name.Source)
	}
	schema.Add("enum", openApiItems)

	result := yamlx.Map()
	result.Add(versionedModelName(model.InVersion, model.Name.Source), schema)
	return result
}
