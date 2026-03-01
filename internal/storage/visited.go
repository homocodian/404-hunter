package storage

import "sync"

type VisitedLink map[string]struct{}

type VisitedStore interface {
	Add(string) bool
}

type InMemoryVisited struct {
	mu   sync.Mutex
	data VisitedLink
}

func NewInMemoryVisited() *InMemoryVisited {
	return &InMemoryVisited{
		data: make(VisitedLink),
	}
}

func (vStore *InMemoryVisited) Add(link string) bool {
	vStore.mu.Lock()
	defer vStore.mu.Unlock()

	if _, ok := vStore.data[link]; ok {
		return false
	}
	vStore.data[link] = struct{}{}
	return true
}
