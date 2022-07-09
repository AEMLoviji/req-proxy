package observer_test

import (
	"req-proxy/observer"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrackAddEntryToTrackList(t *testing.T) {
	requestTracker := observer.NewProxyRequestTracker()

	requestTracker.AddEntry(observer.Entry{
		ClientRequest:      struct{}{},
		ThirdPartyResponse: struct{}{},
	})

	assert.EqualValues(t, 1, len(requestTracker.ListEntries()))
}

func TestTrackConcurentlyAddEntryToTrackList(t *testing.T) {
	requestTracker := observer.NewProxyRequestTracker()

	trackerFunc := func(wg *sync.WaitGroup) {
		requestTracker.AddEntry(observer.Entry{
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

	assert.EqualValues(t, ConcurrentRunCount, len(requestTracker.ListEntries()))
}
