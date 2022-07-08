package proxy_test

import (
	"sync"
	"testing"

	"req-proxy/proxy"

	"github.com/stretchr/testify/assert"
)

func TestTrackAddEntryToTrackList(t *testing.T) {

	proxy.Track(proxy.TrackEntry{
		ClientRequest:      struct{}{},
		ThirdPartyResponse: struct{}{},
	})

	assert.EqualValues(t, 1, len(proxy.TrackList()))
}

func TestTrackConcurentlyAddEntryToTrackList(t *testing.T) {
	trackerFunc := func(wg *sync.WaitGroup) {
		proxy.Track(proxy.TrackEntry{
			ClientRequest:      struct{}{},
			ThirdPartyResponse: struct{}{},
		})
		wg.Done()
	}

	const ConcurrentRunCount int = 100

	var waitG sync.WaitGroup
	waitG.Add(ConcurrentRunCount)
	for i := 0; i < ConcurrentRunCount; i++ {
		go trackerFunc(&waitG)
	}
	waitG.Wait()

	assert.EqualValues(t, ConcurrentRunCount, len(proxy.TrackList()))
}
