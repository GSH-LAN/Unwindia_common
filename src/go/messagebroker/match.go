package messagebroker

import (
	"fmt"
	"github.com/GSH-LAN/Unwindia_common/src/go/matchservice"
	"strconv"
)

const TOPIC = "UNWINDIA_MATCH"

type MatchMessage struct {
	Message
	SubType MatchEvent `json:"subtype"`
	Data    *matchservice.MatchInfo
}

type MatchEvent int

const (
	UNWINDIA_MATCH_NEW MatchEvent = iota
	UNWINDIA_MATCH_READY_A
	UNWINDIA_MATCH_READY_B
	UNWINDIA_MATCH_READY_ALL
	UNWINDIA_MATCH_UPDATE
	UNWINDIA_MATCH_FINISHED
	_max_eventid
)

var Events_name = map[int]string{
	0: "UNWINDIA_MATCH_NEW",
	1: "UNWINDIA_MATCH_READY_A",
	2: "UNWINDIA_MATCH_READY_B",
	3: "UNWINDIA_MATCH_READY_ALL",
	4: "UNWINDIA_MATCH_FINISHED",
}

var Events_value = map[string]MatchEvent{
	"UNWINDIA_MATCH_NEW":       UNWINDIA_MATCH_NEW,
	"UNWINDIA_MATCH_READY_A":   UNWINDIA_MATCH_READY_A,
	"UNWINDIA_MATCH_READY_B":   UNWINDIA_MATCH_READY_B,
	"UNWINDIA_MATCH_READY_ALL": UNWINDIA_MATCH_READY_ALL,
	"UNWINDIA_MATCH_FINISHED":  UNWINDIA_MATCH_FINISHED,
}

func (e MatchEvent) String() string {
	s, ok := Events_name[int(e)]
	if ok {
		return s
	}
	return strconv.Itoa(int(e))
}

type NewContest struct {
}
type Response struct {
	value interface{}
}

// UnmarshalJSON unmarshals b into MatchEvent.
func (e *MatchEvent) UnmarshalJSON(b []byte) error {
	// From json.Unmarshaler: By convention, to approximate the behavior of
	// Unmarshal itself, Unmarshalers implement UnmarshalJSON([]byte("null")) as
	// a no-op.
	if string(b) == "null" {
		return nil
	}
	if e == nil {
		return fmt.Errorf("nil receiver passed to UnmarshalJSON")
	}

	if ci, err := strconv.ParseUint(string(b), 10, 32); err == nil {
		if ci >= uint64(_max_eventid) {
			return fmt.Errorf("invalid code: %q", ci)
		}

		*e = MatchEvent(ci)
		return nil
	}

	if ev, ok := Events_value[string(b)]; ok {
		*e = ev
		return nil
	}
	return fmt.Errorf("invalid code: %q", string(b))
}