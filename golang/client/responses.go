package client

import (
	"fmt"
	"generator"
	"golang/types"
	"golang/writer"
	"spec"
)

func generateResponseStruct(w generator.Writer, types *types.Types, operation *spec.NamedOperation) {
	w.Line(`type %s struct {`, responseTypeName(operation))
	responses := [][]string{}
	for _, response := range operation.Responses {
		responses = append(responses, []string{
			response.Name.PascalCase(),
			types.GoType(spec.Nullable(&response.Type.Definition)),
		})
	}
	writer.WriteAlignedLines(w.Indented(), responses)
	w.Line(`}`)
}

func responseTypeName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func (g *NetHttpGenerator) GenerateResponseFunctions() *generator.CodeFile {
	w := writer.New(g.Modules.Response, `response.go`)
	w.Lines(`
import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func Json(logFields log.Fields, resp *http.Response, result any) error {
	log.WithFields(logFields).WithField("status", resp.StatusCode).Info("Received response")
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		log.WithFields(logFields).Error("Failed to parse response JSON", err.Error())
		return err
	}
	return nil
}

func Text(logFields log.Fields, resp *http.Response) ([]byte, error) {
	log.WithFields(logFields).WithField("status", resp.StatusCode).Info("Received response")
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		log.WithFields(logFields).Error("Reading request body failed", err.Error())
		return nil, err
	}
	return responseBody, nil
}

func Empty(logFields log.Fields, resp *http.Response) {
	log.WithFields(logFields).WithField("status", resp.StatusCode).Info("Received response")
}
`)
	return w.ToCodeFile()
}
