package gents

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
	"path/filepath"
)

func GenerateExpressService(serviceFile string, swaggerPath string, generatePath string, servicesPath string, validation string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	sources := []gen.TextFile{}

	superstruct := generateSuperstruct(filepath.Join(generatePath, "superstruct.ts"))
	sources = append(sources, *superstruct)

	for _, version := range specification.Versions {
		sources = append(sources, *genereateServiceApis(&version, generatePath))
		sources = append(sources, *generateExpressVersionRouting(&version, generatePath))
		sources = append(sources, *generateModels(&version, generatePath))
	}
	sources = append(sources, *genopenapi.GenerateOpenapi(specification, swaggerPath))

	err = gen.WriteFiles(sources, true)
	if err != nil {
		return err
	}

	return nil
}

func generateModels(version *spec.Version, generatePath string) *gen.TextFile {
	w := NewTsWriter()
	generateSuperstructModels(w, version)
	return &gen.TextFile{Path: filepath.Join(generatePath, versionFilename(version, "models", "ts")), Content: w.String()}
}

func versionFilename(version *spec.Version, filename string, ext string) string {
	if version.Version.Source != "" {
		filename = filename + "-" + version.Version.FlatCase()
	}
	if ext != "" {
		filename = filename + "." + ext
	}
	return filename
}