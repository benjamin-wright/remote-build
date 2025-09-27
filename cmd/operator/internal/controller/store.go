package controller

import "pongle-hub.co.uk/remote-build/cmd/operator/internal/k8s"

type Identifiable interface {
	ID() string
}

type Store[T Identifiable] struct {
	objects map[string]T
}

func NewStore[T Identifiable]() Store[T] {
	return Store[T]{
		objects: make(map[string]T),
	}
}

func (s *Store[T]) add(bi T) {
	s.objects[bi.ID()] = bi
}

func (s *Store[T]) Get(id string) (T, bool) {
	bi, ok := s.objects[id]
	return bi, ok
}

func (s *Store[T]) Map() map[string]T {
	return s.objects
}

func (s *Store[T]) remove(id string) {
	delete(s.objects, id)
}

func (s *Store[T]) ProcessEvent(event k8s.WatchEvent[T]) {
	switch event.Type {
	case k8s.Added:
		s.add(event.New)
	case k8s.Updated:
		s.add(event.New)
	case k8s.Deleted:
		s.remove(event.Old.ID())
	}
}
