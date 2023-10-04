package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nordluma/go-bookstore/config"
)

// Start listening for HTTP requests
var StartHTTPServer = startHTTPServer

func startHTTPServer() error {
	mux := http.NewServeMux()
	mux.Handle("/api/", newHandlerAPI())

	server := http.Server{
		ReadTimeout:  config.GetHTTPReadTimeout(),
		WriteTimeout: config.GetHTTPWriteTimeout(),
		Addr:         config.GetHTTPServerAddress(),
		Handler:      mux,
	}

	go func() {
		// Listen for OS interrupt signal
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		<-interrupt
		// Gracefully shut down the server
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down. %v\n", err)
		}
	}()

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		err = nil
	}

	return nil
}
