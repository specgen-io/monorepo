package openapi

import (
	"github.com/specgen-io/specgen/v2/spec"
	"github.com/specgen-io/specgen/yamlx/v2"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func checkOpenApi(t *testing.T, specYaml, expectedOpenApiYaml string) {
	specOptions := spec.SpecOptionsDefault
	specOptions.AddErrors = false
	spec, _, err := spec.ReadSpecWithOptions(specOptions, []byte(specYaml))
	assert.Equal(t, err, nil)

	openapiYaml, err := yamlx.ToYamlString(generateSpecification(spec))
	assert.NilError(t, err)

	assert.Equal(t, strings.TrimSpace(expectedOpenApiYaml), strings.TrimSpace(openapiYaml))
}

func TestEnumModel(t *testing.T) {
	specYaml := `
spec: 2.1
name: bla-api
models:
  Model:
    description: The description
    enum:
      first: FIRST  # First option
      second: SECOND  # Second option
      third: THIRD  # Third option
`

	expectedOpenApiYaml := `
openapi: 3.0.0
info:
  title: bla-api
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

	checkOpenApi(t, specYaml, expectedOpenApiYaml)
}

func TestObjectModel(t *testing.T) {
	specYaml := `
spec: 2.1
name: bla-api
models:
  Model:
    object:
      field1: string
      field2: string  # the description
      field3: string[]
`

	expectedOpenApiYaml := `
openapi: 3.0.0
info:
  title: bla-api
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

	checkOpenApi(t, specYaml, expectedOpenApiYaml)
}

func TestOneOfWrapperModel(t *testing.T) {
	specYaml := `
spec: 2.1
name: bla-api
models:
  Model:
    oneOf:
      one: Model1
      two: Model2
  Model1:
    object:
      field1: string
  Model2:
    object:
      field1: string
`

	expectedOpenApiYaml := `
openapi: 3.0.0
info:
  title: bla-api
  version: ""
paths: {}
components:
  schemas:
    Model:
      oneOf:
        - type: object
          required:
            - one
          properties:
            one:
              $ref: '#/components/schemas/Model1'
        - type: object
          required:
            - two
          properties:
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

	checkOpenApi(t, specYaml, expectedOpenApiYaml)
}

func TestOneOfDiscriminatorModel(t *testing.T) {
	specYaml := `
spec: 2.1
name: bla-api
models:
  Model:
    discriminator: kind
    oneOf:
      one: Model1
      two: Model2
  Model1:
    object:
      field1: string
  Model2:
    object:
      field1: string
`

	expectedOpenApiYaml := `
openapi: 3.0.0
info:
  title: bla-api
  version: ""
paths: {}
components:
  schemas:
    Model:
      anyOf:
        - $ref: '#/components/schemas/ModelOne'
        - $ref: '#/components/schemas/ModelTwo'
      discriminator:
        propertyName: kind
        mapping:
          one: '#/components/schemas/ModelOne'
          two: '#/components/schemas/ModelTwo'
    ModelOne:
      allOf:
        - $ref: '#/components/schemas/Model1'
        - type: object
          required:
            - kind
          properties:
            kind:
              type: string
    ModelTwo:
      allOf:
        - $ref: '#/components/schemas/Model2'
        - type: object
          required:
            - kind
          properties:
            kind:
              type: string
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

	checkOpenApi(t, specYaml, expectedOpenApiYaml)
}

func TestApis(t *testing.T) {
	specYaml := `
spec: 2.1
name: bla-api
http:
    mine:
        create:
            endpoint: POST /create/{id:uuid}
            description: the description
            body: MyModel  # the description
            header:
                Authorization: string
            query:
                uuid_param: uuid  # the description
                str_param: string = the default value  # the description
            response:
                ok: MyModel
models:
    MyModel:
        object:
            field1: string
`

	expectedOpenApiYaml := `
openapi: 3.0.0
info:
  title: bla-api
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

	checkOpenApi(t, specYaml, expectedOpenApiYaml)
}

func TestSpecification(t *testing.T) {
	specYaml := `
spec: 2.1
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

	checkOpenApi(t, specYaml, expectedOpenApiYaml)
}

func TestFullSpecificationNoVersions(t *testing.T) {
	specYaml := `
spec: 2.1
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
    object:
      prop1: string
  Model2:
    object:
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

	checkOpenApi(t, specYaml, expectedOpenApiYaml)
}

func TestFullSpecificationWithVersions(t *testing.T) {
	specYaml := `
spec: 2.1
name: bla-api
title: Bla API
description: Some Bla API service
version: 0

v2:
  http:
    test:
      some_url:
        endpoint: GET /some/url
        response:
          ok: Message

  models:
    Message:
      object:
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

	checkOpenApi(t, specYaml, expectedOpenApiYaml)
}
