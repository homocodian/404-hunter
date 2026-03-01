package storage

import "sync"

type BrokenURL map[string]struct{}

type BrokenStore interface {
	Add(string) bool
}

type InMemoryBroken struct {
	mu   sync.Mutex
	data BrokenURL
}

func NewInMemoryBroken() *InMemoryBroken {
	return &InMemoryBroken{
		data: make(BrokenURL),
	}
}

func (s *InMemoryBroken) Add(url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[url]; ok {
		return false
	}

	s.data[url] = struct{}{}
	return true
}

func (s *InMemoryBroken) GetAll() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var urls []string
	for url := range s.data {
		urls = append(urls, url)
	}

	return urls
}
