package genopenapi

import (
	"github.com/pinzolo/casee"
	spec "github.com/specgen-io/spec.v2"
	"specgen/gen"
	"strings"
)

func GenerateSpecification(serviceFile string, outFile string) (err error) {
	spec, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return
	}

	openapi := generateOpenapi(spec)

	data, err := ToYamlString(openapi)
	if err != nil {
		return
	}

	openApiFile := gen.TextFile{outFile, data}

	err = gen.WriteFile(&openApiFile, true)

	return
}

func generateOpenapi(spec *spec.Spec) *YamlMap {
	info := Map()
	if spec.Title != nil {
		info.Set("title", spec.Title)
	}
	if spec.Description != nil {
		info.Set("description", spec.Description)
	}
	info.Set("version", spec.Version)

	openapi :=
		Map().
			Set("openapi", "3.0.0").
			Set("info", info)

	paths := generateApis(spec)
	openapi.Set("paths", paths)

	schemas := Map()

	for _, version := range spec.Versions {
		for _, model := range version.Models {
			schemas.Set(versionedModelName(version.Version.Source, model.Name.Source), generateModel(model.Model))
		}
	}
	openapi.Set("components", Map().Set("schemas", schemas))

	return openapi
}

func versionedModelName(version string, modelName string) string {
	if version != "" {
		return version + "." + modelName
	}
	return modelName
}

func generateApis(spec *spec.Spec) *YamlMap {
	paths := Map()
	groups := OperationsByUrl(spec)
	for _, group := range groups {
		path := Map()
		for _, o := range group.Operations {
			path.Set(strings.ToLower(o.Operation.Endpoint.Method), generateOperation(o))
		}
		paths.Set(group.Url, path)
	}
	return paths
}

func generateOperation(o *spec.NamedOperation) *YamlMap {
	version := o.Api.Apis.Version.Version
	operationId := casee.ToCamelCase(version.PascalCase() + o.Api.Name.PascalCase() + o.Name.PascalCase())
	operation := Map().Set("operationId", operationId).Set("tags", Array().Add(o.Api.Name.Source))

	if o.Operation.Description != nil {
		operation.Set("description", o.Operation.Description)
	}
	if o.Operation.Body != nil {
		body := o.Operation.Body
		request := Map()
		if body.Description != nil {
			request.Set("description", body.Description)
		}
		request.Set("required", !body.Type.Definition.IsNullable())
		request.Set("content", Map().Set("application/json", Map().Set("schema", OpenApiType(&body.Type.Definition, nil))))
		operation.Set("requestBody", request)
	}

	parameters := Array()

	addParameters(parameters, "path", o.Operation.Endpoint.UrlParams)
	addParameters(parameters, "header", o.Operation.HeaderParams)
	addParameters(parameters, "query", o.Operation.QueryParams)

	if parameters.Length() > 0 {
		operation.Set("parameters", parameters)
	}

	responses := Map()
	for _, r := range o.Operation.Responses {
		responses.Set(spec.HttpStatusCode(r.Name), generateResponse(r.Definition))
	}
	operation.Set("responses", responses)
	return operation
}

func addParameters(parameters *YamlArray, in string, params []spec.NamedParam) {
	for _, p := range params {
		param :=
			Map().
				Set("in", in).
				Set("name", p.Name.Source).
				Set("required", !p.Type.Definition.IsNullable()).
				Set("schema", OpenApiType(&p.Type.Definition, p.Default))
		if p.Description != nil {
			param.Set("description", *p.Description)
		}
		parameters.Add(param)
	}
}

func generateResponse(r spec.Definition) *YamlMap {
	response := Map()
	description := ""
	if r.Description != nil {
		description = *r.Description
	}
	response.Set("description", description)
	if !r.Type.Definition.IsEmpty() {
		response.Set("content", Map().Set("application/json", Map().Set("schema", OpenApiType(&r.Type.Definition, nil))))
	}
	return response
}

func generateModel(model spec.Model) *YamlMap {
	if model.IsObject() {
		return generateObjectModel(model)
	} else if model.IsEnum() {
		return generateEnumModel(model)
	} else {
		return generateUnionModel(model)
	}
}

func generateUnionModel(model spec.Model) *YamlMap {
	schema := Map().Set("type", "object")

	if model.OneOf.Description != nil {
		schema.Set("description", model.OneOf.Description)
	}

	properties := Map()
	for _, item := range model.OneOf.Items {
		property := OpenApiType(&item.Type.Definition, nil)
		if item.Description != nil {
			property.Set("description", item.Description)
		}
		properties.Set(item.Name.Source, property)
	}
	schema.Set("properties", properties)

	return schema
}

func generateObjectModel(model spec.Model) *YamlMap {
	schema := Map().Set("type", "object")

	if model.Object.Description != nil {
		schema.Set("description", model.Object.Description)
	}

	required := Array()
	for _, field := range model.Object.Fields {
		if !field.Type.Definition.IsNullable() {
			required.Add(field.Name.Source)
		}
	}

	if required.Length() > 0 {
		schema.Set("required", required)
	}

	properties := Map()
	for _, field := range model.Object.Fields {
		property := OpenApiType(&field.Type.Definition, nil)
		if field.Description != nil {
			property.Set("description", field.Description)
		}
		properties.Set(field.Name.Source, property)
	}
	schema.Set("properties", properties)

	return schema
}

func generateEnumModel(model spec.Model) *YamlMap {
	schema := Map().Set("type", "string")

	if model.Enum.Description != nil {
		schema.Set("description", model.Enum.Description)
	}

	openApiItems := Array()
	for _, item := range model.Enum.Items {
		openApiItems.Add(item.Name.Source)
	}
	schema.Set("enum", openApiItems)

	return schema
}
