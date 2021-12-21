package genscala

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	Path        string
	PackageName string
	PackageStar string
}

func NewPackage(rootPath, rootPackageName, packageName string) Package {
	path := rootPath
	if packageName != "" {
		packagePath := strings.Join(strings.Split(packageName, "."), string(os.PathSeparator))
		path = fmt.Sprintf(`%s/%s`, rootPath, packagePath)
	}
	packageName = getFullPackageName(rootPackageName, packageName)
	return Package{Path: path, PackageName: packageName, PackageStar: packageName + "._"}
}

func getFullPackageName(rootPackageName, packageName string) string {
	if rootPackageName == "" {
		return packageName
	}
	if packageName == "" {
		return rootPackageName
	}
	return fmt.Sprintf(`%s.%s`, rootPackageName, packageName)
}

func (m Package) GetPath(filename string) string {
	path := filepath.Join(m.Path, filename)
	return path
}

func (m Package) Subpackage(name string) Package {
	if name != "" {
		subpackageName := fmt.Sprintf(`%s.%s`, m.PackageName, name)
		subpackagePath := fmt.Sprintf(`%s/%s`, m.Path, name)
		return Package{Path: subpackagePath, PackageName: subpackageName, PackageStar: subpackageName + "._"}
	}
	return m
}
