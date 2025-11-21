package request

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/adnanahmady/go-rest-api-blog/config"
	_ "github.com/adnanahmady/go-rest-api-blog/docs"
	"github.com/adnanahmady/go-rest-api-blog/pkg/applog"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Router interface {
	GetEngine() *chi.Mux
}

type ServerImpl struct {
	cfg    *config.Config
	lgr    applog.Logger
	mux    *chi.Mux
	server *http.Server
}

func NewServer(cfg *config.Config, lgr applog.Logger) *ServerImpl {
	return &ServerImpl{
		cfg:    cfg,
		lgr:    lgr,
		mux:    prepareFramework(cfg, lgr),
		server: &http.Server{},
	}
}

func prepareFramework(cfg *config.Config, lgr applog.Logger) *chi.Mux {
	mux := chi.NewMux()
	mux.Use(NewMiddlewares(lgr)...)

	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", cfg.App.Host, cfg.App.Port)),
	))

	return mux
}

func (s *ServerImpl) GetEngine() *chi.Mux {
	return s.mux
}

func (s *ServerImpl) Start() error {
	host := fmt.Sprintf("%s:%d", s.cfg.App.Host, s.cfg.App.Port)
	s.server.Addr = host
	s.server.Handler = s.mux
	s.lgr.Info("starting server", "host", host)

	if err := s.server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		s.lgr.Error("failed to start server", err)
		return err
	}
	return nil
}

func (s *ServerImpl) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.lgr.Info("shutting down server")
	if err := s.server.Shutdown(ctx); err != nil {
		s.lgr.Error("failed to shutdown server", err)
		return err
	}
	s.lgr.Info("server shutdown completed")
	return nil
}
