package k8s

import (
	"os"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	client        *dynamic.DynamicClient
	clientset     *kubernetes.Clientset
	labelSelector string
	labels        map[string]string
}

func NewClient(labels map[string]string) (*Client, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	var config *rest.Config
	var err error

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	labelSelector := ""
	for key, value := range labels {
		if labelSelector != "" {
			labelSelector += ","
		}
		labelSelector += key + "=" + value
	}

	return &Client{
		client:        dynClient,
		clientset:     clientset,
		labelSelector: labelSelector,
		labels:        labels,
	}, nil
}
