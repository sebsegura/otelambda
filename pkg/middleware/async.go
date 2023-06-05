package middleware

import (
	"context"
	"sebsegura/otelambda/pkg/logger"
)

func StartAsyncLambda[T any](next func(ctx context.Context, evt T) error, opt ...Options) func(ctx context.Context, evt T) error {
	return func(ctx context.Context, evt T) error {
		var endMsg string
		ctx, log := logger.WithLogger(ctx)

		if opt == nil {
			endMsg = _defaultEndMsg
		} else {
			endMsg = opt[0].EndMsg
		}

		log.WithField("event", evt).Info("starting lambda execution")

		err := next(ctx, evt)
		if err != nil {
			ErrorLogger(ctx, err)
			return err
		}

		log.Info(endMsg)

		return err
	}
}
