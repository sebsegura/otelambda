package middleware

import (
	"context"
	"sebsegura/otelambda/pkg/logger"
)

func StartAsyncLambda[T any](next func(ctx context.Context, evt T) error) func(ctx context.Context, evt T) error {
	return func(ctx context.Context, evt T) error {
		ctx, log := logger.WithLogger(ctx)
		log.WithField("event", evt).Info("starting lambda execution")

		err := next(ctx, evt)
		if err != nil {
			ErrorLogger(ctx, err)
			return err
		}

		log.Info("end of execution")

		return err
	}
}
