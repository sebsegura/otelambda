package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	jsoniter "github.com/json-iterator/go"
	"sebsegura/otelambda/internal/ddb"
	"sebsegura/otelambda/internal/models"
)

type UpdateContact struct {
	repo ddb.Repository
}

func NewUpdateContactHandler(repo ddb.Repository) *UpdateContact {
	return &UpdateContact{repo: repo}
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (h *UpdateContact) Handle(ctx context.Context, evt events.SQSEvent) error {
	var contact models.Contact
	if err := json.Unmarshal([]byte(evt.Records[0].Body), &contact); err != nil {
		return &CustomError{
			Code:    _validationError,
			Cause:   "validation error",
			Message: "invalid event format",
			Detail:  err.Error(),
		}
	}

	if err := h.repo.Update(ctx, &contact); err != nil {
		return &CustomError{
			Code:    _dbError,
			Cause:   "dynamodb",
			Message: "cannot update a contact",
			Detail:  err.Error(),
		}
	}

	return nil
}
