package module

import (
	"path/filepath"
	"strings"
)

type Module struct {
	RootModule string
	Path       string
	Package    string
	Name       string
}

func New(rootModule string, path string) Module {
	packageName := CreatePackageName(rootModule, strings.TrimPrefix(path, "./"))
	parts := strings.Split(packageName, "/")
	name := parts[len(parts)-1]
	return Module{RootModule: rootModule, Path: path, Package: packageName, Name: name}
}

func (m Module) GetPath(filename string) string {
	path := filepath.Join(m.Path, filename)
	return path
}

func (m Module) Submodule(name string) Module {
	if name != "" {
		return New(m.RootModule, filepath.Join(m.Path, name))
	}
	return m
}

func CreatePackageName(args ...string) string {
	parts := []string{}
	for _, arg := range args {
		arg = strings.TrimPrefix(arg, "./")
		if arg != "" {
			parts = append(parts, arg)
		}
	}
	return strings.Join(parts, "/")
}

func GetShortPackageName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
