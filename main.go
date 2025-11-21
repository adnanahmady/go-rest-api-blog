package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/adnanahmady/go-rest-api-blog/internal"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app, err := internal.WireUpApp()
	if err != nil {
		log.Fatalf("failed to wire up app: %v", err)
	}
	defer app.Database.Close()
	app.Database.Migrate()

	go func() {
		if err := app.Server.Start(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	stop()

	if err := app.Server.Shutdown(); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}
	app.Logger.Info("application shutdown completed")
}
