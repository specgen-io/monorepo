package genjava

import (
	"github.com/specgen-io/specgen/v2/genjava/client"
	"github.com/specgen-io/specgen/v2/genjava/models"
	"github.com/specgen-io/specgen/v2/genjava/service"
)

var GenerateClient = client.Generate
var GenerateService = service.Generate
var GenerateModels = models.Generate
