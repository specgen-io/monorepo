package models

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func addObjectModelMethods(w *sources.Writer, model *spec.NamedModel) {
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
	hashParams := []string{}
	var arrayFieldCount, nonArrayFieldCount int
	for _, field := range model.Object.Fields {
		if isJavaArrayType(&field.Type.Definition) {
			hashCodeParams = append(hashCodeParams, fmt.Sprintf(`%s()`, getterName(&field)))
			arrayFieldCount++
		} else {
			hashParams = append(hashParams, fmt.Sprintf(`%s()`, getterName(&field)))
			nonArrayFieldCount++
		}
	}
	if arrayFieldCount > 0 && nonArrayFieldCount == 0 {
		w.Line(`  int result = Arrays.hashCode(%s);`, hashCodeParams[0])
		for _, param := range hashCodeParams[1:] {
			w.Line(`  result = 31 * result + Arrays.hashCode(%s);`, param)
		}
		w.Line(`  return result;`)
	} else if arrayFieldCount > 0 && nonArrayFieldCount > 0 {
		w.Line(`  int result = Objects.hash(%s);`, joinParams(hashParams))
		for _, param := range hashCodeParams {
			w.Line(`  result = 31 * result + Arrays.hashCode(%s);`, param)
		}
		w.Line(`  return result;`)
	} else if arrayFieldCount == 0 && nonArrayFieldCount > 0 {
		w.Line(`  return Objects.hash(%s);`, joinParams(hashParams))
	}
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`@Override`)
	w.Line(`public String toString() {`)
	formatParams := []string{}
	formatArgs := []string{}
	for _, field := range model.Object.Fields {
		formatParams = append(formatParams, fmt.Sprintf(`%s=%s`, field.Name.CamelCase(), "%s"))
		formatArgs = append(formatArgs, toStringExpression(&field.Type.Definition, field.Name.CamelCase()))
	}
	w.Line(`  return String.format("%s{%s}", %s);`, model.Name.PascalCase(), joinParams(formatParams), joinParams(formatArgs))
	w.Line(`}`)
}

func oneOfItemClassName(item *spec.NamedDefinition) string {
	return item.Name.PascalCase()
}

func addOneOfModelMethods(w *sources.Writer, item *spec.NamedDefinition) {
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
