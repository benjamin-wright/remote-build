package k8s

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"pongle-hub.co.uk/remote-build/cmd/operator/internal/model"
)

func anyToService(obj any) model.Service {
	svc := obj.(*corev1.Service)

	// Convert the obj to model.Service
	return model.Service{
		Name:      svc.Name,
		Namespace: svc.Namespace,
	}
}

func (c *Client) WatchServices(ctx context.Context) *Watcher[model.Service] {
	watcher := NewWatcher(anyToService)

	go func() {
		defer close(watcher.done)
		informer := informers.NewSharedInformerFactoryWithOptions(
			c.clientset,
			time.Second*30,
			informers.WithTweakListOptions(func(options *v1.ListOptions) {
				options.LabelSelector = c.labelSelector
			}),
		).Core().V1().Services().Informer()

		informer.AddEventHandler(watcher.GetEventHandler())
		informer.Run(ctx.Done())
	}()

	return watcher
}
