package spec

import "strconv"

type Response struct {
	Name Name
	Definition
}

func (response *Response) BodyKind() BodyKind {
	return kindOf(&response.Definition)
}

func (response *Response) BodyIs(kind BodyKind) bool {
	return kindOf(&response.Definition) == kind
}

func (response *Response) IsSuccess() bool {
	statusCode, _ := strconv.Atoi(HttpStatusCode(response.Name))
	return statusCode >= 200 && statusCode <= 399
}

func (response *Response) IsError() bool {
	statusCode, _ := strconv.Atoi(HttpStatusCode(response.Name))
	return statusCode >= 400 && statusCode <= 599
}
