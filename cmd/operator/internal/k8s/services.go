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

func (c *Client) CreateService(ctx context.Context, svc *model.Service) error {
	_, err := c.clientset.CoreV1().Services(svc.Namespace).Create(ctx, &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Labels:    c.labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": svc.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Port: 1234,
				},
			},
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}, v1.CreateOptions{})
	return err
}

func (c *Client) DeleteService(ctx context.Context, svc *model.Service) error {
	return c.clientset.CoreV1().Services(svc.Namespace).Delete(ctx, svc.Name, v1.DeleteOptions{})
}
