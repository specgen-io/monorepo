package module

import (
	"fmt"
	"strings"
)

type Module struct {
	path    string
	isIndex bool
}

func New(folderPath string) Module {
	return Module{path: folderPath, isIndex: false}
}

func NewIndex(folderPath string) Module {
	return Module{path: folderPath, isIndex: true}
}

func (m Module) ItIsIndex() {
	m.isIndex = true
}

func (m Module) GetPath() string {
	if m.isIndex {
		return fmt.Sprintf(`%s/index.ts`, m.path)
	} else {
		return fmt.Sprintf(`%s.ts`, m.path)
	}
}

func (m Module) GetImport(toModule Module) string {
	return importPath(m.GetPath(), toModule.GetPath())
}

func (m Module) Submodule(name string) Module {
	if name != "" {
		return New(fmt.Sprintf(`%s/%s`, m.path, name))
	}
	return m
}

func (m Module) SubmoduleIndex(name string) Module {
	if name != "" {
		return NewIndex(fmt.Sprintf(`%s/%s`, m.path, name))
	}
	return m
}

func commonPrefixPath(s1, s2 string) string {
	result := ""
	p1 := strings.Split(s1, "/")
	p2 := strings.Split(s2, "/")
	minLength := len(p2)
	if len(p1) < minLength {
		minLength = len(p1)
	}
	for i := 0; i < minLength; i++ {
		if p1[i] != p2[i] {
			break
		}
		result += p1[i] + "/"
	}
	return result
}

func importPath(whatPath string, toPath string) string {
	if strings.HasSuffix(whatPath, "/index.ts") {
		whatPath = strings.TrimSuffix(whatPath, "/index.ts")
	}
	prefix := commonPrefixPath(whatPath, toPath)
	pathSegmentsCount := strings.Count(strings.TrimPrefix(toPath, prefix), "/")
	backwardsPath := strings.Repeat("../", pathSegmentsCount)
	result := "./" + backwardsPath + strings.TrimPrefix(whatPath, prefix)
	if strings.HasSuffix(result, ".ts") {
		result = result[:len(result)-3]
	}
	return result
}
