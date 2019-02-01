// cmd/serve/main.go
package main

import (
	"net/http"

	"github.com/emilyzhang/goyagi/pkg/server"
	logger "github.com/lob/logger-go"
)

func main() {
	log := logger.New()

	srv := server.New()

	log.Info("server started")

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Err(err).Fatal("server stopped")
	}

	log.Info("server stopped")
}