package gengo

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func GenerateGoClient(moduleName string, serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	generatedFiles := []gen.TextFile{}

	for _, version := range specification.Versions {
		versionPath := createPath(generatePath, version.Version.FlatCase())

		generatedFiles = append(generatedFiles, *generateConverter(versionPath))
		generatedFiles = append(generatedFiles, generateVersionModels(&version, createPath(versionPath, modelsPackage))...)
		generatedFiles = append(generatedFiles, generateClientsImplementations(&version, moduleName, versionPath)...)
		generatedFiles = append(generatedFiles, *generateServicesResponses(&version, moduleName, versionPath))
	}
	err = gen.WriteFiles(generatedFiles, true)
	return err
}

func generateServicesResponses(version *spec.Version, rootPackage string, generatePath string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", getShortPackageName(generatePath))

	imports := Imports().Add(createPackageName(rootPackage, generatePath, modelsPackage))
	imports.Write(w)

	w.EmptyLine()
	w.Line(`type EmptyDef struct{}`)
	w.EmptyLine()
	w.Line(`var Empty = EmptyDef{}`)

	for _, api := range version.Http.Apis {
		for _, operation := range api.Operations {
			w.EmptyLine()
			generateOperationResponseStruct(w, operation)
		}
	}

	return &gen.TextFile{
		Path:    filepath.Join(generatePath, "responses.go"),
		Content: w.String(),
	}
}
