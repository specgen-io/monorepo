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

	packageName := "spec"
	for _, version := range specification.Versions {
		versionPackageName := versionedPackage(version.Version, packageName)
		versionPath := versionedFolder(version.Version, filepath.Join(generatePath, packageName))

		generatedFiles = append(generatedFiles, *generateConverter(versionPackageName, filepath.Join(versionPath, "converter.go")))
		generatedFiles = append(generatedFiles, generateVersionModels(&version, createPath(versionPath, modelsPackage))...)
		generatedFiles = append(generatedFiles, generateClientsImplementations(&version, moduleName, versionPackageName, versionPath)...)
		generatedFiles = append(generatedFiles, *generateServicesResponses(&version, moduleName, versionPackageName, versionPath, "responses.go"))
	}
	err = gen.WriteFiles(generatedFiles, true)
	return err
}

func generateServicesResponses(version *spec.Version, rootPackage string, packageName string, generatePath string, fileName string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)

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
		Path:    filepath.Join(generatePath, fileName),
		Content: w.String(),
	}
}
