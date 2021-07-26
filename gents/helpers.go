package gents

import "github.com/specgen-io/spec"

func versionFilename(version *spec.Version, filename string, ext string) string {
	if version.Version.Source != "" {
		filename = filename + "_" + version.Version.FlatCase()
	}
	if ext != "" {
		filename = filename + "." + ext
	}
	return filename
}

func versionModule(version *spec.Version, filename string) string {
	if version.Version.Source != "" {
		filename = filename + "_" + version.Version.FlatCase()
	}
	return filename
}