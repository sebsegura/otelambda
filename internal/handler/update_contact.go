package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	jsoniter "github.com/json-iterator/go"
	"sebsegura/otelambda/internal/ddb"
	"sebsegura/otelambda/internal/models"
	"sebsegura/otelambda/pkg/logger"
)

type UpdateContact struct {
	repo ddb.Repository
}

func NewUpdateContactHandler(repo ddb.Repository) *UpdateContact {
	return &UpdateContact{repo: repo}
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (h *UpdateContact) Handle(ctx context.Context, evt events.SQSEvent) error {
	log := logger.New(ctx)

	var contact models.Contact
	if err := json.Unmarshal([]byte(evt.Records[0].Body), &contact); err != nil {
		log.WithField("cause", err.Error()).Error("parse error")
		return err
	}
	log.WithField("contact", contact).Info("event received")

	if err := h.repo.Update(ctx, &contact); err != nil {
		log.WithField("cause", err.Error()).Error("error updating contact")
		return err
	}

	return nil
}
