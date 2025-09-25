package k8s

import "pongle-hub.co.uk/remote-build/internal/model"

type StatefulSetWatcher interface {
	Added(set model.StatefulSet)
	Updated(old model.StatefulSet, new model.StatefulSet)
	Deleted(set model.StatefulSet)
}

func (c *Client) WatchStatefulSets(watcher StatefulSetWatcher) error {
	// Implementation to watch for changes in stateful sets
	return nil
}
