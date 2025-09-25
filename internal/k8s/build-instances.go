package k8s

import "pongle-hub.co.uk/remote-build/internal/model"

type BuildInstanceWatcher interface {
	Added(instance model.BuildInstance)
	Updated(old model.BuildInstance, new model.BuildInstance)
	Deleted(instance model.BuildInstance)
}

func (c *Client) WatchBuildInstances(watcher BuildInstanceWatcher) error {
	// Implementation to watch for changes in build instances
	return nil
}
