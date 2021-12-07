package conopenapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/specgen-io/specgen/v2/console"
)

func ConvertFromOpenapi(inFile, outFile string) error {
	doc, err := openapi3.NewLoader().LoadFromFile(inFile)
	if err != nil { return err }
	console.Print(doc)
	return nil
}
