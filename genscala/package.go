package genscala

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

func NewPackage(rootPath, rootPackageName, packageName string) Package {
	path := rootPath
	if packageName == "" {
		packagePath := strings.Join(strings.Split(packageName, "."), string(os.PathSeparator))
		path = fmt.Sprintf(`%s/%s`, rootPath, packagePath)
	}
	if rootPackageName != "" {
		packageName = fmt.Sprintf(`%s.%s`, rootPackageName, packageName)
	}
	if packageName == "" {
		packageName = rootPackageName
	}
	return Package{RootPath: rootPath, Path: path, PackageName: packageName, PackageStar: packageName + "._"}
}

func (m Package) GetPath(filename string) string {
	path := filepath.Join(m.Path, filename)
	return path
}

func (m Package) Subpackage(name string) Package {
	if name != "" {
		subpackageName := fmt.Sprintf(`%s.%s`, m.PackageName, name)
		subpackagePath := fmt.Sprintf(`%s/%s`, m.Path, name)
		return Package{RootPath: m.RootPath, Path: subpackagePath, PackageName: subpackageName, PackageStar: subpackageName + "._"}
	}
	return m
}
