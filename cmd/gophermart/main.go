package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sandor-clegane/go-market/internal/app"
	"github.com/sandor-clegane/go-market/internal/config"
)

func main() {
	var cfg config.Config
	cfg.Init()

	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	idleConnectionsClosed := make(chan struct{})

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals

		if err := a.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	if err = a.Start(); err != http.ErrServerClosed && err != nil {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnectionsClosed
}
