package k8s

import "pongle-hub.co.uk/remote-build/cmd/operator/internal/model"

type ServiceWatcher interface {
	Added(service model.Service)
	Updated(old model.Service, new model.Service)
	Deleted(service model.Service)
}

func (c *Client) WatchServices(watcher ServiceWatcher) error {
	// Implementation to watch for changes in services
	return nil
}
