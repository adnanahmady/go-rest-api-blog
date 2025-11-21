package request

import (
	"context"

	"github.com/adnanahmady/go-rest-api-blog/pkg/applog"
)

var (
	ctxKeyLogger = &struct{ uint8 }{}
)

func WithLogger(ctx context.Context, lgr applog.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, lgr)
}

func GetLogger(ctx context.Context) applog.Logger {
	return ctx.Value(ctxKeyLogger).(applog.Logger)
}
