package client

import (
	"generator"
	"golang/module"
	"strings"
)

func generateResponseFunctions(module module.Module) *generator.CodeFile {
	code := `
package [[.PackageName]]

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
`

	code, _ = generator.ExecuteTemplate(code, struct{ PackageName string }{module.Name})
	return &generator.CodeFile{module.GetPath("response.go"), strings.TrimSpace(code)}
}
