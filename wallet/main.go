package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//not sure about how i want to handle this part
	//split between printing to standard out and letting another service
	//take care of this. We'll see...
	// _, err := os.OpenFile("logFile", os.O_RDWR|os.O_CREATE, 0606)
	// if err != nil {
	// 	fmt.Println("error setting up logfile")
	// 	return
	// }
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	handlers := handlerConfig{
		logger: logger,
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:              ":8081",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	mux.HandleFunc("GET /healthz", handlers.healthCheck)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("starting server on port :8081")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down gracefully")

	//include an additional context to ensure no new
	//connection is accepted and all existing connections
	//get a chance to shutdown properly

	shutDownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutDownContext); err != nil {
		logger.Error("shutdown failed", "err", err)
	}
	logger.Info("server stopped")

}
