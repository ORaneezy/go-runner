package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ORaneezy/go-runner/internal"
)

func main() {
	api := internal.InitAPI()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := api.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start api: %v", err)
		}
	}()

	<-ctx.Done()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := api.Stop(timeoutCtx); err != nil {
		log.Fatalf("failed to stop api: %v", err)
	}

	log.Println("api gracefully stopped")
}
