package api

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"sebsegura/otelambda/pkg/instr"
	"sebsegura/otelambda/pkg/logger"
	"strings"
)

type Client interface {
	Post(ctx context.Context, body any) error
}

type client struct {
	rest    *http.Client
	baseURL string
}

func NewAPIClient(baseURL string) Client {
	// Only for calling own services
	rest := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	return &client{
		rest:    rest,
		baseURL: baseURL,
	}
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (c *client) Post(ctx context.Context, body any) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(instr.GetTracerName()).Start(ctx, "POST "+c.baseURL)
	defer span.End()

	log := logger.New(ctx)
	b, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, strings.NewReader(string(b)))
	if err != nil {
		return err
	}
	log.WithField("request.url", req.URL.String()).WithField("request.body", string(b)).Debug("making request")

	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := c.rest.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, "request error")
		return err
	}

	if resp.StatusCode != http.StatusOK {
		span.SetStatus(codes.Error, "wrong status")
		return fmt.Errorf("wrong status %d", resp.StatusCode)
	}

	return nil
}
