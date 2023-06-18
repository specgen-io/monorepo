package spec

import "strconv"

type Response struct {
	Name Name
	ResponseBody
}

func (response *Response) BodyKind() ResponseBodyKind {
	return kindOfResponseBody(&response.ResponseBody)
}

func (response *Response) BodyIs(kind ResponseBodyKind) bool {
	return kindOfResponseBody(&response.ResponseBody) == kind
}

func (response *Response) IsSuccess() bool {
	statusCode, _ := strconv.Atoi(HttpStatusCode(response.Name))
	return statusCode >= 200 && statusCode <= 399
}

func (response *Response) IsError() bool {
	statusCode, _ := strconv.Atoi(HttpStatusCode(response.Name))
	return statusCode >= 400 && statusCode <= 599
}
