package service_nw_server

import (
	"context"
	"errors"
	"net"
	"net/http"

	"raiashpanda007/local-cloud-cli/cmd/core/env"
)

func MasterServer(listener net.Listener, centralContext, errorContext context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Handler: mux,
	}

	serverErrors := make(chan error, 1)

	env.Log.Info("http server starting", "addr", listener.Addr().String())

	go func() {
		serverErrors <- server.Serve(listener)
	}()

	var shutdownReason string

	select {
	case <-centralContext.Done():
		shutdownReason = "shutdown signal"
	case <-errorContext.Done():
		shutdownReason = "error context canceled"
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			env.Log.Info("http server stopped", "addr", listener.Addr().String())
			return nil
		}
		env.Log.Error("http server failed", "addr", listener.Addr().String(), "err", err)
		return err
	}

	env.Log.Info("http server shutting down", "addr", listener.Addr().String(), "reason", shutdownReason)

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5_000_000_000,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		env.Log.Error("http server shutdown failed", "err", err)
		return err
	}

	err := <-serverErrors

	if errors.Is(err, http.ErrServerClosed) {
		env.Log.Info("http server stopped", "addr", listener.Addr().String())
		return nil
	}

	env.Log.Error("http server failed", "addr", listener.Addr().String(), "err", err)
	return err
}
