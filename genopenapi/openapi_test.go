package genopenapi

import (
	spec "github.com/specgen-io/spec.v2"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func TestEnumModel(t *testing.T) {
	specYaml := `
idl_version: 2
name: bla-api
models:
  Model:
    description: The description
    enum:
      first:
        value: FIRST
        description: First option
      second:
        value: SECOND
        description: Second option
      third:
        value: THIRD
        description: Third option
`

	expectedOpenApiYaml := `
openapi: 3.0.0
info:
  version: ""
paths: {}
components:
  schemas:
    Model:
      type: string
      description: The description
      enum:
        - first
        - second
        - third
`

	spec, err := spec.ParseSpec([]byte(specYaml))
	assert.Equal(t, err, nil)

	openapiYaml, err := ToYamlString(generateOpenapi(spec))
	assert.NilError(t, err)

	assert.Equal(t, strings.TrimSpace(expectedOpenApiYaml), strings.TrimSpace(openapiYaml))
}

func TestObjectModel(t *testing.T) {
	specYaml := `
idl_version: 2
name: bla-api
models:
  Model:
    fields:
      field1: 
        type: string
      field2:
        type: string
        description: the description
      field3:
        type: string[]
`

	expectedOpenApiYaml := `
openapi: 3.0.0
info:
  version: ""
paths: {}
components:
  schemas:
    Model:
      type: object
      required:
        - field1
        - field2
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

	spec, err := spec.ParseSpec([]byte(specYaml))
	assert.Equal(t, err, nil)

	openapiYaml, err := ToYamlString(generateOpenapi(spec))
	assert.NilError(t, err)

	assert.Equal(t, strings.TrimSpace(expectedOpenApiYaml), strings.TrimSpace(openapiYaml))
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
	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(openapiYaml))
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
	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(openapiYaml))
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
	apis := []spec.Api{api}
	versionedApi := spec.VersionedApis{Version: NewName(""), Apis: apis}
	versionedApis := []spec.VersionedApis{versionedApi}

	openapiYaml, err := ToYamlString(generateApis(versionedApis))
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
	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(openapiYaml))
}

func TestSpecification(t *testing.T) {
	title := "The Service"
	description := "The service with description"
	spec := spec.Spec{
		IdlVersion:  "0",
		Name:        NewName("the-service"),
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
	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(openapiYaml))
}

func TestFullSpecificationNoVersions(t *testing.T) {
	specYaml := `
idl_version: 2
name: bla-api
title: Bla API
description: Some Bla API service
version: 0

http:
    test:
        some_url:
            endpoint: GET /some/url
            response:
                ok: Model1
        ping:
            endpoint: GET /ping
            query:
                message: string?
            response:
                ok: empty

models:
  Model1:
    prop1: string
  Model2:
    prop1: string
    prop2: int32
`

	expectedOpenApiYaml := `
openapi: 3.0.0
info:
  title: Bla API
  description: Some Bla API service
  version: "0"
paths:
  /some/url:
    get:
      operationId: testSomeUrl
      tags:
        - test
      responses:
        "200":
          description: ""
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Model1'
  /ping:
    get:
      operationId: testPing
      tags:
        - test
      parameters:
        - in: query
          name: message
          required: false
          schema:
            type: string
      responses:
        "200":
          description: ""
components:
  schemas:
    Model1:
      type: object
      required:
        - prop1
      properties:
        prop1:
          type: string
    Model2:
      type: object
      required:
        - prop1
        - prop2
      properties:
        prop1:
          type: string
        prop2:
          type: integer
          format: int32
`

	spec, err := spec.ParseSpec([]byte(specYaml))
	assert.Equal(t, err, nil)

	openapiYaml, err := ToYamlString(generateOpenapi(spec))
	assert.NilError(t, err)

	assert.Equal(t, strings.TrimSpace(expectedOpenApiYaml), strings.TrimSpace(openapiYaml))
}

func TestFullSpecificationWithVersions(t *testing.T) {
	specYaml := `
idl_version: 2
name: bla-api
title: Bla API
description: Some Bla API service
version: 0

http:
    v2:
        test:
            some_url:
                endpoint: GET /some/url
                response:
                    ok: Message

models:
  v2:
    Message:
      prop1: string
`

	expectedOpenApiYaml := `
openapi: 3.0.0
info:
  title: Bla API
  description: Some Bla API service
  version: "0"
paths:
  /v2/some/url:
    get:
      operationId: v2TestSomeUrl
      tags:
        - test
      responses:
        "200":
          description: ""
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/v2.Message'
components:
  schemas:
    v2.Message:
      type: object
      required:
        - prop1
      properties:
        prop1:
          type: string
`

	spec, err := spec.ParseSpec([]byte(specYaml))
	assert.Equal(t, err, nil)

	openapiYaml, err := ToYamlString(generateOpenapi(spec))
	assert.NilError(t, err)

	assert.Equal(t, strings.TrimSpace(expectedOpenApiYaml), strings.TrimSpace(openapiYaml))
}

