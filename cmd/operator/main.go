package main

import (
	"context"
	"log/slog"
	"os"

	"pongle-hub.co.uk/remote-build/cmd/operator/internal/k8s"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)))

	slog.Info("Starting operator...")

	client, err := k8s.NewClient(map[string]string{
		"app.kubernetes.io/managed-by": "remote-build-operator",
	})

	if err != nil {
		panic(err)
	}

	instanceWatcher := client.WatchBuildInstances(context.Background())

	slog.Info("Running...")

	<-instanceWatcher.Done()
}
