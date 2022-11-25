// Package workitemLock provides a locking mechanism for workitems
package workitemLock

import (
	"context"
	"time"
)

const (
	defaultTTL = time.Second * 120
)

// WorkItemLock allows the unwindia services to lock workitems to prevent concurrent processing
type WorkItemLock interface {
	Lock(ctx context.Context, workitemID string, ttl *time.Duration) error
	Unlock(ctx context.Context, workitemID string) error
}
