package handler

import (
	"context"
	"github.com/oklog/ulid/v2"
	"sebsegura/otelambda/internal/api"
	"sebsegura/otelambda/internal/ddb"
	"sebsegura/otelambda/internal/models"
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
	contact := &models.Contact{
		ID:        ulid.Make().String(),
		FirstName: req.Name,
	}
	if err := h.repo.Create(ctx, contact); err != nil {
		return Response{}, &CustomError{
			Code:    _dbError,
			Cause:   "dynamodb error",
			Message: "cannot create a new item",
			Detail:  err.Error(),
		}
	}

	if err := h.api.Post(ctx, contact); err != nil {
		return Response{}, &CustomError{
			Code:    _apiError,
			Cause:   "API communication",
			Message: "cannot POST a new request",
			Detail:  err.Error(),
		}
	}

	return Response{
		Ok: true,
	}, nil
}
