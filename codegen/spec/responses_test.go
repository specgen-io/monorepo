package spec

import (
	"errors"
	"gotest.tools/assert"
	"testing"
)

func Test_Respnses(t *testing.T) {
	runReadSpecificationCases(t, responsesCases)
}

var responsesCases = []ReadSpecificationCase{
	{
		`non required error response`,
		`
errors:
  responses:
    forbidden: ForbiddenError
  models:
    ForbiddenError:
      object:
        message: string
http:
  test:
    some_url:
      endpoint: GET /some/url
      body: string
      response:
        ok: empty
        bad_request: BadRequestError
        forbidden: ForbiddenError
`,
		nil,
		[]Message{},
		func(t *testing.T, spec *Spec) {
			bad_request := spec.HttpErrors.Responses.GetByStatusName(HttpStatusBadRequest)
			assert.Equal(t, bad_request != nil, true)
			assert.Equal(t, bad_request.Required, true)
			internal_server_error := spec.HttpErrors.Responses.GetByStatusName(HttpStatusInternalServerError)
			assert.Equal(t, internal_server_error != nil, true)
			assert.Equal(t, internal_server_error.Required, true)
			not_found := spec.HttpErrors.Responses.GetByStatusName(HttpStatusNotFound)
			assert.Equal(t, not_found != nil, true)
			assert.Equal(t, not_found.Required, true)
			forbidden := spec.HttpErrors.Responses.GetByStatusName("forbidden")
			assert.Equal(t, forbidden != nil, true)
			assert.Equal(t, forbidden.Required, false)
		},
	},
	{
		`non required error response not declared`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      body: string
      response:
        ok: empty
        forbidden: ForbiddenError
`,
		errors.New("failed to parse specification"),
		[]Message{
			Error(`unknown type: ForbiddenError`).At(&Location{specificationMetaLines + 8, 20}),
		},
		nil,
	},
	{
		`non required error response different body`,
		`
errors:
  responses:
    forbidden: ForbiddenError
  models:
    ForbiddenError:
      object:
        message: string
http:
  test:
    some_url:
      endpoint: GET /some/url
      body: string
      response:
        ok: empty
        forbidden: empty
`,
		errors.New("failed to validate specification"),
		[]Message{
			Error(`response forbidden is declared with body type: empty, however errors section declares it with body type: ForbiddenError`).At(&Location{specificationMetaLines + 15, 20}),
		},
		nil,
	},
}
