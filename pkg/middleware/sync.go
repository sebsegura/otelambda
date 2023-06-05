package middleware

import (
	"context"
	"sebsegura/otelambda/pkg/logger"
)

func StartSyncLambda[X, Y any](next func(ctx context.Context, req X) (Y, error)) func(ctx context.Context, req X) (Y, error) {
	return func(ctx context.Context, req X) (Y, error) {
		ctx, log := logger.WithLogger(ctx)
		log.WithField("event", req).Info("starting lambda execution")

		resp, err := next(ctx, req)
		if err != nil {
			ErrorLogger(ctx, err)
			return resp, err
		}

		log.WithField("event.response", resp).Info("end of execution")

		return resp, err
	}
}
