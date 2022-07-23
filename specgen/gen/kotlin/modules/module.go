package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Module struct {
	RootPath    string
	Path        string
	PackageName string
	PackageStar string
}

func Package(rootPath string, packageName string) Module {
	path := fmt.Sprintf(`%s/%s`, rootPath, packageToPath(packageName))
	return Module{RootPath: rootPath, Path: path, PackageName: packageName, PackageStar: packageName + ".*"}
}

func (m Module) GetPath(filename string) string {
	path := filepath.Join(m.Path, filename)
	return path
}

func (m Module) Subpackage(name string) Module {
	if name != "" {
		return Package(m.RootPath, fmt.Sprintf(`%s.%s`, m.PackageName, name))
	}
	return m
}

func (m Module) Get(name string) string {
	return fmt.Sprintf(`%s.%s`, m.PackageName, name)
}

func packageToPath(packageName string) string {
	parts := strings.Split(packageName, ".")
	return strings.Join(parts, string(os.PathSeparator))
}
