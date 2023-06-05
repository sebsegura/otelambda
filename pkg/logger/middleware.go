package logger

import "context"

func WithLogger(ctx context.Context) (context.Context, *Logger) {
	log := New(ctx)
	return context.WithValue(ctx, "logger", log), log
}

func GetLogger(ctx context.Context) *Logger {
	log, _ := ctx.Value("logger").(*Logger)
	return log
}
