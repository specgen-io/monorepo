package gengo

import (
	"fmt"
	"path/filepath"
)

type Package struct {
	root string
	packageName string
}

func NewPackage (root string, packageName string) Package {
	return Package{root, packageName}
}

func (p Package) Import() string {
	return fmt.Sprintf(`"%s/%s"`, p.root, p.packageName)
}

func (p Package) ImportAlias(alias string) string {
	return fmt.Sprintf(`%s %s`, alias, p.Import())
}

func (p Package) Path(filename string) string {
	return filepath.Join(p.packageName, filename)
}