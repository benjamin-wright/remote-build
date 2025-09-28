package controller

import (
	"context"
	"log/slog"
	"time"

	"pongle-hub.co.uk/remote-build/cmd/operator/internal/controller/actions"
	"pongle-hub.co.uk/remote-build/cmd/operator/internal/k8s"
	"pongle-hub.co.uk/remote-build/cmd/operator/internal/model"
)

const DEBOUNCE_DURATION = time.Second * 5

type Controller struct {
	client       *k8s.Client
	instances    Store[*model.BuildInstance]
	services     Store[*model.Service]
	statefulSets Store[*model.StatefulSet]
}

func New(client *k8s.Client) *Controller {
	return &Controller{
		client:       client,
		instances:    NewStore[*model.BuildInstance](),
		services:     NewStore[*model.Service](),
		statefulSets: NewStore[*model.StatefulSet](),
	}
}

func (c *Controller) Start(ctx context.Context) {
	instanceWatcher := c.client.WatchBuildInstances(ctx)
	serviceWatcher := c.client.WatchServices(ctx)
	statefulSetWatcher := c.client.WatchStatefulSets(ctx)
	debounce := time.After(DEBOUNCE_DURATION)

	for {
		select {
		case <-ctx.Done():
			return
		case <-debounce:
			slog.Info("Reconciling state")
			numErrors := c.reconcile(ctx)
			if numErrors > 0 {
				slog.Error("Reconciliation failed", "errors", numErrors)
			} else {
				slog.Info("Reconciliation complete")
			}
		case event, ok := <-instanceWatcher.Events():
			if !ok {
				slog.Info("Instance watcher closed, stopping controller")
				return
			}

			c.instances.ProcessEvent(event)
			debounce = time.After(DEBOUNCE_DURATION)
		case event, ok := <-serviceWatcher.Events():
			if !ok {
				slog.Info("Service watcher closed, stopping controller")
				return
			}

			c.services.ProcessEvent(event)
			debounce = time.After(DEBOUNCE_DURATION)
		case event, ok := <-statefulSetWatcher.Events():
			if !ok {
				slog.Info("StatefulSet watcher closed, stopping controller")
				return
			}

			c.statefulSets.ProcessEvent(event)
			debounce = time.After(DEBOUNCE_DURATION)
		}
	}
}

func (c *Controller) reconcile(ctx context.Context) int {
	actions := actions.NewActionMap()

	// Reconcile services
	actions.Append(c.reconcileServices()...)
	actions.Append(c.reconcileOrphanedServices()...)
	actions.Append(c.reconcileStatefulSets()...)
	actions.Append(c.reconcileOrphanedStatefulSets()...)

	return actions.Run(ctx, 5)
}

func (c *Controller) reconcileServices() []*actions.ActionNode {
	var nodes []*actions.ActionNode

	for id, instance := range c.instances.Map() {
		_, serviceExists := c.services.Get(id)

		if !serviceExists {
			nodes = append(nodes, &actions.ActionNode{
				Action: actions.Action{
					Name: "Create Service " + id,
					Do: func(ctx context.Context) error {
						return c.client.CreateService(ctx, &model.Service{
							Name:      instance.Name,
							Namespace: instance.Namespace,
						})
					},
				},
			})
		}
	}

	return nodes
}

func (c *Controller) reconcileOrphanedServices() []*actions.ActionNode {
	var nodes []*actions.ActionNode

	for id, service := range c.services.Map() {
		_, instanceExists := c.instances.Get(id)

		if !instanceExists {
			nodes = append(nodes, &actions.ActionNode{
				Action: actions.Action{
					Name: "Delete Service " + id,
					Do: func(ctx context.Context) error {
						return c.client.DeleteService(ctx, service)
					},
				},
			})
		}
	}

	return nodes
}

func (c *Controller) reconcileStatefulSets() []*actions.ActionNode {
	var nodes []*actions.ActionNode

	for id, instance := range c.instances.Map() {
		existing, statefulSetExists := c.statefulSets.Get(id)
		desired := &model.StatefulSet{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Image:     instance.Image,
			CPU:       instance.CPU,
			Memory:    instance.Memory,
			Disk:      instance.Disk,
		}

		if !statefulSetExists {
			nodes = append(nodes, &actions.ActionNode{
				Action: actions.Action{
					Name: "Create StatefulSet " + id,
					Do: func(ctx context.Context) error {
						return c.client.CreateStatefulSet(ctx, desired)
					},
				},
			})

			continue
		}

		// Replace this equals with an invalidation-only version, so that we can separate "needs to update the store" from "needs to update k8s"
		if existing.NeedsUpdate(*desired) {
			nodes = append(nodes, &actions.ActionNode{
				Action: actions.Action{
					Name: "Delete StatefulSet " + id,
					Do: func(ctx context.Context) error {
						return c.client.DeleteStatefulSet(ctx, existing)
					},
				},
				Next: []*actions.ActionNode{
					{
						Action: actions.Action{
							Name: "Recreate StatefulSet " + id,
							Do: func(ctx context.Context) error {
								return c.client.CreateStatefulSet(ctx, desired)
							},
						},
					},
				},
			})
		}
	}

	return nodes
}

func (c *Controller) reconcileOrphanedStatefulSets() []*actions.ActionNode {
	var nodes []*actions.ActionNode

	for id, statefulSet := range c.statefulSets.Map() {
		_, instanceExists := c.instances.Get(id)

		if !instanceExists {
			nodes = append(nodes, &actions.ActionNode{
				Action: actions.Action{
					Name: "Delete StatefulSet " + id,
					Do: func(ctx context.Context) error {
						return c.client.DeleteStatefulSet(ctx, statefulSet)
					},
				},
			})
		}
	}

	return nodes
}
