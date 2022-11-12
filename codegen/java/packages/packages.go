package packages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	RootPath    string
	Path        string
	PackageName string
	PackageStar string
}

func New(rootPath string, packageName string) Package {
	path := fmt.Sprintf(`%s/%s`, rootPath, packageToPath(packageName))
	return Package{RootPath: rootPath, Path: path, PackageName: packageName, PackageStar: packageName + ".*"}
}

func (m Package) GetPath(filename string) string {
	path := filepath.Join(m.Path, filename)
	return path
}

func (m Package) Subpackage(name string) Package {
	if name != "" {
		return New(m.RootPath, fmt.Sprintf(`%s.%s`, m.PackageName, name))
	}
	return m
}

func packageToPath(packageName string) string {
	parts := strings.Split(packageName, ".")
	return strings.Join(parts, string(os.PathSeparator))
}
