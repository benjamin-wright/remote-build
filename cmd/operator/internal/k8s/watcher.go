package k8s

import (
	"log/slog"

	"k8s.io/client-go/tools/cache"
)

type Watcheable[T any] interface {
	Equal(old T) bool
}

type Watcher[T Watcheable[T]] struct {
	objects map[string]T
	done    chan struct{}
	convert func(any) T
}

func NewWatcher[T Watcheable[T]](convert func(any) T) *Watcher[T] {
	return &Watcher[T]{
		objects: make(map[string]T),
		convert: convert,
		done:    make(chan struct{}),
	}
}

func (w *Watcher[T]) Added(obj any) {
	converted := w.convert(obj)
	slog.Debug("Adding new object", "object", getID(obj))
	w.objects[getID(obj)] = converted
}

func (w *Watcher[T]) Updated(oldObj, newObj any) {
	convertedOld := w.convert(oldObj)
	convertedNew := w.convert(newObj)
	if convertedOld.Equal(convertedNew) {
		slog.Debug("No changes detected", "object", getID(newObj))
		return
	}
	slog.Debug("Changes detected", "object", getID(newObj))
	w.objects[getID(newObj)] = convertedNew
}

func (w *Watcher[T]) Deleted(obj any) {
	slog.Debug("Deleting object", "object", getID(obj))
	delete(w.objects, getID(obj))
}

func (w *Watcher[T]) GetEventHandler() cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		AddFunc:    w.Added,
		UpdateFunc: w.Updated,
		DeleteFunc: w.Deleted,
	}
}

func (w *Watcher[T]) Done() <-chan struct{} {
	return w.done
}
