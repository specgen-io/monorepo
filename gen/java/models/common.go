package models

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/generator"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func addObjectModelMethods(w *generator.Writer, model *spec.NamedModel) {
	w.Line(`@Override`)
	w.Line(`public boolean equals(Object o) {`)
	w.Line(`  if (this == o) return true;`)
	w.Line(`  if (!(o instanceof %s)) return false;`, model.Name.PascalCase())
	w.Line(`  %s that = (%s) o;`, model.Name.PascalCase(), model.Name.PascalCase())
	equalsParams := []string{}
	for _, field := range model.Object.Fields {
		equalsParam := equalsExpression(&field.Type.Definition, fmt.Sprintf(`%s()`, getterName(&field)), fmt.Sprintf(`%s()`, getterName(&field)))
		equalsParams = append(equalsParams, equalsParam)
	}
	w.Line(`  return %s;`, strings.Join(equalsParams, " && "))
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`@Override`)
	w.Line(`public int hashCode() {`)
	hashCodeParams := []string{}
	for _, field := range model.Object.Fields {
		hashCodeParams = append(hashCodeParams, fmt.Sprintf(`%s()`, getterName(&field)))
	}
	w.Line(`  return Objects.hash(%s);`, joinParams(hashCodeParams))
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`@Override`)
	w.Line(`public String toString() {`)
	formatParams := []string{}
	formatArgs := []string{}
	for _, field := range model.Object.Fields {
		formatParams = append(formatParams, fmt.Sprintf(`%s=%s`, field.Name.CamelCase(), "%s"))
		formatArgs = append(formatArgs, fmt.Sprintf(`%s`, field.Name.CamelCase()))
	}
	w.Line(`  return String.format("%s{%s}", %s);`, model.Name.PascalCase(), joinParams(formatParams), joinParams(formatArgs))
	w.Line(`}`)
}

func oneOfItemClassName(item *spec.NamedDefinition) string {
	return item.Name.PascalCase()
}

func addOneOfModelMethods(w *generator.Writer, item *spec.NamedDefinition) {
	w.Line(`@Override`)
	w.Line(`public boolean equals(Object o) {`)
	w.Line(`  if (this == o) return true;`)
	w.Line(`  if (!(o instanceof %s)) return false;`, oneOfItemClassName(item))
	w.Line(`  %s that = (%s) o;`, oneOfItemClassName(item), oneOfItemClassName(item))
	w.Line(`  return Objects.equals(getData(), that.getData());`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`@Override`)
	w.Line(`public int hashCode() {`)
	w.Line(`  return Objects.hash(getData());`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`@Override`)
	w.Line(`public String toString() {`)
	w.Line(`  return String.format("%s{data=%s}", data);`, oneOfItemClassName(item), "%s")
	w.Line(`}`)
}
