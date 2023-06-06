package middleware

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
)

func StartOTELSyncLambda[X, Y any](ctx context.Context, handlerFn func(ctx context.Context, req X) (Y, error), opts ...Options) {
	tp, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		fmt.Printf("error creating tracer provider: %v", err)
	}

	defer func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			fmt.Printf("error shutting down tracer provider: %v", err)
		}
	}(ctx)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	lambda.Start(otellambda.InstrumentHandler(StartSyncLambda[X, Y](handlerFn, opts...), xrayconfig.WithRecommendedOptions(tp)...))
}

func StartOTELAsyncLambda[T any](ctx context.Context, handlerFn func(ctx context.Context, evt T) error, opts ...Options) {
	tp, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		fmt.Printf("error creating tracer provider: %v", err)
	}

	defer func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			fmt.Printf("error shutting down tracer provider: %v", err)
		}
	}(ctx)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	lambda.Start(otellambda.InstrumentHandler(StartAsyncLambda[T](handlerFn, opts...), xrayconfig.WithRecommendedOptions(tp)...))
}
