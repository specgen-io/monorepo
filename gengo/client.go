package gengo

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"path/filepath"
)

func GenerateGoClient(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}
	generatedFiles := []gen.TextFile{}

	packageName := "spec"
	for _, version := range specification.Versions {
		versionPackageName := versionedPackage(version.Version, packageName)
		generatedPath := filepath.Join(generatePath, packageName)
		versionPath := createPath(version.Version.FlatCase(), generatedPath)

		generatedFiles = append(generatedFiles, *generateConverter(versionPackageName, filepath.Join(versionPath, "converter.go")))
		generatedFiles = append(generatedFiles, generateVersionModels(&version, createPath(versionPath, modelsPackage))...)
		generatedFiles = append(generatedFiles, generateClientsImplementations(&version, versionPackageName, versionPath)...)
		generatedFiles = append(generatedFiles, *generateServicesResponses(&version, versionPackageName, versionPath, "responses.go"))
	}
	err = gen.WriteFiles(generatedFiles, true)
	return err
}

func generateServicesResponses(version *spec.Version, packageName string, generatePath string, fileName string) *gen.TextFile {
	w := NewGoWriter()
	w.Line("package %s", packageName)

	imports := []string{}
	imports = generateVersionImports(version, imports)
	if len(imports) > 0 {
		w.EmptyLine()
		for _, imp := range imports {
			w.Line(`import %s`, imp)
		}
	}

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
