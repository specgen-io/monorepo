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
			*NewEnumItem("first", nil),
			*NewEnumItem("second", nil),
			*NewEnumItem("third", nil),
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
	defaultValue := "the default value"
	description := "the description"
	fields := []spec.NamedField{
		*NewField("field1", *spec.Plain(spec.TypeString), nil, nil),
		*NewField("field2", *spec.Nullable(spec.Plain(spec.TypeString)), &defaultValue, &description),
		*NewField("field3", *spec.Array(spec.Plain(spec.TypeString)), nil, nil),
	}
	model := spec.Model{Object: NewObject(fields, nil)}
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
    default: the default value
    description: the description
  field3:
    type: array
    items:
      type: string
`
	assert.Equal(t, strings.TrimSpace(openapiYaml.String()), strings.TrimSpace(expected))
}

func TestResponse(t *testing.T) {
	response := spec.Definition{spec.Type{*spec.Plain("SomeModel"), nil}, nil, nil}
	openapiYaml := generateResponse(response)
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
	defaultValue := "the default value"
	operation := spec.Operation{
		*spec.ParseEndpoint("POST /create/{id:uuid}"),
		&description,
		&spec.Definition{spec.Type{*myModel, nil}, &description, nil},
		spec.HeaderParams{*NewParam("Authorization", *spec.Plain(spec.TypeString), nil, &description)},
		spec.QueryParams{
			*NewParam("id", *spec.Plain(spec.TypeUuid), nil, &description),
			*NewParam("str_param", *spec.Plain(spec.TypeString), &defaultValue, &description),
		},
		spec.Responses{*NewResponse("ok", *myModel, nil)},
	}
	api := spec.Api{
		Name: NewName("mine"),
		Operations: spec.Operations{
			spec.NamedOperation{
				Name:      NewName("create"),
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
    - in: query
      name: str_param
      required: true
      schema:
        type: string
        default: the default value
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
		ServiceName: NewName("the-service"),
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
