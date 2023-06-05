package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"os"
	"sebsegura/otelambda/internal/api"
	"sebsegura/otelambda/internal/ddb"
	"sebsegura/otelambda/internal/handler"
	"sebsegura/otelambda/pkg/middleware"
)

func main() {
	ctx := context.Background()
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

	repo := ddb.New(ctx)
	client := api.NewAPIClient(os.Getenv("BASE_URL"))
	h := handler.NewCreateContactHandler(repo, client)

	lambda.Start(otellambda.InstrumentHandler(middleware.StartSyncLambda[handler.Request, handler.Response](h.Handle), xrayconfig.WithRecommendedOptions(tp)...))
}
