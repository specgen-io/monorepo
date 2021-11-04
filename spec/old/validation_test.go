package old

import (
	"gotest.tools/assert"
	"strings"
	"testing"
)

func Test_Body_NonObject_Error(t *testing.T) {
	data := `
http:
    test:
      some_url:
        endpoint: GET /some/url
        body: string
        response:
          ok: empty
`

	old, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(old)
	assert.Equal(t, len(errors), 0)

	warnings, errors := validate(old)
	assert.Equal(t, len(warnings), 0)
	assert.Equal(t, len(errors), 1)
	assert.Equal(t, strings.Contains(errors[0].Message, "body"), true)
}

func Test_Response_NonObject_Error(t *testing.T) {
	data := `
http:
    test:
      some_url:
        endpoint: GET /some/url
        response:
          ok: string
`

	old, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(old)
	assert.Equal(t, len(errors), 0)

	warnings, errors := validate(old)
	assert.Equal(t, len(warnings), 0)
	assert.Equal(t, len(errors), 1)
	assert.Equal(t, strings.Contains(errors[0].Message, "response"), true)
}

func Test_Response_NonEmpty_SpecialStatus_Error(t *testing.T) {
	data := `
http:
    test:
      some_url:
        endpoint: GET /some/url
        response:
          internal_server_error: string{}
          not_found: string{}
          bad_request: string{}
`

	old, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(old)
	assert.Equal(t, len(errors), 0)

	warnings, errors := validate(old)
	assert.Equal(t, len(warnings), 3)
	assert.Equal(t, len(errors), 0)
	assert.Equal(t, strings.Contains(warnings[0].Message, "internal_server_error"), true)
	assert.Equal(t, strings.Contains(warnings[1].Message, "not_found"), true)
	assert.Equal(t, strings.Contains(warnings[2].Message, "bad_request"), true)
}

func Test_Query_Param_Array_Allowed(t *testing.T) {
	data := `
http:
    test:
      some_url:
        endpoint: GET /some/url
        query:
          param1: int[]
        response:
          ok: empty
`

	old, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(old)
	assert.Equal(t, len(errors), 0)

	warnings, errors := validate(old)
	assert.Equal(t, len(warnings), 0)
	assert.Equal(t, len(errors), 0)
}

func Test_Query_Param_Object_Errors(t *testing.T) {
	data := `
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
   field: string
`

	old, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(old)
	assert.Equal(t, len(errors), 0)

	warnings, errors := validate(old)
	assert.Equal(t, len(warnings), 0)
	assert.Equal(t, len(errors), 2)
	assert.Equal(t, strings.Contains(errors[0].Message, "param1"), true)
	assert.Equal(t, strings.Contains(errors[1].Message, "param2"), true)
}

func Test_Params_SameName_Error(t *testing.T) {
	data := `
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
`

	old, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(old)
	assert.Equal(t, len(errors), 0)

	warnings, errors := validate(old)
	assert.Equal(t, len(warnings), 0)
	assert.Equal(t, len(errors), 1)
	assert.Equal(t, strings.Contains(errors[0].Message, "the_param"), true)
}

func Test_NonDefaultable_Type_Error(t *testing.T) {
	data := `
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
`

	old, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(old)
	assert.Equal(t, len(errors), 0)

	warnings, errors := validate(old)
	assert.Equal(t, len(warnings), 0)
	assert.Equal(t, len(errors), 2)
	assert.Equal(t, strings.Contains(errors[0].Message, "string?"), true)
	assert.Equal(t, strings.Contains(errors[1].Message, "date?"), true)
}

func Test_Defaulted_Format_Pass(t *testing.T) {
	data := `
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
`

	old, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(old)
	assert.Equal(t, len(errors), 0)

	warnings, errors := validate(old)
	assert.Equal(t, len(warnings), 0)
	assert.Equal(t, len(errors), 0)
}

func Test_Defaulted_Format_Fail(t *testing.T) {
	data := `
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
`

	old, err := unmarshalSpec([]byte(data))
	assert.Equal(t, err, nil)

	errors := enrichSpec(old)
	assert.Equal(t, len(errors), 0)

	warnings, errors := validate(old)
	assert.Equal(t, len(warnings), 0)
	assert.Equal(t, len(errors), 10)
}
