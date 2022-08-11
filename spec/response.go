package spec

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
