package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/gents/modules"
	"github.com/specgen-io/specgen/v2/gents/validation"
	"github.com/specgen-io/specgen/v2/sources"
	"github.com/specgen-io/specgen/v2/spec"
)

type clientGenerator interface {
	generateApiClient(api spec.Api, validationModule, modelsModule, paramsModule, module modules.Module) *sources.CodeFile
}

func newClientGenerator(client string, validation validation.Validation) clientGenerator {
	if client == "axios" {
		return &axiosGenerator{validation}
	}
	if client == "node-fetch" {
		return &fetchGenerator{true, validation}
	}
	if client == "browser-fetch" {
		return &fetchGenerator{false, validation}
	}
	panic(fmt.Sprintf("Unknown client: %s", client))
}
