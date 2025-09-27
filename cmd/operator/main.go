package main

import (
	"context"
	"log/slog"
	"os"

	"pongle-hub.co.uk/remote-build/cmd/operator/internal/controller"
	"pongle-hub.co.uk/remote-build/cmd/operator/internal/k8s"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)))

	slog.Info("Starting operator...")
	client, err := k8s.NewClient(map[string]string{
		"app.kubernetes.io/managed-by": "remote-build",
	})
	if err != nil {
		slog.Error("Failed to create k8s client", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	ctrl := controller.New(client)
	go func() {
		defer cancel()
		ctrl.Start(ctx)
	}()

	slog.Info("Running...")
	<-ctx.Done()
}
