package workitemLock

import (
	"fmt"
	"strconv"
	"time"
)

const (
	housekeepingInterval = time.Minute
)

type WorkItemLockEntry struct {
	ID        string     `json:"id,omitempty" bson:"_id,omitempty"`
	LockedBy  string     `json:"lockedBy" bson:"lockedBy"`
	CreatedAt time.Time  `json:"createdAt" bson:"createdAt"`
	ExpiresAt *time.Time `json:"expiresAt" bson:"expiresAt"`
}

type WorkItemLockType int

const (
	LOCK_MEMORY WorkItemLockType = iota
	LOCK_MONGODB
	_maxLock
)

var LockTypeName = map[int]string{
	0: "memory",
	1: "mongodb",
}

var LockTypeValue = map[string]WorkItemLockType{
	LockTypeName[0]: LOCK_MEMORY,
	LockTypeName[1]: LOCK_MONGODB,
}

func (p WorkItemLockType) String() string {
	s, ok := LockTypeName[int(p)]
	if ok {
		return s
	}
	return strconv.Itoa(int(p))
}

// UnmarshalJSON unmarshals b into PulsarAuth.
func (p *WorkItemLockType) UnmarshalJSON(b []byte) error {
	// From json.Unmarshaler: By convention, to approximate the behavior of
	// Unmarshal itself, Unmarshalers implement UnmarshalJSON([]byte("null")) as
	// a no-op.
	if string(b) == "null" {
		return nil
	}
	if p == nil {
		return fmt.Errorf("nil receiver passed to UnmarshalJSON")
	}

	if ci, err := strconv.ParseUint(string(b), 10, 32); err == nil {
		if ci >= uint64(_maxLock) {
			return fmt.Errorf("invalid code: %q", ci)
		}

		*p = WorkItemLockType(ci)
		return nil
	}

	s := string(b)
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}

	if ev, ok := LockTypeValue[s]; ok {
		*p = ev
		return nil
	}
	return fmt.Errorf("invalid code: %q", string(b))
}

func (p *WorkItemLockType) Unmarshal(data string) error {
	return p.UnmarshalJSON([]byte(data))
}

func (p *WorkItemLockType) UnmarshalText(text []byte) error {
	return p.UnmarshalJSON(text)
}
