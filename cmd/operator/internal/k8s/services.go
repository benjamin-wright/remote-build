package k8s

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"pongle-hub.co.uk/remote-build/cmd/operator/internal/model"
)

type ServiceWatcher struct {
	Added   func(service model.Service)
	Updated func(old model.Service, new model.Service)
	Deleted func(service model.Service)
}

func (c *Client) WatchServices(ctx context.Context, watcher ServiceWatcher) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)
		informer := informers.NewSharedInformerFactoryWithOptions(
			c.clientset,
			time.Second*30,
			informers.WithTweakListOptions(func(options *v1.ListOptions) {
				options.LabelSelector = c.labelSelector
			}),
		).Core().V1().Services().Informer()

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(svc any) {
				new := svc.(*corev1.Service)

				watcher.Added(model.Service{
					Name:      new.Name,
					Namespace: new.Namespace,
				})
			},
			UpdateFunc: func(svc1 any, svc2 any) {
				old := svc1.(*corev1.Service)
				new := svc2.(*corev1.Service)

				watcher.Updated(model.Service{
					Name:      new.Name,
					Namespace: new.Namespace,
				}, model.Service{
					Name:      old.Name,
					Namespace: old.Namespace,
				})
			},
			DeleteFunc: func(svc any) {
				cast := svc.(*corev1.Service)
				watcher.Deleted(model.Service{
					Name:      cast.Name,
					Namespace: cast.Namespace,
				})
			},
		})
	}()

	return errCh
}
