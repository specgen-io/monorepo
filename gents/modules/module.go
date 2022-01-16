package modules

import (
	"fmt"
	"strings"
)

type Module struct {
	path string
}

func NewModule(folderPath string) Module {
	return Module{path: folderPath}
}

func (m Module) GetPath() string {
	return fmt.Sprintf(`%s.ts`, m.path)
}

func (m Module) GetImport(toModule Module) string {
	return importPath(m.GetPath(), toModule.GetPath())
}

func (m Module) Submodule(name string) Module {
	if name != "" {
		return NewModule(fmt.Sprintf(`%s/%s`, m.path, name))
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
	prefix := commonPrefixPath(whatPath, toPath)
	pathSegmentsCount := strings.Count(strings.TrimPrefix(toPath, prefix), "/")
	backwardsPath := strings.Repeat("../", pathSegmentsCount)
	result := "./" + backwardsPath + strings.TrimPrefix(whatPath, prefix)
	if strings.HasSuffix(result, ".ts") {
		result = result[:len(result)-3]
	}
	return result
}
