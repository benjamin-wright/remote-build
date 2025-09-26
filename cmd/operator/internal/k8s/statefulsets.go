package k8s

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"pongle-hub.co.uk/remote-build/cmd/operator/internal/model"
)

func anyToStatefulSet(obj any) model.StatefulSet {
	svc := obj.(*appsv1.StatefulSet)

	// Convert the obj to model.StatefulSet
	return model.StatefulSet{
		Name:      svc.Name,
		Namespace: svc.Namespace,
		CPU:       svc.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String(),
		Memory:    svc.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String(),
		Disk:      svc.Spec.Template.Spec.Containers[0].Resources.Requests.Storage().String(),
		Image:     svc.Spec.Template.Spec.Containers[0].Image,
		Ready:     svc.Status.ReadyReplicas == 1,
	}
}

func (c *Client) WatchStatefulSets(ctx context.Context) *Watcher[model.StatefulSet] {
	watcher := NewWatcher(anyToStatefulSet)

	go func() {
		defer close(watcher.done)
		informer := informers.NewSharedInformerFactoryWithOptions(
			c.clientset,
			time.Second*30,
			informers.WithTweakListOptions(func(options *v1.ListOptions) {
				options.LabelSelector = c.labelSelector
			}),
		).Apps().V1().StatefulSets().Informer()

		informer.AddEventHandler(watcher.GetEventHandler())
		informer.Run(ctx.Done())
	}()

	return watcher
}
