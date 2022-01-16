package gents

import (
	"github.com/specgen-io/specgen/v2/gents/client"
	"github.com/specgen-io/specgen/v2/gents/service"
	"github.com/specgen-io/specgen/v2/gents/validation"
)

var GenerateClient = client.GenerateClient
var GenerateService = service.GenerateService
var GenerateModels = validation.GenerateModels
