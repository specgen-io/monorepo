package spec

import (
	"errors"
	"gotest.tools/assert"
	"testing"
)

func Test_Enrichment(t *testing.T) {
	runReadSpecificationCases(t, enrichmentCases)
}

var enrichmentCases = []ReadSpecificationCase{
	{
		`resolve operations builtin type no errors`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url/{id:string}
      query:
        the_query: string
      header:
        The-Header: string
      response:
        ok: empty
`,
		nil,
		[]Message{},
		nil,
	},
	{
		`resolve operations unknown type error`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url/{id:nonexisting1}
      query:
        the_query: nonexisting2
      header:
        The-Header: nonexisting3
      response:
        ok: empty
`,
		errors.New(`failed to parse specification`),
		[]Message{
			Error(`unknown type: nonexisting1`).At(&Location{specificationMetaLines + 4, 17}),
			Error(`unknown type: nonexisting2`).At(&Location{specificationMetaLines + 6, 20}),
			Error(`unknown type: nonexisting3`).At(&Location{specificationMetaLines + 8, 21}),
		},
		nil,
	},
	{
		`resolve operations custom type no errors`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      body: Custom1
      response:
        ok: Custom2
models:
  Custom1:
    object:
      field: string
  Custom2:
    object:
      field: string
`,
		nil,
		[]Message{},
		nil,
	},
	{
		`resolve models object field custom type no errors`,
		`
models:
  Custom1:
    object:
      field1: string
      field2: Custom2
  Custom2:
    object:
      field: Custom3
  Custom3:
    enum:
      - first
      - second
`,
		nil,
		[]Message{},
		nil,
	},
	{
		`resolve models object field unknown type error`,
		`
models:
  Custom:
    object:
      field1: NonExisting
`,
		errors.New(`failed to parse specification`),
		[]Message{Error(`unknown type: NonExisting`).At(&Location{specificationMetaLines + 4, 15})},
		nil,
	},
	{
		`resolve models union item no errors`,
		`
models:
  Custom1:
    object:
      field1: string
      field2: Custom2
  Custom2:
    oneOf:
      one: string
      two: boolean[]
      three: int?
`,
		nil,
		[]Message{},
		nil,
	},
	{
		`resolve models union item unknown type error`,
		`
models:
  Custom:
    oneOf:
      nope: NonExisting
`,
		errors.New(`failed to parse specification`),
		[]Message{Error(`unknown type: NonExisting`).At(&Location{specificationMetaLines + 4, 13})},
		nil,
	},
	{
		`resolve models normal order no errors`,
		`
models:
  Model1:
    object:
      field: string
  Model2:
    object:
      field: string
`,
		nil,
		[]Message{},
		func(t *testing.T, spec *Spec) {
			assert.Equal(t, len(spec.Versions), 1)
			models := spec.Versions[0].ResolvedModels
			assert.Equal(t, len(models), 2)
			assert.Equal(t, models[0].Name.Source, "Model1")
			assert.Equal(t, models[1].Name.Source, "Model2")
		},
	},
	{
		`resolve models reversed order no errors`,
		`
models:
  Model1:
    object:
      field: Model2
  Model2:
    object:
      field: Model3
  Model3:
    object:
      field: string
`,
		nil,
		[]Message{},
		func(t *testing.T, spec *Spec) {
			assert.Equal(t, len(spec.Versions), 1)
			version := &spec.Versions[0]
			models := version.ResolvedModels
			assert.Equal(t, len(models), 3)
			assert.Equal(t, models[0].Name.Source, "Model3")
			assert.Equal(t, models[0].InVersion, version)
			assert.Equal(t, models[1].Name.Source, "Model2")
			assert.Equal(t, models[1].InVersion, version)
			assert.Equal(t, models[2].Name.Source, "Model1")
			assert.Equal(t, models[2].InVersion, version)
		},
	},
	{
		`resolve operations no errors`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      response:
        ok: empty
`,
		nil,
		[]Message{},
		func(t *testing.T, spec *Spec) {
			version := &spec.Versions[0]
			apis := &version.Http
			api := &apis.Apis[0]
			operation := &api.Operations[0]
			response := operation.Responses[0]
			assert.Equal(t, apis.InVersion, version)
			assert.Equal(t, api.InHttp, apis)
			assert.Equal(t, operation.InApi, api)
			assert.Equal(t, response.Operation, operation)
		},
	},
}
