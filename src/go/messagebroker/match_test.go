package messagebroker

import (
	"strconv"
	"testing"
)

func TestMatchEvent_String(t *testing.T) {
	tests := []struct {
		name string
		e    MatchEvent
		want string
	}{
		{
			name: "test_UNWINDIA_MATCH_NEW",
			e:    UNWINDIA_MATCH_NEW,
			want: "UNWINDIA_MATCH_NEW",
		},
		{
			name: "test_max_eventid",
			e:    _max_eventid,
			want: strconv.Itoa(int(_max_eventid)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchEvent_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		e       MatchEvent
		args    args
		wantErr bool
	}{
		{
			name:    "test_ok",
			e:       MatchEvent(0),
			args:    struct{ b []byte }{b: []byte("UNWINDIA_MATCH_NEW")},
			wantErr: false,
		},
		{
			name:    "test_unknown",
			e:       MatchEvent(0),
			args:    struct{ b []byte }{b: []byte("UNKNOWN_MATCH_EVENT_WHICH_SHOULD_FAIL")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
