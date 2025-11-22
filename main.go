package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/adnanahmady/go-rest-api-blog/internal"
)

// @title Go Rest API Blog
// @version 1.0.0
// @description Go Rest API Blog
// @contact.name Adnan Ahmady
// @contact.email adnanahmady@gmail.com
// @contact.url https://github.com/adnanahmady
// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT
func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app, err := internal.WireUpApp()
	if err != nil {
		app.Logger.Fatal("failed to wire up app", err)
	}
	defer app.Database.Close()
	app.Database.Migrate()

	go func() {
		if err := app.Server.Start(); err != nil {
			app.Logger.Fatal("failed to start server", err)
		}
	}()

	<-ctx.Done()
	stop()

	if err := app.Server.Shutdown(); err != nil {
		app.Logger.Fatal("failed to shutdown server", err)
	}
	app.Logger.Info("application shutdown completed")
}
