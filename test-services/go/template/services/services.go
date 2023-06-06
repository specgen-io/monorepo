package services

import (
	v2 "{{project.value}}/services/v2"
	"{{project.value}}/spec"
)

func Create() spec.Services {
	echoServiceV2 := &v2.EchoService{}
	echoService := &EchoService{}
	checkService := &CheckService{}
	return spec.Services{echoServiceV2, echoService, checkService}
}
