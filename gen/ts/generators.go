package ts

import (
	"github.com/specgen-io/specgen/v2/gen/ts/client"
	"github.com/specgen-io/specgen/v2/gen/ts/service"
	"github.com/specgen-io/specgen/v2/gen/ts/validations"
)

var GenerateClient = client.GenerateClient
var GenerateService = service.GenerateService
var GenerateModels = validations.GenerateModels
