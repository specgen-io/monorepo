package genjava

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	RootPath     string
	Path        string
	PackageName string
	PackageStar string
}

func PackageInfo(rootPath string, packageName string) Package {
	path := fmt.Sprintf(`%s/%s`, rootPath, packageToPath(packageName))
	return Package{RootPath: rootPath, Path: path, PackageName: packageName, PackageStar: packageName+".*"}
}

func (p Package) GetPath(filename string) string {
	path := filepath.Join(p.Path, filename)
	return path
}

func (p Package) Subpackage(name string) Package {
	if name != "" {
		return PackageInfo(p.RootPath, fmt.Sprintf(`%s.%s`, p.PackageName, name))
	}
	return p
}

func packageToPath(packageName string) string {
	parts := strings.Split(packageName, ".")
	return strings.Join(parts, string(os.PathSeparator))
}