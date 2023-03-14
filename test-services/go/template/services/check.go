package services

import (
	"the-service/spec/check"
	"the-service/spec/empty"
	"the-service/spec/httperrors/errmodels"
	"the-service/spec/models"
)

type CheckService struct{}

func (service *CheckService) CheckEmpty() error {
	return nil
}

func (service *CheckService) CheckEmptyResponse(body *models.Message) error {
	return nil
}

func (service *CheckService) CheckForbidden() (*check.CheckForbiddenResponse, error) {
	return &check.CheckForbiddenResponse{Forbidden: &empty.Value}, nil
}

func (service *CheckService) SameOperationName() (*check.SameOperationNameResponse, error) {
	return &check.SameOperationNameResponse{Ok: &empty.Value}, nil
}

func (service *CheckService) CheckBadRequest() (*check.CheckBadRequestResponse, error) {
	badRequest := errmodels.BadRequestError{"Error returned from service implementation", errmodels.ErrorLocationUnknown, nil}
	return &check.CheckBadRequestResponse{BadRequest: &badRequest}, nil
}
