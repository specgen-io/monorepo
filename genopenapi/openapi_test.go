package genopenapi

import (
	"github.com/ModaOperandi/spec"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func TestEnumModel(t *testing.T) {
	description := "The description"
	model := spec.Model{Enum: &spec.Enum{
		Description: &description,
		Items: []spec.NamedEnumItem{
			{Name: spec.Name{"first"}, EnumItem: spec.EnumItem{}},
			{Name: spec.Name{"second"}, EnumItem: spec.EnumItem{}},
			{Name: spec.Name{"third"}, EnumItem: spec.EnumItem{}},
		},
	}}
	openapiYaml := generateModel(model)
	expected := `
type: string
description: The description
enum:
- first
- second
- third
`
	assert.Equal(t, strings.TrimSpace(openapiYaml.String()), strings.TrimSpace(expected))
}

func TestObjectModel(t *testing.T) {
	description := "the description"
	fields := []spec.NamedField{
		*spec.NewField("field1", *spec.Plain(spec.TypeString), nil, nil),
		*spec.NewField("field2", *spec.Nullable(spec.Plain(spec.TypeString)), nil, &description),
		*spec.NewField("field3", *spec.Array(spec.Plain(spec.TypeString)), nil, nil),
	}
	model := spec.Model{Object: spec.NewObject(fields, nil)}
	openapiYaml := generateModel(model)
	expected := `
type: object
required:
- field1
- field3
properties:
  field1:
    type: string
  field2:
    type: string
    description: the description
  field3:
    type: array
    items:
      type: string
`
	assert.Equal(t, strings.TrimSpace(openapiYaml.String()), strings.TrimSpace(expected))
}

func TestResponse(t *testing.T) {
	response := spec.NewDefinition(*spec.Plain("SomeModel"), nil)
	openapiYaml := generateResponse(*response)
	expected := `
description: ""
content:
  application/json:
    schema:
      $ref: '#/components/schemas/SomeModel'
`
	assert.Equal(t, strings.TrimSpace(openapiYaml.String()), strings.TrimSpace(expected))
}

func TestApis(t *testing.T) {
	myModel := spec.Plain("MyModel")
	description := "the description"
	operation := spec.NewOperation(
		spec.NewEndpoint("POST /create/{id:uuid}"),
		&description,
		spec.NewDefinition(*myModel, &description),
		spec.HeaderParams{*spec.NewParam("Authorization", *spec.Plain(spec.TypeString), nil, &description)},
		spec.QueryParams{*spec.NewParam("id", *spec.Plain(spec.TypeUuid), nil, &description)},
		spec.Responses{spec.NamedResponse{Name: spec.Name{"ok"}, Definition: *spec.NewDefinition(*myModel, nil)}},
	)
	api := spec.Api{
		Name: spec.Name{"mine"},
		Operations: spec.Operations{
			spec.NamedOperation{
				Name:      spec.Name{"create"},
				Operation: *operation,
			},
		},
	}
	apis := spec.Apis{api}
	openapiYaml := generateApis(apis)
	expected := `
/create/{id}:
  post:
    operationId: mineCreate
    tags:
    - mine
    description: the description
    requestBody:
      description: the description
      required: true
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/MyModel'
    parameters:
    - in: path
      name: id
      required: true
      schema:
        type: string
        format: uuid
    - in: header
      name: Authorization
      required: true
      schema:
        type: string
      description: the description
    - in: query
      name: id
      required: true
      schema:
        type: string
        format: uuid
      description: the description
    responses:
      "200":
        description: ""
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/MyModel'
`
	assert.Equal(t, strings.TrimSpace(openapiYaml.String()), strings.TrimSpace(expected))
}

func TestSpecification(t *testing.T) {
	idlVersion := "0"
	title := "The Service"
	description := "The service with description"
	spec := spec.Spec{
		IdlVersion:  &idlVersion,
		ServiceName: spec.Name{"the-service"},
		Title:       &title,
		Description: &description,
		Version:     "0",
	}

	openapi := generateOpenapi(&spec)

	expected := `
openapi: 3.0.0
info:
  title: The Service
  description: The service with description
  version: "0"
paths: {}
components:
  schemas: {}
`
	assert.Equal(t, strings.TrimSpace(openapi.String()), strings.TrimSpace(expected))
}
