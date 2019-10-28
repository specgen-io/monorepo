package genopenapi

import (
	"github.com/ModaOperandi/spec"
	"gopkg.in/yaml.v2"
	"specgen/gen"
	"strings"
)

func GenerateSpecification(serviceFile string, outFile string) (err error) {
	spec, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return
	}

	openapi := generateOpenapi(spec)

	data, err := yaml.Marshal(openapi.Yaml)
	if err != nil {
		return
	}

	openApiFile := gen.TextFile{outFile, string(data)}

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
			Set("info", info.Yaml)

	paths := generateApis(spec.Apis)
	openapi.Set("paths", paths.Yaml)

	schemas := Map()
	for _, model := range spec.Models {
		schemas.Set(model.Name.Source, generateModel(model.Model).Yaml)
	}
	openapi.Set("components", Map().Set("schemas", schemas.Yaml).Yaml)

	return openapi
}

func generateApis(apis spec.Apis) *YamlMap {
	paths := Map()
	groups := Groups(apis)
	for _, group := range groups {
		path := Map()
		for _, o := range group.Operations {
			path.Set(strings.ToLower(o.Operation.Method), generateOperation(o).Yaml)
		}
		paths.Set(group.Url, path.Yaml)
	}
	return paths
}

func generateOperation(o Operation) *YamlMap {
	operationId := o.Api.Name.Source + o.Operation.Name.PascalCase()
	operation := Map()
	operation.Set("operationId", operationId)
	operation.Set("tags", Array().Add(o.Api.Name.Source).Yaml)
	if o.Operation.Description != nil {
		operation.Set("description", o.Operation.Description)
	}
	if o.Operation.Body != nil {
		body := o.Operation.Body
		request := Map()
		if body.Description != nil {
			request.Set("description", body.Description)
		}
		request.Set("required", !body.Type.IsNullable())
		request.Set("content", Map().Set("application/json", Map().Set("schema", OpenApiType(&body.Type, nil).Yaml).Yaml).Yaml)
		operation.Set("requestBody", request.Yaml)
	}

	parameters := Array()

	addParameters(parameters, "path", o.Operation.UrlParams)
	addParameters(parameters, "header", o.Operation.HeaderParams)
	addParameters(parameters, "query", o.Operation.QueryParams)

	if parameters.Length() > 0 {
		operation.Set("parameters", parameters.Yaml)
	}

	responses := Map()
	for _, r := range o.Operation.Responses {
		responses.Set(spec.HttpStatusCode(r.Name), generateResponse(r.Response).Yaml)
	}
	operation.Set("responses", responses.Yaml)
	return operation
}

func addParameters(parameters *YamlArray, in string, params []spec.NamedParam) {
	for _, p := range params {
		param :=
			Map().
				Set("in", in).
				Set("name", p.Name.Source).
				Set("required", !p.Type.IsNullable()).
				Set("schema", OpenApiType(&p.Type, p.Default).Yaml)
		if p.Description != nil {
			param.Set("description", *p.Description)
		}
		parameters.Add(param.Yaml)
	}
}

func generateResponse(r spec.Response) *YamlMap {
	response := Map()
	if r.Description != nil {
		response.Set("description", r.Description)
	}
	if !r.Type.IsEmpty() {
		response.Set("content", Map().Set("application/json", Map().Set("schema", OpenApiType(&r.Type, nil).Yaml).Yaml).Yaml)
	}
	return response
}

func generateModel(model spec.Model) *YamlMap {
	if model.IsObject() {
		return generateObjectModel(model)
	} else {
		return generateEnumModel(model)
	}
}

func generateObjectModel(model spec.Model) *YamlMap {
	schema := Map().Set("type", "object")

	if model.Object.Description != nil {
		schema.Set("description", model.Object.Description)
	}

	required := Array()
	for _, field := range model.Object.Fields {
		if !field.Type.IsNullable() {
			required.Add(field.Name.Source)
		}
	}

	if required.Length() > 0 {
		schema.Set("required", required.Yaml)
	}

	properties := Map()
	for _, field := range model.Object.Fields {
		property := OpenApiType(&field.Type, field.Default)
		if field.Description != nil {
			property.Set("description", field.Description)
		}
		properties.Set(field.Name.Source, property.Yaml)
	}
	schema.Set("properties", properties.Yaml)

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
	schema.Set("enum", openApiItems.Yaml)

	return schema
}
