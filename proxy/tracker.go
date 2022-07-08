package proxy

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

type TrackEntry struct {
	ClientRequest      interface{}
	ThirdPartyResponse interface{}
}

var (
	trackList map[uuid.UUID]TrackEntry = make(map[uuid.UUID]TrackEntry)
	mu        sync.Mutex
)

func Track(te TrackEntry) {
	uuid := uuid.New()
	log.Printf("request is being tracked with id %s", uuid)

	mu.Lock()
	trackList[uuid] = te
	mu.Unlock()
}

func TrackList() map[uuid.UUID]TrackEntry {
	return trackList
}
