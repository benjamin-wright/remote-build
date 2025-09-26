package k8s

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"pongle-hub.co.uk/remote-build/cmd/operator/internal/model"
)

const GroupName = "remote-build.pongle-hub.co.uk"
const GroupVersion = "v1alpha1"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

func (c *Client) WatchBuildInstances(ctx context.Context) *Watcher[model.BuildInstance] {
	watcher := NewWatcher(func(obj any) model.BuildInstance {
		instance := obj.(*unstructured.Unstructured)

		output := model.BuildInstance{
			Name:      instance.GetName(),
			Namespace: instance.GetNamespace(),
			CPU:       nestedString(instance.Object, "spec", "cpu"),
			Memory:    nestedString(instance.Object, "spec", "memory"),
			Disk:      nestedString(instance.Object, "spec", "disk"),
			Image:     nestedString(instance.Object, "spec", "image"),
			State:     nestedString(instance.Object, "status", "state"),
			Active:    nestedBool(instance.Object, "status", "active"),
		}

		return output
	})

	go func() {
		defer close(watcher.done)
		informer := dynamicinformer.NewFilteredDynamicSharedInformerFactory(c.client, time.Second*30, "", func(lo *v1.ListOptions) {
			lo.LabelSelector = c.labelSelector
		}).ForResource(SchemeGroupVersion.WithResource("buildinstances")).Informer()

		informer.AddEventHandler(watcher.GetEventHandler())

		informer.Run(ctx.Done())
	}()

	return watcher
}
