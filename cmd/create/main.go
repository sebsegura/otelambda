package main

import (
	"context"
	"os"
	"sebsegura/otelambda/internal/api"
	"sebsegura/otelambda/internal/ddb"
	"sebsegura/otelambda/internal/handler"
	"sebsegura/otelambda/pkg/middleware"
)

func main() {
	ctx := context.Background()
	repo := ddb.New(ctx)
	client := api.NewAPIClient(os.Getenv("BASE_URL"))
	h := handler.NewCreateContactHandler(repo, client)
	middleware.StartOTELSyncLambda[handler.Request, handler.Response](ctx, h.Handle)
}
