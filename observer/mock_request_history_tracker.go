package observer

import (
	"github.com/google/uuid"
)

type MockRequestHistoryTracker struct {
	AddEntryFunc    func(te Entry)
	ListEntriesFunc func() map[uuid.UUID]Entry
}

func (m *MockRequestHistoryTracker) AddEntry(te Entry) {
	m.AddEntryFunc(te)
}

func (m *MockRequestHistoryTracker) ListEntries() map[uuid.UUID]Entry {
	return m.ListEntriesFunc()
}
