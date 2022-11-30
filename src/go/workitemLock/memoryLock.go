package workitemLock

import (
	"context"
	"errors"
	"sync"
	"time"
)

// MemoryWorkItemLock is a simple in-memory implementation of the WorkItemLock interface. Cannot be used for horizontal scaling.
type MemoryWorkItemLock struct {
	locks map[string]time.Time
	lock  sync.RWMutex
}

func NewMemoryWorkItemLock() *MemoryWorkItemLock {
	lock := &MemoryWorkItemLock{
		locks: make(map[string]time.Time),
	}

	go lock.StartHousekeeping()

	return lock
}

func (w *MemoryWorkItemLock) housekeeping() {
	w.lock.Lock()
	defer w.lock.Unlock()
	
	for workitemID, expiresAt := range w.locks {
		if time.Now().After(expiresAt) {
			delete(w.locks, workitemID)
		}
	}

}

func (w *MemoryWorkItemLock) Lock(_ context.Context, workitemID string, ttl *time.Duration) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if lock, ok := w.locks[workitemID]; ok {
		if time.Now().Before(lock) {
			return errors.New("workitem is already locked")
		}
	}

	expiresAt := time.Now().Add(defaultTTL)
	if ttl != nil {
		expiresAt = time.Now().Add(*ttl)
	}

	w.locks[workitemID] = expiresAt
	return nil
}

func (w *MemoryWorkItemLock) Unlock(_ context.Context, workitemID string) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	delete(w.locks, workitemID)
	return nil
}

func (w *MemoryWorkItemLock) StartHousekeeping() {
	ticker := time.NewTicker(housekeepingInterval)
	for {
		select {
		case <-ticker.C:
			w.housekeeping()
		}
	}
}
