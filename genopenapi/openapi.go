package genopenapi

import (
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/spec"
	"github.com/specgen-io/specgen/v2/yamlx"
	"strings"
)

func GenerateAndWriteOpenapi(specification *spec.Spec, outFile string) (err error) {
	openApiFile := GenerateOpenapi(specification, outFile)
	err = gen.WriteFile(openApiFile, true)
	return
}

func GenerateOpenapi(spec *spec.Spec, outFile string) *gen.TextFile {
	openapi := generateSpecification(spec)

	data, err := yamlx.ToYamlString(openapi)
	if err != nil {
		panic(err)
	}

	return &gen.TextFile{outFile, data}
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
			schemas.Add(versionedModelName(version.Version.Source, model.Name.Source), generateModel(model.Model))
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

func versionedModelName(version string, modelName string) string {
	if version != "" {
		return version + "." + modelName
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
	version := o.Api.Apis.Version.Version
	operationId := casee.ToCamelCase(version.PascalCase() + o.Api.Name.PascalCase() + o.Name.PascalCase())
	operation := yamlx.Map()
	operation.Add("operationId", operationId)
	operation.Add("tags", yamlx.Array(o.Api.Name.Source))

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

	responses := yamlx.Map()
	allResponses := addSpecialResponses(&o.Operation)
	for _, r := range allResponses {
		responses.Add(spec.HttpStatusCode(r.Name), generateResponse(r.Definition))
	}
	operation.Add("responses", responses)
	return operation
}

func addParameters(parameters *yamlx.YamlArray, in string, params []spec.NamedParam) {
	for _, p := range params {
		param := yamlx.Map()
		param.Add("in", in)
		param.Add("name", p.Name.Source)
		param.Add("required", !p.Type.Definition.IsNullable())
		param.Add("schema", OpenApiType(&p.Type.Definition, p.Default))
		if p.Description != nil {
			param.Add("description", *p.Description)
		}
		parameters.Add(param)
	}
}

func generateJsonContent(typ *spec.TypeDef) *yamlx.YamlMap {
	schema := yamlx.Map(yamlx.Pair{"schema", OpenApiType(typ, nil)})
	content := yamlx.Map(yamlx.Pair{"application/json", schema})
	return content
}

func generateResponse(r spec.Definition) *yamlx.YamlMap {
	response := yamlx.Map()
	description := ""
	if r.Description != nil {
		description = *r.Description
	}
	response.Add("description", description)
	if !r.Type.Definition.IsEmpty() {
		response.Add("content", generateJsonContent(&r.Type.Definition))
	}
	return response
}

func generateModel(model spec.Model) *yamlx.YamlMap {
	if model.IsObject() {
		return generateObjectModel(model)
	} else if model.IsEnum() {
		return generateEnumModel(model)
	} else {
		return generateUnionModel(model)
	}
}

func generateUnionModel(model spec.Model) *yamlx.YamlMap {
	schema := yamlx.Map()
	schema.Add("type", "object")

	if model.Description != nil {
		schema.Add("description", model.Description)
	}

	properties := yamlx.Map()
	for _, item := range model.OneOf.Items {
		property := OpenApiType(&item.Type.Definition, nil)
		if item.Description != nil {
			property.Add("description", item.Description)
		}
		properties.Add(item.Name.Source, property)
	}
	schema.Add("properties", properties)

	return schema
}

func generateObjectModel(model spec.Model) *yamlx.YamlMap {
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
		property := OpenApiType(&field.Type.Definition, nil)
		if field.Description != nil {
			property.Add("description", field.Description)
		}
		properties.Add(field.Name.Source, property)
	}
	schema.Add("properties", properties)

	return schema
}

func generateEnumModel(model spec.Model) *yamlx.YamlMap {
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

	return schema
}
