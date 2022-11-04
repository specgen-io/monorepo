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
	if client == Axios {
		return &Generator{validation, &axiosGenerator{modules, validation}, modules}
	}
	if client == NodeFetch {
		return &Generator{validation, &fetchGenerator{modules, true, validation}, modules}
	}
	if client == BrowserFetch {
		return &Generator{validation, &fetchGenerator{modules, false, validation}, modules}
	}
	panic(fmt.Sprintf("Unknown client: %s", client))
}

func (g *Generator) ParamsBuilder() *generator.CodeFile {
	return generateParamsBuilder(g.Modules.Params)
}

var Axios = "axios"
var NodeFetch = "node-fetch"
var BrowserFetch = "browser-fetch"

var Clients = []string{Axios, NodeFetch, BrowserFetch}
