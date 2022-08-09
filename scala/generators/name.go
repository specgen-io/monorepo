package generators

import (
	"fmt"
	"spec"
)

func scalaCamelCase(name spec.Name) string {
	return scalaName(name.CamelCase())
}

func scalaName(name string) string {
	for _, keyword := range keywords {
		if name == keyword {
			return fmt.Sprintf("`%s`", name)
		}
	}
	return name
}

var keywords = []string{
	"abstract",
	"case",
	"catch",
	"class",
	"def",
	"do",
	"else",
	"extends",
	"false",
	"final",
	"finally",
	"for",
	"forSome",
	"if",
	"implicit",
	"import",
	"lazy",
	"match",
	"new",
	"null",
	"object",
	"override",
	"package",
	"private",
	"protected",
	"return",
	"sealed",
	"super",
	"this",
	"throw",
	"trait",
	"try",
	"true",
	"type",
	"val",
	"var",
	"while",
	"with",
	"yield",
}
