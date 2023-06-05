package services

import (
	"{{project.value}}/spec"
)

func Create() spec.Services {
	sampleService := &SampleService{}
	return Services{sampleService}
}
