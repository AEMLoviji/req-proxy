package observer

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

type RequestTracker interface {
	AddEntry(te Entry)
	ListEntries() map[uuid.UUID]Entry
}

type Entry struct {
	ClientRequest      interface{}
	ThirdPartyResponse interface{}
}

type ProxyRequestTracker struct {
	store map[uuid.UUID]Entry
	mu    sync.Mutex
}

func NewProxyRequestTracker() *ProxyRequestTracker {
	return &ProxyRequestTracker{
		store: make(map[uuid.UUID]Entry),
	}
}

func (t *ProxyRequestTracker) AddEntry(te Entry) {
	uuid := uuid.New()
	log.Printf("request is being tracked with id %s", uuid)

	t.mu.Lock()
	t.store[uuid] = te
	t.mu.Unlock()
}

func (t *ProxyRequestTracker) ListEntries() map[uuid.UUID]Entry {
	return t.store
}
