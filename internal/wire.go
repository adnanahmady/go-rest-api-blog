//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/adnanahmady/go-rest-api-blog/config"
	"github.com/adnanahmady/go-rest-api-blog/internal/application"
	"github.com/adnanahmady/go-rest-api-blog/internal/domain"
	"github.com/adnanahmady/go-rest-api-blog/internal/infra"
	"github.com/adnanahmady/go-rest-api-blog/internal/presentation"
	"github.com/adnanahmady/go-rest-api-blog/pkg/applog"
	"github.com/adnanahmady/go-rest-api-blog/pkg/database"
	"github.com/adnanahmady/go-rest-api-blog/pkg/request"
	"github.com/google/wire"
)

type App struct {
	Config   *config.Config
	Logger   *applog.AppLog
	Database *database.DatabaseManagerImpl
	Server   *request.ServerImpl
	Rotues   *presentation.V1Routes
}

var AppSet = wire.NewSet(
	config.GetConfig,

	request.NewServer,
	wire.Bind(new(request.Router), new(*request.ServerImpl)),

	applog.NewLog,
	wire.Bind(new(applog.Logger), new(*applog.AppLog)),

	database.NewDatabaseManager,
	wire.Bind(new(database.DatabaseManager), new(*database.DatabaseManagerImpl)),

	presentation.NewV1Handlers,
	presentation.NewV1Routes,

	infra.NewSqlitePostRepository,
	wire.Bind(new(domain.PostRepository), new(*infra.SqlitePostRepository)),

	application.NewCreatePostUseCase,
	wire.Bind(new(application.CreatePostUseCase), new(*application.CreatePostUseCaseImpl)),
	application.NewListPostsUseCase,
	wire.Bind(new(application.ListPostsUseCase), new(*application.ListPostsUseCaseImpl)),
	application.NewShowPostUseCase,
	wire.Bind(new(application.ShowPostUseCase), new(*application.ShowPostUseCaseImpl)),
	application.NewUpdatePostUseCase,
	wire.Bind(new(application.UpdatePostUseCase), new(*application.UpdatePostUseCaseImpl)),
	application.NewDeletePostUseCase,
	wire.Bind(new(application.DeletePostUseCase), new(*application.DeletePostUseCaseImpl)),

	wire.Struct(new(App), "*"),
)

func WireUpApp() (*App, error) {
	wire.Build(AppSet)
	return &App{}, nil
}
