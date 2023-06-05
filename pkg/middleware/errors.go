package middleware

import (
	"context"
	"sebsegura/otelambda/internal/handler"
	"sebsegura/otelambda/pkg/logger"
)

func ErrorLogger(ctx context.Context, err error) {
	log := logger.GetLogger(ctx)
	customErr := err.(*handler.CustomError)

	log.
		WithField("error.code", customErr.Code).
		WithField("error.cause", customErr.Cause).
		WithField("error.detail", customErr.Detail).
		Error(customErr.Message)
}
