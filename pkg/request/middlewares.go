package request

import (
	"bytes"
	"io"
	"net/http"

	"github.com/adnanahmady/go-rest-api-blog/pkg/applog"
	"github.com/go-chi/chi/v5/middleware"
)

func NewMiddlewares(lgr applog.Logger) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		middleware.RequestID,
		NewLoggerMiddleware(lgr),
		middleware.Recoverer,
		middleware.RealIP,
		middleware.RedirectSlashes,
		middleware.StripSlashes,
		middleware.CleanPath,
		middleware.NoCache,
	}
}

type appWriter struct {
	lgr applog.Logger
	http.ResponseWriter
	statusCode int
}

func newWriter(
	lgr applog.Logger,
	w http.ResponseWriter,
) *appWriter {
	return &appWriter{
		lgr: lgr,
		ResponseWriter: w,
		statusCode:     200,
	}
}

func (w *appWriter) WriteHeader(statusCode int) {
	w.lgr.Info("response headers", "status_code", statusCode)
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *appWriter) Write(b []byte) (int, error) {
	w.lgr.Info("response", "body", string(b))
	return w.ResponseWriter.Write(b)
}

func NewLoggerMiddleware(lgr applog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			reqLgr := lgr.With(
				"request_id", middleware.GetReqID(ctx),
				"method", r.Method,
				"path", r.URL.Path,
			)
			reqLgr.Info(
				"request received",
				"remote_addr", r.RemoteAddr, "user_agent", r.UserAgent(),
				"host", r.Host, "referer", r.Referer(),
			)

			body, err := io.ReadAll(r.Body)
			if err != nil {
				reqLgr.Error("failed to read request body", err)
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			reqLgr.Info("request params", "params", r.URL.Query())
			reqLgr.Info("request headers", "headers", r.Header)
			reqLgr.Info("request body", "body", string(body),
				"content_type", r.Header.Get("Content-Type"),
				"content_length", r.ContentLength)

			ctx = WithLogger(ctx, reqLgr)
			r = r.WithContext(ctx)

			next.ServeHTTP(newWriter(reqLgr, w), r)
		})
	}
}
