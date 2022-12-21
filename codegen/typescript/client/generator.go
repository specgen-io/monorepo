package client

import (
	"fmt"
	"generator"
	"spec"
	"typescript/validations"
)

type ClientGenerator interface {
	ApiClient(api *spec.Api) *generator.CodeFile
	ErrorResponses(httpErrors *spec.HttpErrors) *generator.CodeFile
}

type Generator struct {
	validations.Validation
	CommonGenerator
	ClientGenerator
	Modules *Modules
}

func NewClientGenerator(client, validationName string, modules *Modules) *Generator {
	validation := validations.New(validationName, &(modules.Modules))
	commonGenerator := CommonGenerator{modules, validation}
	var clientGenerator ClientGenerator = nil
	switch client {
	case Axios:
		clientGenerator = &AxiosGenerator{modules, validation, commonGenerator}
	case NodeFetch:
		clientGenerator = &FetchGenerator{modules, true, validation, commonGenerator}
	case BrowserFetch:
		clientGenerator = &FetchGenerator{modules, false, validation, commonGenerator}
	default:
		panic(fmt.Sprintf("Unknown client: %s", client))
	}
	return &Generator{validation, commonGenerator, clientGenerator, modules}
}

var Axios = "axios"
var NodeFetch = "node-fetch"
var BrowserFetch = "browser-fetch"

var Clients = []string{Axios, NodeFetch, BrowserFetch}
