package module

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Module struct {
	RootModule string
	Path       string
	Package    string
	Name       string
	Alias      string
}

func New(rootModule string, path string) Module {
	return NewAliased(rootModule, path, "")
}

func NewAliased(rootModule string, path string, alias string) Module {
	packageName := createPackageName(rootModule, strings.TrimPrefix(path, "./"))
	parts := strings.Split(packageName, "/")
	name := parts[len(parts)-1]
	return Module{RootModule: rootModule, Path: path, Package: packageName, Name: name, Alias: alias}
}

func (m Module) GetPath(filename string) string {
	path := filepath.Join(m.Path, filename)
	return path
}

func (m Module) Submodule(name string) Module {
	return m.SubmoduleAliased(name, "")
}

func (m Module) SubmoduleAliased(name string, alias string) Module {
	if name != "" {
		return NewAliased(m.RootModule, filepath.Join(m.Path, name), alias)
	}
	return m
}

func (m Module) Use() string {
	if m.Alias != "" {
		return m.Alias
	}
	return m.Name
}

func (m Module) Get(name string) string {
	return fmt.Sprintf("%s.%s", m.Use(), name)
}

func createPackageName(args ...string) string {
	parts := []string{}
	for _, arg := range args {
		arg = strings.TrimPrefix(arg, "./")
		if arg != "" {
			parts = append(parts, arg)
		}
	}
	return strings.Join(parts, "/")
}
