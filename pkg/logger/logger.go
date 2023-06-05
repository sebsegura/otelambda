package logger

import (
	"context"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"sebsegura/otelambda/pkg/env"
	"time"
)

var (
	_defaultWriter io.Writer = os.Stdout
)

type Logger struct {
	log *zap.Logger
}

func New(ctx context.Context) *Logger {
	var (
		cfg   = env.GetVars()
		level zapcore.Level
	)

	switch cfg.Env {
	case "prod":
		if !cfg.Debug {
			level = zap.DebugLevel
		} else {
			level = zap.InfoLevel
		}
	case "dev":
		level = zap.DebugLevel
	}

	var serviceARN string
	if ctx != nil {
		ct, _ := lambdacontext.FromContext(ctx)
		serviceARN = ct.InvokedFunctionArn
	}

	encoderConfig := zapcore.EncoderConfig{
		LevelKey:    "level",
		TimeKey:     "timestamp",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(time.RFC3339))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(_defaultWriter), level)
	log := zap.New(core)
	log = log.
		With(zap.String("service.name", os.Getenv("AWS_LAMBDA_FUNCTION_NAME"))).
		With(zap.String("service.version", os.Getenv("AWS_LAMBDA_FUNCTION_VERSION"))).
		With(zap.String("service.arn", serviceARN)).
		With(zap.String("env", cfg.Env)).
		With(zap.String("trace.id", os.Getenv("_X_AMZN_TRACE_ID"))).
		With(zap.String("event.flow", os.Getenv("FLOW")))

	return &Logger{
		log: log,
	}
}

func (z *Logger) Debug(msg string) {
	z.log.Debug(msg)
}

func (z *Logger) Info(msg string) {
	z.log.Info(msg)
}

func (z *Logger) Error(msg string) {
	z.log.Error(msg)
}

func (z *Logger) WithField(key string, v any) *Logger {
	return &Logger{
		log: z.log.With(zap.Any(key, v)),
	}
}
