package middleware

import (
	"context"
	"sebsegura/otelambda/pkg/logger"
)

func StartSyncLambda[X, Y any](next func(ctx context.Context, req X) (Y, error), opt ...Options) func(ctx context.Context, req X) (Y, error) {
	return func(ctx context.Context, req X) (Y, error) {
		var endMsg string
		ctx, log := logger.WithLogger(ctx)

		if opt == nil {
			endMsg = _defaultEndMsg
		} else {
			endMsg = opt[0].EndMsg
		}

		log.WithField("event", req).Info("starting lambda execution")

		resp, err := next(ctx, req)
		if err != nil {
			ErrorLogger(ctx, err)
			return resp, err
		}

		log.WithField("event.response", resp).Info(endMsg)

		return resp, err
	}
}
