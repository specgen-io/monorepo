package spec

type BodyKind string

const (
	BodyEmpty          BodyKind = "empty"
	BodyText           BodyKind = "string"
	BodyJson           BodyKind = "json"
	BodyFormData       BodyKind = "form-data"
	BodyFormUrlEncoded BodyKind = "form-urlencoded"
)
