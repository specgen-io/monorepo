package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/specgen-io/specgen/spec/v2"
)

func (c *Converter) Specification(doc *openapi3.T) *spec.Spec {
	meta := spec.Meta{spec.SpecVersion, name(doc.Info.Title), &doc.Info.Title, &doc.Info.Description, doc.Info.Version}
	version := spec.Version{
		name(""),
		spec.VersionSpecification{c.apis(doc), c.models(doc.Components.Schemas)},
		nil,
	}
	specification := spec.Spec{meta, []spec.Version{version}}
	return &specification
}
