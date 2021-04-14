package genopenapi

import (
	spec "github.com/specgen-io/spec.v1"
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
	openapiYaml, err := ToYamlString(generateModel(model))
	assert.NilError(t, err)
	expected := `
type: string
description: The description
enum:
  - first
  - second
  - third
`
	assert.Equal(t, strings.TrimSpace(openapiYaml), strings.TrimSpace(expected))
}

func TestObjectModel(t *testing.T) {
	description := "the description"
	fields := []spec.NamedDefinition{
		*NewField("field1", *spec.Plain(spec.TypeString), nil),
		*NewField("field2", *spec.Nullable(spec.Plain(spec.TypeString)), &description),
		*NewField("field3", *spec.Array(spec.Plain(spec.TypeString)), nil),
	}
	model := spec.Model{Object: NewObject(fields, nil)}
	openapiYaml, err := ToYamlString(generateModel(model))
	assert.NilError(t, err)
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
	assert.Equal(t, strings.TrimSpace(openapiYaml), strings.TrimSpace(expected))
}

func TestUnionModel(t *testing.T) {
	items := []spec.NamedDefinition{
		*NewField("field1", *spec.Plain("Model1"), nil),
		*NewField("field2", *spec.Plain("Model2"), nil),
		*NewField("field3", *spec.Plain("Model3"), nil),
	}
	model := spec.Model{OneOf: NewOneOf(items, nil)}
	openapiYaml, err := ToYamlString(generateModel(model))
	assert.NilError(t, err)
	expected := `
type: object
properties:
  field1:
    $ref: '#/components/schemas/Model1'
  field2:
    $ref: '#/components/schemas/Model2'
  field3:
    $ref: '#/components/schemas/Model3'
`
	assert.Equal(t, strings.TrimSpace(openapiYaml), strings.TrimSpace(expected))
}

func TestResponse(t *testing.T) {
	response := spec.Definition{spec.Type{*spec.Plain("SomeModel"), nil}, nil, nil}
	openapiYaml, err := ToYamlString(generateResponse(response))
	assert.NilError(t, err)
	expected := `
description: ""
content:
  application/json:
    schema:
      $ref: '#/components/schemas/SomeModel'
`
	assert.Equal(t, strings.TrimSpace(openapiYaml), strings.TrimSpace(expected))
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
	openapiYaml, err := ToYamlString(generateApis(apis))
	assert.NilError(t, err)
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
	assert.Equal(t, strings.TrimSpace(openapiYaml), strings.TrimSpace(expected))
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

	openapiYaml, err := ToYamlString(generateOpenapi(&spec))
	assert.NilError(t, err)

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
	assert.Equal(t, strings.TrimSpace(openapiYaml), strings.TrimSpace(expected))
}
