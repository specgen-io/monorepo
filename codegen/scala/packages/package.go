package packages

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

func New(rootPath, rootPackageName, packageName string) Package {
	path := getFullPath(rootPath, packageNameToPath(packageName))
	packageName = getFullPackageName(rootPackageName, packageName)
	return Package{Path: path, PackageName: packageName, PackageStar: packageName + "._"}
}

func (m Package) GetPath(filename string) string {
	path := filepath.Join(m.Path, filename)
	return path
}

func (m Package) Subpackage(name string) Package {
	subpackageName := getFullPackageName(m.PackageName, name)
	subpackagePath := getFullPath(m.Path, packageNameToPath(name))
	return Package{Path: subpackagePath, PackageName: subpackageName, PackageStar: subpackageName + "._"}
}

func packageNameToPath(packageName string) string {
	return strings.Join(strings.Split(packageName, "."), string(os.PathSeparator))
}

func getFullPackageName(first, second string) string {
	if first == "" {
		return second
	}
	if second == "" {
		return first
	}
	return fmt.Sprintf(`%s.%s`, first, second)
}

func getFullPath(first, second string) string {
	if first == "" {
		return second
	}
	if second == "" {
		return first
	}
	return fmt.Sprintf(`%s%s%s`, first, string(os.PathSeparator), second)
}
