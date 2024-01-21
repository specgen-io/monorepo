package spec

type BodyKind string

const (
	BodyEmpty          BodyKind = "empty"
	BodyText           BodyKind = "string"
	BodyBinary         BodyKind = "binary"
	BodyFile           BodyKind = "file"
	BodyJson           BodyKind = "json"
	BodyFormData       BodyKind = "form-data"
	BodyFormUrlEncoded BodyKind = "form-urlencoded"
)
