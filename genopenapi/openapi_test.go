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
		{Name: spec.Name{"field1"}, Field: *spec.NewField(*spec.Plain(spec.TypeString), nil)},
		{Name: spec.Name{"field2"}, Field: *spec.NewField(*spec.Nullable(spec.Plain(spec.TypeString)), &description)},
		{Name: spec.Name{"field3"}, Field: *spec.NewField(*spec.Array(spec.Plain(spec.TypeString)), nil)},
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
	response := spec.NewResponse(*spec.Plain("SomeModel"), nil)
	openapiYaml := generateResponse(*response)
	expected := `
content:
  application/json:
    schema:
      $ref: '#/components/schemas/SomeModel'
`
	assert.Equal(t, strings.TrimSpace(openapiYaml.String()), strings.TrimSpace(expected))
}

func TestApis(t *testing.T) {
	myModel := spec.Plain("MyModel")
	response := spec.NewResponse(*myModel, nil)
	description := "the description"
	operation := spec.Operation{
		Endpoint:     "POST /create/{id:uuid}",
		Description:  &description,
		Body:         spec.NewBody(*myModel, &description),
		HeaderParams: spec.HeaderParams{spec.NamedParam{Name: spec.Name{"Authorization"}, Param: *spec.NewParam(*spec.Plain(spec.TypeString), &description)}},
		QueryParams:  spec.QueryParams{spec.NamedParam{Name: spec.Name{"id"}, Param: *spec.NewParam(*spec.Plain(spec.TypeUuid), &description)}},
		Responses:    spec.Responses{spec.NamedResponse{Name: spec.Name{"ok"}, Response: *response}},
	}
	operation.Init()
	api := spec.Api{
		Name: spec.Name{"mine"},
		Operations: spec.Operations{
			spec.NamedOperation{
				Name:      spec.Name{"create"},
				Operation: operation,
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
