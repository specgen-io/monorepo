package gents

import (
	"fmt"
	"strings"
)

type module struct {
	path        string
}

func Module(folderPath string) module {
	return module{path: folderPath}
}

func (m module) GetPath() string {
	return fmt.Sprintf(`%s.ts`, m.path)
}

func (m module) GetImport(toModule module) string {
	return importPath(m.GetPath(), toModule.GetPath())
}

func (m module) Submodule(name string) module {
	if name != "" {
		return Module(fmt.Sprintf(`%s/%s`, m.path, name))
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