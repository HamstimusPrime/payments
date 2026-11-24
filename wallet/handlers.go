package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

type handlerConfig struct {
	logger *slog.Logger
}

func (h *handlerConfig) healthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("health check", "remote", r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status" : "ok"}` + "\n"))
	fmt.Fprintf(w, "you have reached the server")

}
