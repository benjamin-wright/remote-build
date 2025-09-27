package k8s

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
		Disk:      svc.Spec.VolumeClaimTemplates[0].Spec.Resources.Requests.Storage().String(),
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

func (c *Client) CreateStatefulSet(ctx context.Context, sts *model.StatefulSet) error {
	replicas := int32(1)
	_, err := c.clientset.AppsV1().StatefulSets(sts.Namespace).Create(ctx, &appsv1.StatefulSet{
		ObjectMeta: v1.ObjectMeta{
			Name:      sts.Name,
			Namespace: sts.Namespace,
			Labels:    c.labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{"app": sts.Name},
			},
			ServiceName: sts.Name,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"app": sts.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  sts.Name,
						Image: sts.Image,
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse(sts.CPU),
								corev1.ResourceMemory: resource.MustParse(sts.Memory),
							},
						},
					}},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: v1.ObjectMeta{
					Name: "data",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.PersistentVolumeAccessMode("ReadWriteOnce"),
					},
					Resources: corev1.VolumeResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: resource.MustParse(sts.Disk),
						},
					},
				},
			}},
		},
	}, v1.CreateOptions{})
	return err
}

func (c *Client) DeleteStatefulSet(ctx context.Context, sts *model.StatefulSet) error {
	return c.clientset.AppsV1().StatefulSets(sts.Namespace).Delete(ctx, sts.Name, v1.DeleteOptions{})
}
