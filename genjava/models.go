package genjava

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
	"strings"
)

func GenerateModels(specification *spec.Spec, packageName string, generatePath string, jsonlib string) *sources.Sources {
	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}
	mainPackage := Package(generatePath, packageName)
	modelsPackage := mainPackage.Subpackage("models")

	newSources := sources.NewSources()
	newSources.AddGeneratedAll(generateModels(specification, modelsPackage, jsonlib))

	return newSources
}

func generateModels(specification *spec.Spec, thePackage Module, jsonlib string) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, version := range specification.Versions {
		versionPackage := thePackage.Subpackage(version.Version.FlatCase())
		files = append(files, generateVersionModels(&version, versionPackage, jsonlib)...)
	}
	files = append(files, *generateJson(thePackage))
	return files
}

func generateVersionModels(version *spec.Version, thePackage Module, jsonlib string) []sources.CodeFile {
	files := []sources.CodeFile{}
	for _, model := range version.ResolvedModels {
		if model.IsObject() {
			if jsonlib == Jackson {
				files = append(files, *generateJacksonObjectModel(model, thePackage, jsonlib))
			}
			if jsonlib == Moshi {
				files = append(files, *generateMoshiObjectModel(model, thePackage, jsonlib))
			}
		} else if model.IsOneOf() {
			if jsonlib == Jackson {
				files = append(files, *generateJacksonOneOfModels(model, thePackage, jsonlib))
			}
			if jsonlib == Moshi {
				files = append(files, *generateMoshiOneOfModels(model, thePackage, jsonlib))
			}
		} else if model.IsEnum() {
			files = append(files, *generateEnumModel(model, thePackage))
		}
	}
	return files
}

func addObjectModelProperties(w *sources.Writer, jsonlib string, model *spec.NamedModel) {
	for _, field := range model.Object.Fields {
		w.EmptyLine()
		w.Line(`public %s %s() {`, JavaType(&field.Type.Definition, jsonlib), getterName(&field))
		w.Line(`  return %s;`, field.Name.CamelCase())
		w.Line(`}`)
		w.EmptyLine()
		w.Line(`public void %s(%s %s) {`, setterName(&field), JavaType(&field.Type.Definition, jsonlib), field.Name.CamelCase())
		w.Line(`  this.%s = %s;`, field.Name.CamelCase(), field.Name.CamelCase())
		w.Line(`}`)
	}
}

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
		w.Line(`  int result = Objects.hash(%s);`, JoinDelimParams(hashParams))
		for _, param := range hashCodeParams {
			w.Line(`  result = 31 * result + Arrays.hashCode(%s);`, param)
		}
		w.Line(`  return result;`)
	} else if arrayFieldCount == 0 && nonArrayFieldCount > 0 {
		w.Line(`  return Objects.hash(%s);`, JoinDelimParams(hashParams))
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
	w.Line(`  return String.format("%s{%s}", %s);`, model.Name.PascalCase(), JoinDelimParams(formatParams), JoinDelimParams(formatArgs))
	w.Line(`}`)
}

func generateEnumModel(model *spec.NamedModel, thePackage Module) *sources.CodeFile {
	w := NewJavaWriter()
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	addJacksonImports(w)
	w.EmptyLine()
	enumName := model.Name.PascalCase()
	w.Line(`public enum %s {`, enumName)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @JsonProperty("%s") %s,`, enumItem.Value, enumItem.Name.UpperCase())
	}
	w.Line(`}`)

	return &sources.CodeFile{
		Path:    thePackage.GetPath(enumName + ".java"),
		Content: w.String(),
	}
}

func oneOfItemClassName(item *spec.NamedDefinition) string {
	return item.Name.PascalCase()
}

func addOneOfModelMethods(w *sources.Writer, jsonlib string, item *spec.NamedDefinition) {
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
	w.Line(`  return String.format("%s{data=%s}", data);`, JavaType(&item.Type.Definition, jsonlib), "%s")
	w.Line(`}`)
}
