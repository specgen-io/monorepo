package golang

import (
	"path/filepath"
	"strings"
)

type module struct {
	RootModule string
	Path       string
	Package    string
	Name       string
}

func Module(rootModule string, path string) module {
	packageName := createPackageName(rootModule, strings.TrimPrefix(path, "./"))
	parts := strings.Split(packageName, "/")
	name := parts[len(parts)-1]
	return module{RootModule: rootModule, Path: path, Package: packageName, Name: name}
}

func (m module) GetPath(filename string) string {
	path := filepath.Join(m.Path, filename)
	return path
}

func (m module) Submodule(name string) module {
	if name != "" {
		return Module(m.RootModule, filepath.Join(m.Path, name))
	}
	return m
}
