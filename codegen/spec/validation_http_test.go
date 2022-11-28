package spec

import (
	"errors"
	"testing"
)

func Test_Validation_Http(t *testing.T) {
	runReadSpecificationCases(t, validationHttpCases)
}

var validationHttpCases = []ReadSpecificationCase{
	{
		`string request body no errors`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      body: string
      response:
        ok: empty
`,
		nil,
		[]Message{},
		nil,
	},
	{
		`string response body no errors`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      response:
        ok: string
`,
		nil,
		[]Message{},
		nil,
	},
	{
		`scalar request body error`,
		`
http:
    test:
      some_url:
        endpoint: GET /some/url
        body: int
        response:
          ok: empty
`,
		errors.New("failed to validate specification"),
		[]Message{Error("body should be object, array or string type, found int").At(&Location{specificationMetaLines + 5, 15})},
		nil,
	},
	{
		`scalar request body error`,
		`
http:
    test:
      some_url:
        endpoint: GET /some/url
        body: int
        response:
          ok: empty
`,
		errors.New("failed to validate specification"),
		[]Message{Error("body should be object, array or string type, found int").At(&Location{specificationMetaLines + 5, 15})},
		nil,
	},
	{
		`scalar response body error`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      response:
        ok: int
`,
		errors.New("failed to validate specification"),
		[]Message{Error("response ok should be either empty or some type with structure of an object or array, found int").At(&Location{specificationMetaLines + 6, 13})},
		nil,
	},
	{
		`special responses body non empty warnings`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      response:
        internal_server_error: string{}
        not_found: string{}
        bad_request: string{}
`,
		errors.New(`failed to validate specification`),
		[]Message{
			Error(`response internal_server_error is declared with body type: string{}, however errors section declares it with body type: InternalServerError`).At(&Location{specificationMetaLines + 6, 32}),
			Error(`response not_found is declared with body type: string{}, however errors section declares it with body type: NotFoundError`).At(&Location{specificationMetaLines + 7, 20}),
			Error(`response bad_request is declared with body type: string{}, however errors section declares it with body type: BadRequestError`).At(&Location{specificationMetaLines + 8, 22}),
		},
		nil,
	},
	{
		`array query param no errors`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      query:
        param1: int[]
      response:
        ok: empty
`,
		nil,
		[]Message{},
		nil,
	},
	{
		`object and nullable query param errors`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      query:
        param1: int[]?
        param2: TheModel
      response:
        ok: empty
models:
  TheModel:
    object:
      field: string
`,
		errors.New("failed to validate specification"),
		[]Message{
			Error("parameter param1 should be of scalar type or array of scalar type, found int[]?").At(&Location{specificationMetaLines + 6, 17}),
			Error("parameter param2 should be of scalar type or array of scalar type, found TheModel").At(&Location{specificationMetaLines + 7, 17}),
		},
		nil,
	},
	{
		`same name params error`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      query:
        the_param: string
      header:
        The-Param: string
      response:
        ok: empty
`,
		errors.New("failed to validate specification"),
		[]Message{
			Error("parameter name 'The-Param' conflicts with the other parameter name 'the_param'").At(&Location{specificationMetaLines + 8, 9}),
		},
		nil,
	},
	{
		`nullable pararms defaulted errors`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      query:
        the_query_param: string? = the default
      header:
        The-Header-Param: date? = the default
      response:
        ok: empty
`,
		errors.New("failed to validate specification"),
		[]Message{
			Error("type string? can not have default value").At(&Location{specificationMetaLines + 6, 26}),
			Error("type date? can not have default value").At(&Location{specificationMetaLines + 8, 27}),
		},
		nil,
	},
	{
		`pararms defaulted no errors`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      query:
        int: int = 123
        long: long = 123
        float: float = 123.4
        double: double = 123.4
        decimal: decimal = 123.4
        boolean: boolean = true
        string: string = the default value
        uuid: uuid = 58d5e212-165b-4ca0-909b-c86b9cee0111
        date: date = 2019-08-07
        datetime: datetime = 2019-08-07T10:20:30
        the_enum: Enum = second
      response:
        ok: empty
models:
  Enum:
    enum:
      - first
      - second
`,
		nil,
		[]Message{},
		nil,
	},
	{
		`pararms defaulted bad format errors`,
		`
http:
  test:
    some_url:
      endpoint: GET /some/url
      query:
        int: int = -
        long: long = 1.2
        float: float = abc
        double: double = .4
        decimal: decimal = -.
        boolean: boolean = yes
        string: string = the default value
        uuid: uuid = 58d5e212165b4ca0909bc86b9cee0111
        date: date = 2019/08/07
        datetime: datetime = 2019/08/07 10:20:30
        the_enum: Enum = nonexisting
      response:
        ok: empty
models:
  Enum:
    enum:
      - first
      - second
`,
		errors.New("failed to validate specification"),
		[]Message{
			Error("default value format error: '-' is in wrong format, should be integer; examples: 123").At(&Location{specificationMetaLines + 6, 14}),
			Error("default value format error: '1.2' is in wrong format, should be integer; examples: 123").At(&Location{specificationMetaLines + 7, 15}),
			Error("default value format error: 'abc' is in wrong format, should be float; examples: 123.4").At(&Location{specificationMetaLines + 8, 16}),
			Error("default value format error: '.4' is in wrong format, should be float; examples: 123.4").At(&Location{specificationMetaLines + 9, 17}),
			Error("default value format error: '-.' is in wrong format, should be float; examples: 123.4").At(&Location{specificationMetaLines + 10, 18}),
			Error("default value format error: 'yes' is in wrong format, should be boolean; examples: true or false").At(&Location{specificationMetaLines + 11, 18}),
			Error("default value format error: '58d5e212165b4ca0909bc86b9cee0111' is in wrong format, should be uuid; examples: fbd3036f-0f1c-4e98-b71c-d4cd61213f90").At(&Location{specificationMetaLines + 13, 15}),
			Error("default value format error: '2019/08/07' is in wrong format, should be date; examples: 2019-12-31").At(&Location{specificationMetaLines + 14, 15}),
			Error("default value format error: '2019/08/07 10:20:30' is in wrong format, should be datetime; examples: 2019-12-31T15:53:45").At(&Location{specificationMetaLines + 15, 19}),
			Error("default value nonexisting is not defined in the enum Enum").At(&Location{specificationMetaLines + 16, 19}),
		},
		nil,
	},
	{
		`duplicated url error`,
		`
http:
  test:
    first:
      endpoint: GET /some/url
      response:
        ok: empty
    second:
      endpoint: GET /some/url
      response:
        ok: empty
`,
		nil,
		[]Message{Warning(`endpoint "GET /some/url" is used for 2 operations: test.first, test.second`).At(&Location{specificationMetaLines + 4, 7})},
		nil,
	},
}
