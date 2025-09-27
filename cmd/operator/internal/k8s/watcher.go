package k8s

import (
	"log/slog"

	"k8s.io/client-go/tools/cache"
)

type Watcheable[T any] interface {
	Equal(old T) bool
}

// enums
type EventType string

const (
	Added   EventType = "ADDED"
	Updated EventType = "UPDATED"
	Deleted EventType = "DELETED"
)

type WatchEvent[T any] struct {
	Type EventType
	Old  T
	New  T
}

type Watcher[T Watcheable[T]] struct {
	events  chan WatchEvent[*T]
	done    chan struct{}
	convert func(any) T
}

func NewWatcher[T Watcheable[T]](convert func(any) T) *Watcher[T] {
	return &Watcher[T]{
		convert: convert,
		done:    make(chan struct{}),
		events:  make(chan WatchEvent[*T], 100),
	}
}

func (w *Watcher[T]) Added(obj any) {
	converted := w.convert(obj)
	slog.Debug("New object event", getLogMeta(obj)...)
	w.events <- WatchEvent[*T]{Type: Added, New: &converted}
}

func (w *Watcher[T]) Updated(oldObj, newObj any) {
	convertedOld := w.convert(oldObj)
	convertedNew := w.convert(newObj)
	if convertedOld.Equal(convertedNew) {
		slog.Debug("No changes detected", getLogMeta(newObj)...)
		return
	}
	slog.Debug("Updated object event", getLogMeta(newObj)...)
	w.events <- WatchEvent[*T]{Type: Updated, Old: &convertedOld, New: &convertedNew}
}

func (w *Watcher[T]) Deleted(obj any) {
	slog.Debug("Deleted object event", getLogMeta(obj)...)
	converted := w.convert(obj)
	w.events <- WatchEvent[*T]{Type: Deleted, Old: &converted}
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

func (w *Watcher[T]) Events() <-chan WatchEvent[*T] {
	return w.events
}
