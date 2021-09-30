package gengo

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
)

func GenerateGoClient(serviceFile string, moduleName string, generatePath string) error {
	result, err := spec.ReadSpecFile(serviceFile)
	if err != nil {
		return err
	}

	specification := result.Spec

	generatedFiles := []gen.TextFile{}

	for _, version := range specification.Versions {
		versionModule := Module(moduleName, createPath(generatePath, version.Version.FlatCase()))
		modelsModule := Module(moduleName, createPath(generatePath, version.Version.FlatCase(), modelsPackage))

		generatedFiles = append(generatedFiles, *generateConverter(versionModule))
		generatedFiles = append(generatedFiles, generateVersionModels(&version, modelsModule)...)
		generatedFiles = append(generatedFiles, generateClientsImplementations(&version, versionModule, modelsModule)...)
		generatedFiles = append(generatedFiles, *generateServicesResponses(&version, versionModule, modelsModule))
	}
	err = gen.WriteFiles(generatedFiles, true)
	return err
}

func generateServicesResponses(version *spec.Version, versionModule module, modelsModule module) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", versionModule.Name)

	imports := Imports().Add(modelsModule.Package)
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
		Path:    versionModule.GetPath("responses.go"),
		Content: w.String(),
	}
}
