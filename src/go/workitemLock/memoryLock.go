package workitemLock

import (
	"context"
	"errors"
	"time"
)

type MemoryWorkItemLock struct {
	locks map[string]time.Time
}

func NewMemoryWorkItemLock() *MemoryWorkItemLock {
	lock := &MemoryWorkItemLock{
		locks: make(map[string]time.Time),
	}

	go lock.StartHousekeeping()

	return lock
}

func (w *MemoryWorkItemLock) housekeeping() {
	for workitemID, expiresAt := range w.locks {
		if time.Now().After(expiresAt) {
			delete(w.locks, workitemID)
		}
	}

}

func (w *MemoryWorkItemLock) Lock(_ context.Context, workitemID string, ttl *time.Duration) error {
	if lock, ok := w.locks[workitemID]; ok {
		if time.Now().Before(lock) {
			return errors.New("workitem is already locked")
		}
	}

	expiresAt := time.Now().Add(defaultTtl)
	if ttl != nil {
		expiresAt = time.Now().Add(*ttl)
	}

	w.locks[workitemID] = expiresAt
	return nil
}

func (w *MemoryWorkItemLock) Unlock(_ context.Context, workitemID string) error {
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
