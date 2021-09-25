package genjava

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type module struct {
	RootPath     string
	Path         string
	Package      string
	PackageStar  string
}

func Module(rootPath string, packageName string) module {
	path := fmt.Sprintf(`%s/%s`, rootPath, packageToPath(packageName))
	return module{RootPath: rootPath, Path: path, Package: packageName, PackageStar: packageName+".*"}
}

func (m module) GetPath(filename string) string {
	path := filepath.Join(m.Path, filename)
	return path
}

func (m module) Submodule(name string) module {
	if name != "" {
		return Module(m.RootPath, fmt.Sprintf(`%s.%s`, m.Package, name))
	}
	return m
}

func packageToPath(packageName string) string {
	parts := strings.Split(packageName, ".")
	return strings.Join(parts, string(os.PathSeparator))
}