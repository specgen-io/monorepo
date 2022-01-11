package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/specgen-io/specgen/v2/spec"
)

func ConvertFromOpenapi(inFile, outFile string) error {
	doc, err := openapi3.NewLoader().LoadFromFile(inFile)
	if err != nil {
		return err
	}
	specification := convertSpec(doc)
	err = spec.WriteSpecFile(specification, outFile)
	return err
}

func convertSpec(doc *openapi3.T) *spec.Spec {
	meta := spec.Meta{spec.SpecVersion, name(doc.Info.Title), &doc.Info.Title, &doc.Info.Description, doc.Info.Version}
	versions := convertVersions(doc)
	specification := spec.Spec{meta, versions}
	return &specification
}

func convertVersions(doc *openapi3.T) []spec.Version {
	version := spec.Version{
		name(""),
		spec.VersionSpecification{convertApis(doc), convertModels(doc.Components.Schemas)},
		nil,
	}
	return []spec.Version{version}
}
