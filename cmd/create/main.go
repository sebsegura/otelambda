package main

import (
	"context"
	"os"
	"sebsegura/otelambda/internal/api"
	"sebsegura/otelambda/internal/ddb"
	"sebsegura/otelambda/internal/handler"
	"sebsegura/otelambda/pkg/instr"
	"sebsegura/otelambda/pkg/middleware"
)

func main() {
	var (
		ctx    = context.Background()
		repo   = ddb.New(ctx)
		client = api.NewAPIClient(os.Getenv("BASE_URL"))
		h      = handler.NewCreateContactHandler(repo, client)
	)

	instr.StartInstrumentedLambda(ctx, middleware.StartSyncLambda[handler.Request, handler.Response](h.Handle))
}
