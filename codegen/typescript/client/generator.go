package client

import (
	"fmt"
	"generator"
	"spec"
	"typescript/validations"
)

type ClientGenerator interface {
	ApiClient(api *spec.Api) *generator.CodeFile
}

type Generator struct {
	validations.Validation
	ClientGenerator
	Modules *Modules
}

func NewClientGenerator(client, validationName string, modules *Modules) *Generator {
	validation := validations.New(validationName, &(modules.Modules))
	var clientGenerator ClientGenerator = nil
	switch client {
	case Axios:
		clientGenerator = &axiosGenerator{modules, validation}
	case NodeFetch:
		clientGenerator = &fetchGenerator{modules, true, validation}
	case BrowserFetch:
		clientGenerator = &fetchGenerator{modules, false, validation}
	default:
		panic(fmt.Sprintf("Unknown client: %s", client))
	}
	return &Generator{validation, clientGenerator, modules}
}

var Axios = "axios"
var NodeFetch = "node-fetch"
var BrowserFetch = "browser-fetch"

var Clients = []string{Axios, NodeFetch, BrowserFetch}
