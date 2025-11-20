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

	// TODO: Set up repository, use cases, and HTTP server...

	<-ctx.Done()
	stop()

	// TODO: Shutdown logic...
}
