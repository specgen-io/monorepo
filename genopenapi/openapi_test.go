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
	specYaml := `
idl_version: 2
name: bla-api
models:
  Model:
    oneOf:
      one: Model1
      two: Model2
  Model1:
    field1: string
  Model2:
    field1: string
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
      properties:
        one:
          $ref: '#/components/schemas/Model1'
        two:
          $ref: '#/components/schemas/Model2'
    Model1:
      type: object
      required:
        - field1
      properties:
        field1:
          type: string
    Model2:
      type: object
      required:
        - field1
      properties:
        field1:
          type: string
`

	spec, err := spec.ParseSpec([]byte(specYaml))
	assert.Equal(t, err, nil)

	openapiYaml, err := ToYamlString(generateOpenapi(spec))
	assert.NilError(t, err)

	assert.Equal(t, strings.TrimSpace(expectedOpenApiYaml), strings.TrimSpace(openapiYaml))
}

func TestApis(t *testing.T) {
	specYaml := `
idl_version: 2
name: bla-api
http:
    mine:
        create:
            endpoint: POST /create/{id:uuid}
            description: the description
            body:
                description: the description
                type: MyModel
            header:
                Authorization: string
            query:
                uuid_param:
                    type: uuid
                    description: the description
                str_param:
                    type: string
                    default: the default value
                    description: the description
            response:
                ok: MyModel
models:
    MyModel:
        field1: string
`

	expectedOpenApiYaml := `
openapi: 3.0.0
info:
  version: ""
paths:
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
          name: uuid_param
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
components:
  schemas:
    MyModel:
      type: object
      required:
        - field1
      properties:
        field1:
          type: string
`
	//assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(openapiYaml))
	spec, err := spec.ParseSpec([]byte(specYaml))
	assert.Equal(t, err, nil)

	openapiYaml, err := ToYamlString(generateOpenapi(spec))
	assert.NilError(t, err)

	assert.Equal(t, strings.TrimSpace(expectedOpenApiYaml), strings.TrimSpace(openapiYaml))
}

func TestSpecification(t *testing.T) {
	specYaml := `
idl_version: 2
name: bla-api
title: The Service
description: The service with description
version: 0
`

	expectedOpenApiYaml := `
openapi: 3.0.0
info:
  title: The Service
  description: The service with description
  version: "0"
paths: {}
components:
  schemas: {}
`

	spec, err := spec.ParseSpec([]byte(specYaml))
	assert.Equal(t, err, nil)

	openapiYaml, err := ToYamlString(generateOpenapi(spec))
	assert.NilError(t, err)

	assert.Equal(t, strings.TrimSpace(expectedOpenApiYaml), strings.TrimSpace(openapiYaml))
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

