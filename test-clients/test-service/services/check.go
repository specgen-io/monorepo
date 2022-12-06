package services

import (
	"test-service/spec/check"
	"test-service/spec/empty"
	"test-service/spec/httperrors/errmodels"
	"test-service/spec/models"
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

func (service *CheckService) CheckConflict() (*check.CheckConflictResponse, error) {
	return &check.CheckConflictResponse{Conflict: &errmodels.ConflictMessage{Message: "Conflict with the current state of the target resource"}}, nil
}

func (service *CheckService) CheckBadRequest() (*check.CheckBadRequestResponse, error) {
	return &check.CheckBadRequestResponse{BadRequest: &errmodels.BadRequestError{Location: errmodels.ErrorLocationUnknown, Message: "Failed to execute request", Errors: nil}}, nil
}

func (service *CheckService) SameOperationName() (*check.SameOperationNameResponse, error) {
	return &check.SameOperationNameResponse{Ok: &empty.Value}, nil
}
