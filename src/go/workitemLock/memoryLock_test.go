package workitemLock

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestMemoryWorkItemLock_Lock(t *testing.T) {
	type fields struct {
		locks map[string]time.Time
	}
	type args struct {
		in0        context.Context
		workitemID string
		ttl        *time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test_ok",
			fields: fields{
				locks: make(map[string]time.Time),
			},
			args: args{
				in0:        context.Background(),
				workitemID: "t1",
				ttl:        nil,
			},
			wantErr: false,
		},
		{
			name: "test_item_already_locked",
			fields: fields{
				locks: map[string]time.Time{
					"t2": time.Now().Add(time.Second * 10),
				},
			},
			args: args{
				in0:        context.Background(),
				workitemID: "t2",
				ttl:        nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &MemoryWorkItemLock{
				locks: tt.fields.locks,
			}
			if err := w.Lock(tt.args.in0, tt.args.workitemID, tt.args.ttl); (err != nil) != tt.wantErr {
				t.Errorf("Lock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryWorkItemLock_StartHousekeeping(t *testing.T) {
	type fields struct {
		locks map[string]time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &MemoryWorkItemLock{
				locks: tt.fields.locks,
			}
			w.StartHousekeeping()
		})
	}
}

func TestMemoryWorkItemLock_Unlock(t *testing.T) {
	type fields struct {
		locks map[string]time.Time
	}
	type args struct {
		in0        context.Context
		workitemID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test_unnlock_ok",
			fields: fields{
				locks: map[string]time.Time{
					"t1": time.Now().Add(time.Second * 10),
				},
			},
			args: struct {
				in0        context.Context
				workitemID string
			}{
				in0:        context.Background(),
				workitemID: "t1",
			},
			wantErr: false,
		},
		{
			name: "test_unnlock_not_existing",
			fields: fields{
				locks: map[string]time.Time{
					"anyOtherWorkitem": time.Now().Add(time.Second * 10),
				},
			},
			args: struct {
				in0        context.Context
				workitemID string
			}{
				in0:        context.Background(),
				workitemID: "t2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &MemoryWorkItemLock{
				locks: tt.fields.locks,
			}
			if err := w.Unlock(tt.args.in0, tt.args.workitemID); (err != nil) != tt.wantErr {
				t.Errorf("Unlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryWorkItemLock_housekeeping(t *testing.T) {
	locks1 := map[string]time.Time{
		"t1": time.Now().Add(time.Second * 10),
	}

	type fields struct {
		locksBefore map[string]time.Time
		locksAfter  map[string]time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test_no_items_to_housekeep_yet",
			fields: fields{
				locksBefore: locks1,
				locksAfter:  locks1,
			},
		},
		{
			name: "test_housekeep_one_item",
			fields: fields{
				locksBefore: map[string]time.Time{
					"t1": time.Now().Add(time.Second * -10),
				},
				locksAfter: map[string]time.Time{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &MemoryWorkItemLock{
				locks: tt.fields.locksBefore,
			}
			w.housekeeping()
			if !reflect.DeepEqual(w.locks, tt.fields.locksAfter) {
				t.Errorf("housekeeping() = %v, want %v", w.locks, tt.fields.locksAfter)
			}
		})
	}
}

func TestNewMemoryWorkItemLock(t *testing.T) {
	tests := []struct {
		name string
		want *MemoryWorkItemLock
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemoryWorkItemLock(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemoryWorkItemLock() = %v, want %v", got, tt.want)
			}
		})
	}
}
