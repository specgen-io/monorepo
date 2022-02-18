package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gen/ts/modules"
	"github.com/specgen-io/specgen/v2/gen/ts/validations"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

type ClientGenerator interface {
	ApiClient(api spec.Api, validationModule, modelsModule, paramsModule, module modules.Module) *sources.CodeFile
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
