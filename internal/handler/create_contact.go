package handler

import (
	"context"
	"github.com/oklog/ulid/v2"
	"sebsegura/otelambda/internal/api"
	"sebsegura/otelambda/internal/ddb"
	"sebsegura/otelambda/internal/models"
	"sebsegura/otelambda/pkg/logger"
)

type (
	CreateContact struct {
		repo ddb.Repository
		api  api.Client
	}
	Request struct {
		Name   string `json:"name"`
		Status bool   `json:"status,omitempty"`
	}
	Response struct {
		Ok bool `json:"ok"`
	}
)

func NewCreateContactHandler(repo ddb.Repository, api api.Client) *CreateContact {
	return &CreateContact{
		repo: repo,
		api:  api,
	}
}

func (h *CreateContact) Handle(ctx context.Context, req Request) (Response, error) {
	ok := true
	log := logger.New(ctx)
	log.WithField("request", req).Info("start")

	id := ulid.Make().String()
	contact := &models.Contact{
		ID:        id,
		FirstName: req.Name,
	}
	if err := h.repo.Create(ctx, contact); err != nil {
		log.WithField("cause", err.Error()).Error("error creating new contact")
		ok = false
	}

	if err := h.api.Post(ctx, contact); err != nil {
		log.WithField("cause", err.Error()).Error("error sending request")
	}

	return Response{
		Ok: ok,
	}, nil
}
