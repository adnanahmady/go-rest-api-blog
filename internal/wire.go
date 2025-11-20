//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/adnanahmady/go-rest-api-blog/config"
	"github.com/adnanahmady/go-rest-api-blog/pkg/applog"
	"github.com/adnanahmady/go-rest-api-blog/pkg/database"
	"github.com/google/wire"
)

type App struct {
	Config   *config.Config
	Logger   *applog.AppLog
	Database *database.DatabaseManagerImpl
}

var AppSet = wire.NewSet(
	config.GetConfig,

	applog.NewLog,
	wire.Bind(new(applog.Logger), new(*applog.AppLog)),

	database.NewDatabaseManager,
	wire.Bind(new(database.DatabaseManager), new(*database.DatabaseManagerImpl)),

	wire.Struct(new(App), "*"),
)

func WireUpApp() (*App, error) {
	wire.Build(AppSet)
	return &App{}, nil
}
