package gents

import (
	"github.com/specgen-io/specgen/v2/gents/client"
	"github.com/specgen-io/specgen/v2/gents/service"
	"github.com/specgen-io/specgen/v2/gents/validations"
)

var GenerateClient = client.GenerateClient
var GenerateService = service.GenerateService
var GenerateModels = validations.GenerateModels
