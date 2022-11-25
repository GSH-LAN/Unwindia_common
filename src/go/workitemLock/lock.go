package workitemLock

import (
	"context"
	"time"
)

const (
	defaultTtl = time.Second * 120
)

// WorkItemLock allows the unwindia services to lock workitems to prevent concurrent processing
type WorkItemLock interface {
	Lock(ctx context.Context, workitemID string, ttl *time.Duration) error
	Unlock(ctx context.Context, workitemID string) error
}
