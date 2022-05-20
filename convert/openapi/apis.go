package openapi

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pinzolo/casee"
	"github.com/specgen-io/specgen/v2/spec"
)

func (c *Converter) apis(doc *openapi3.T) spec.Apis {
	apis := collectApis(doc)
	specApis := []spec.Api{}
	for _, api := range apis {
		specApis = append(specApis, *c.api(api))
	}
	return spec.Apis{nil, specApis, nil}
}

func (c *Converter) api(api *Api) *spec.Api {
	operations := []spec.NamedOperation{}
	for _, pathItem := range api.Items {
		operation := c.pathItem(&pathItem)
		operations = append(operations, spec.NamedOperation{name(casee.ToSnakeCase(pathItem.Operation.OperationID)), *operation, nil})
	}
	return &spec.Api{name(api.Name), operations, nil}
}

func (c *Converter) pathItem(pathItem *PathItem) *spec.Operation {
	pathParams := collectParams(pathItem.Operation.Parameters, "path")
	headerParams := collectParams(pathItem.Operation.Parameters, "header")
	queryParams := collectParams(pathItem.Operation.Parameters, "query")

	endpoint := spec.Endpoint{pathItem.Method, pathItem.Path, c.params(pathParams), nil}
	var description *string = nil
	if pathItem.Operation.Description != "" {
		description = &pathItem.Operation.Description
	}
	operation := spec.Operation{
		endpoint,
		description,
		c.params(headerParams),
		c.params(queryParams),
		c.requestBody(pathItem.Operation.RequestBody),
		c.responses(pathItem.Operation.Responses),
		[]spec.NamedResponse{},
		nil,
	}
	return &operation
}

func (c *Converter) responses(responses openapi3.Responses) []spec.NamedResponse {
	result := []spec.NamedResponse{}
	for status, response := range responses {
		statusName := "ok"
		if status != "default" {
			statusName = spec.HttpStatusName(status)
		}
		result = append(result, spec.NamedResponse{name(statusName), *c.response(response), nil})
	}
	return result
}

func (c *Converter) response(response *openapi3.ResponseRef) *spec.Definition {
	if response.Value == nil {
		return nil //TODO: not sure in this - what if ref is specified here
	}
	definition := &spec.Definition{emptyType, response.Value.Description, nil}
	for mediaType, media := range response.Value.Content {
		switch mediaType {
		case "application/json":
			definition = &spec.Definition{*specType(media.Schema, true), response.Value.Description, nil}
			break
		default:
			panic(fmt.Sprintf("Unsupported media type: %s", mediaType))
		}
	}
	return definition
}

func (c *Converter) params(parameters openapi3.Parameters) []spec.NamedParam {
	result := []spec.NamedParam{}
	for _, parameter := range parameters {
		result = append(result, *c.param(parameter))
	}
	return result
}

func (c *Converter) param(parameter *openapi3.ParameterRef) *spec.NamedParam {
	p := parameter.Value
	return &spec.NamedParam{
		name(p.Name),
		spec.DefinitionDefault{
			*specType(p.Schema, p.Required),
			nil,
			&p.Description,
			nil,
		},
	}
}

func (c *Converter) requestBody(body *openapi3.RequestBodyRef) *spec.Definition {
	if body == nil {
		return nil // this is fair - no body means nil definition
	}
	if body.Value == nil {
		return nil //TODO: not sure in this - what if ref is specified here
	}
	media := body.Value.Content.Get("application/json")
	if media == nil {
		return nil
	}
	//TODO: check if non-required body is allowed
	definition := spec.Definition{*specType(media.Schema, body.Value.Required), &body.Value.Description, nil}
	return &definition
}

func collectParams(parameters openapi3.Parameters, in string) openapi3.Parameters {
	result := openapi3.Parameters{}
	for _, parameter := range parameters {
		if parameter.Value.In == in {
			result = append(result, parameter)
		}
	}
	return result
}

func useTagsAsApis(doc *openapi3.T) bool {
	for _, pathItem := range doc.Paths {
		for _, operation := range pathItem.Operations() {
			if len(operation.Tags) != 1 {
				return false //TODO: message here
			}
		}
	}
	return true
}

func collectApis(doc *openapi3.T) []*Api {
	useTagsAsApis := useTagsAsApis(doc)
	defaultApiName := name(doc.Info.Title).FlatCase()
	apisMap := map[string]*Api{defaultApiName: {defaultApiName, []PathItem{}}}
	for path, pathItem := range doc.Paths {
		for method, operation := range pathItem.Operations() {
			apiName := defaultApiName
			if useTagsAsApis {
				apiName = operation.Tags[0]
			}
			if _, existing := apisMap[apiName]; !existing {
				apisMap[apiName] = &Api{apiName, []PathItem{}}
			}
			api := apisMap[apiName]
			pathItem := PathItem{path, method, operation}
			api.Items = append(api.Items, pathItem)
		}
	}
	result := []*Api{}
	for _, api := range apisMap {
		result = append(result, api)
	}
	return result
}

type Api struct {
	Name  string
	Items []PathItem
}

type PathItem struct {
	Path      string
	Method    string
	Operation *openapi3.Operation
}
