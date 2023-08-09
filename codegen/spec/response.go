package spec

import "strconv"

type Response struct {
	Name        Name
	Body        ResponseBody
	Description *string
}

func (response *Response) IsSuccess() bool {
	statusCode, _ := strconv.Atoi(HttpStatusCode(response.Name))
	return statusCode >= 200 && statusCode <= 399
}

func (response *Response) IsError() bool {
	statusCode, _ := strconv.Atoi(HttpStatusCode(response.Name))
	return statusCode >= 400 && statusCode <= 599
}
