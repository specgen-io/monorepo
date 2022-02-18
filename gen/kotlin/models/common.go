package models

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

func addObjectModelMethods(w *sources.Writer, model *spec.NamedModel) {
	w.Line(`override fun equals(other: Any?): Boolean {`)
	w.Line(`  if (this === other) return true`)
	w.Line(`  if (other !is %s) return false`, model.Name.PascalCase())
	w.EmptyLine()
	for _, field := range model.Object.Fields {
		if _isKotlinArrayType(&field.Type.Definition) {
			w.Line(`  if (!%s.contentEquals(other.%s)) return false`, field.Name.CamelCase(), field.Name.CamelCase())
		} else {
			w.Line(`  if (%s != other.%s) return false`, field.Name.CamelCase(), field.Name.CamelCase())
		}
	}
	w.EmptyLine()
	w.Line(`  return true`)
	w.Line(`}`)
	w.EmptyLine()
	w.Line(`override fun hashCode(): Int {`)
	hashCodeArrayParams := []string{}
	hashCodeNonArrayParams := []string{}
	var arrayFieldCount, nonArrayFieldCount int
	for _, field := range model.Object.Fields {
		if _isKotlinArrayType(&field.Type.Definition) {
			hashCodeArrayParams = append(hashCodeArrayParams, fmt.Sprintf(`%s.contentHashCode()`, field.Name.CamelCase()))
			arrayFieldCount++
		} else {
			if field.Type.Definition.IsNullable() {
				hashCodeNonArrayParams = append(hashCodeNonArrayParams, fmt.Sprintf(`(%s?.hashCode() ?: 0)`, field.Name.CamelCase()))
			} else {
				hashCodeNonArrayParams = append(hashCodeNonArrayParams, fmt.Sprintf(`%s.hashCode()`, field.Name.CamelCase()))
			}
			nonArrayFieldCount++
		}
	}
	if arrayFieldCount == 1 && nonArrayFieldCount == 0 {
		w.Line(`  return %s`, hashCodeArrayParams[0])
	} else if arrayFieldCount > 1 && nonArrayFieldCount == 0 {
		w.Line(`  var result = %s`, hashCodeArrayParams[0])
		for _, param := range hashCodeArrayParams[1:] {
			w.Line(`  result = 31 * result + %s`, param)
		}
		w.Line(`  return result`)
	} else if arrayFieldCount > 0 && nonArrayFieldCount > 0 {
		w.Line(`  var result = %s`, hashCodeArrayParams[0])
		for _, param := range hashCodeArrayParams[1:] {
			w.Line(`  result = 31 * result + %s`, param)
		}
		for _, param := range hashCodeNonArrayParams {
			w.Line(`  result = 31 * result + %s`, param)
		}
		w.Line(`  return result`)
	}
	w.Line(`}`)
}

func oneOfItemClassName(item *spec.NamedDefinition) string {
	return item.Name.PascalCase()
}
