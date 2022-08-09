package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	"spec"
	"io/ioutil"
)

func ConvertFromOpenapi(inFile, outFile string) error {
	doc, err := openapi3.NewLoader().LoadFromFile(inFile)
	if err != nil {
		return err
	}
	converter := NewConverter()
	specification := converter.Specification(doc)
	data, err := spec.WriteSpec(specification)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outFile, data, 0644)
	return err
}

type Converter struct {
	errors   []string
	warnings []string
}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) Error(message string) {
	c.errors = append(c.errors, message)
}

func (c *Converter) Warning(message string) {
	c.warnings = append(c.warnings, message)
}
