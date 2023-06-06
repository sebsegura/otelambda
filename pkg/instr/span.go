package instr

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"os"
)

const (
	_lambdaName    = "AWS_LAMBDA_FUNCTION_NAME"
	_lambdaVersion = "AWS_LAMBDA_FUNCTION_VERSION"
)

func GetTracerName() string {
	return fmt.Sprintf("%s:%s", os.Getenv(_lambdaName), os.Getenv(_lambdaVersion))
}

func NewSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return trace.SpanFromContext(ctx).TracerProvider().Tracer(GetTracerName()).Start(ctx, name)
}
