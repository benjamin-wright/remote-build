package controller

import (
	"context"

	"pongle-hub.co.uk/remote-build/cmd/operator/internal/k8s"
)

type Controller struct {
	client *k8s.Client
}

func New(client *k8s.Client) *Controller {
	return &Controller{
		client: client,
	}
}

func (c *Controller) Start(ctx context.Context) {
	instanceWatcher := c.client.WatchBuildInstances(ctx)
	serviceWatcher := c.client.WatchServices(ctx)
	statefulSetWatcher := c.client.WatchStatefulSets(ctx)

	<-instanceWatcher.Done()
	<-serviceWatcher.Done()
	<-statefulSetWatcher.Done()
}
