package models

import (
	"java/packages"
	"spec"
)

type Packages struct {
	models       map[string]packages.Package
	Json         packages.Package
	Errors       packages.Package
	ErrorsModels packages.Package
}

func NewPackages(packageName, generatePath string, specification *spec.Spec) *Packages {
	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}
	generated := packages.New(generatePath, packageName)
	json := generated.Subpackage("json")
	errors := generated.Subpackage("errors")
	errorsModels := errors.Subpackage("models")

	models := map[string]packages.Package{}
	for _, version := range specification.Versions {
		models[version.Name.Source] = generated.Subpackage(version.Name.FlatCase()).Subpackage("models")
	}

	return &Packages{
		models,
		json,
		errors,
		errorsModels,
	}
}

func (p *Packages) Models(version *spec.Version) packages.Package {
	return p.models[version.Name.Source]
}
