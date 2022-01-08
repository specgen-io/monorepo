package conopenapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/specgen-io/specgen/v2/spec"
)

func ConvertFromOpenapi(inFile, outFile string) error {
	doc, err := openapi3.NewLoader().LoadFromFile(inFile)
	if err != nil {
		return err
	}
	specification := generateSpec(doc)
	err = spec.WriteSpecFile(specification, outFile)
	return err
}

func generateSpec(doc *openapi3.T) *spec.Spec {
	api := Api{"default", []PathItem{}}
	for path, pathItem := range doc.Paths {
		for method, o := range pathItem.Operations() {
			pathItem := PathItem{path, method, o}
			api.Items = append(api.Items, pathItem)
		}
	}

	specApi := generateSpecApi(&api)

	apis := spec.Apis{nil, []spec.Api{*specApi}, nil}
	version := spec.Version{spec.Name{"", nil}, spec.VersionSpecification{apis, []spec.NamedModel{}}, []*spec.NamedModel{}}

	meta := spec.Meta{"2.1", spec.Name{doc.Info.Title, nil}, &doc.Info.Title, &doc.Info.Description, doc.Info.Version}
	specification := spec.Spec{meta, []spec.Version{version}}
	return &specification
}

func generateSpecApi(api *Api) *spec.Api {
	operations := []spec.NamedOperation{}
	for _, pathItem := range api.Items {
		operation := generateSpecOperation(&pathItem)
		operations = append(operations, spec.NamedOperation{spec.Name{pathItem.Operation.OperationID, nil}, *operation, nil})
	}
	return &spec.Api{spec.Name{api.Name, nil}, operations, nil}
}

func generateSpecOperation(pathItem *PathItem) *spec.Operation {
	endpoint := spec.Endpoint{pathItem.Method, pathItem.Path, []spec.NamedParam{}, nil}
	var description *string = nil
	if pathItem.Operation.Description != "" {
		description = &pathItem.Operation.Description
	}
	operation := spec.Operation{
		endpoint,
		description,
		[]spec.NamedParam{},
		[]spec.NamedParam{},
		getBody(pathItem.Operation.RequestBody),
		[]spec.NamedResponse{},
		nil,
	}
	return &operation
}

func getBody(body *openapi3.RequestBodyRef) *spec.Definition {
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
	definition := spec.Definition{spec.Type{*SpecTypeFromSchemaRef(media.Schema), nil}, &body.Value.Description, nil}
	return &definition
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
