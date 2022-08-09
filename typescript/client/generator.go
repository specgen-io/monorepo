package client

import (
	"fmt"

	"generator"
	"spec"
	"typescript/modules"
	"typescript/validations"
)

type ClientGenerator interface {
	ApiClient(api spec.Api, validationModule, modelsModule, paramsModule, module modules.Module) *generator.CodeFile
}

func NewClientGenerator(client string, validation validations.Validation) ClientGenerator {
	if client == Axios {
		return &axiosGenerator{validation}
	}
	if client == NodeFetch {
		return &fetchGenerator{true, validation}
	}
	if client == BrowserFetch {
		return &fetchGenerator{false, validation}
	}
	panic(fmt.Sprintf("Unknown client: %s", client))
}

var Axios = "axios"
var NodeFetch = "node-fetch"
var BrowserFetch = "browser-fetch"

var Clients = []string{Axios, NodeFetch, BrowserFetch}
