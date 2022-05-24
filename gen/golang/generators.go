package golang

import (
	"github.com/specgen-io/specgen/v2/gen/golang/client"
	"github.com/specgen-io/specgen/v2/gen/golang/models"
	"github.com/specgen-io/specgen/v2/gen/golang/service"
)

var GenerateClient = client.GenerateClient
var GenerateService = service.GenerateService
var GenerateModels = models.GenerateModels
