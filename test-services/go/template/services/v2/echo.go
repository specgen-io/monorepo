package v2

import (
	"the-service/spec/v2/models"
)

type EchoService struct{}

func (service *EchoService) EchoBodyModel(body *models.Message) (*models.Message, error) {
	return body, nil
}
