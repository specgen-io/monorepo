package genopenapi

import (
	"github.com/ModaOperandi/spec"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func TestEnumModel(t *testing.T) {
	model := spec.Model{Enum: &spec.Enum{Items: []spec.Name{spec.Name{"first"}, spec.Name{"second"}, spec.Name{"third"}}}}
	openapiYaml := generateModel(model)
	expected := `
type: string
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
		HeaderParams: spec.HeaderParams{spec.NamedParam{Name: spec.Name{"Authorization"}, Param: *spec.NewParam(*spec.Plain(spec.TypeString), nil)}},
		QueryParams:  spec.QueryParams{spec.NamedParam{Name: spec.Name{"id"}, Param: *spec.NewParam(*spec.Plain(spec.TypeUuid), nil)}},
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
    - in: query
      name: id
      required: true
      schema:
        type: string
        format: uuid
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
