package client

import (
	"fmt"
	"kotlin/writer"
	"spec"
)

func addBuilderParam(param *spec.NamedParam) string {
	if param.Type.Definition.IsNullable() {
		return fmt.Sprintf(`%s!!`, param.Name.CamelCase())
	}
	return param.Name.CamelCase()
}

func generateThrowClientException(w *writer.Writer, errorMessage string, wrapException string) {
	w.Line(`val errorMessage = %s`, errorMessage)
	w.Line(`logger.error(errorMessage)`)
	params := "errorMessage"
	if wrapException != "" {
		params += ", " + wrapException
	}
	w.Line(`throw ClientException(%s)`, params)
}
