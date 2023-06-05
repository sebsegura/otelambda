package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"sebsegura/otelambda/internal/ddb"
	"sebsegura/otelambda/internal/handler"
	"sebsegura/otelambda/pkg/middleware"
)

func main() {
	ctx := context.Background()
	repo := ddb.New(ctx)
	h := handler.NewUpdateContactHandler(repo)
	middleware.StartOTELAsyncLambda[events.SQSEvent](ctx, h.Handle)
}
