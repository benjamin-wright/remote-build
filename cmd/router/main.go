package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
}

func serve(ctx context.Context) <-chan error {
	router := gin.Default()

	router.GET("/metrics", handlers.PostMetricsHandler(client, cfg))

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router.Handler(),
	}

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	go func() {
		<-ctx.Done()
		done <- srv.Close()
	}()

	return done
}
