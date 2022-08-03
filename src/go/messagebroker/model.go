package messagebroker

import (
	"fmt"
	"strconv"
)

type Events int

const (
	UNWINDIA_CMS_CONTEST_NEW Events = iota
	UNWINDIA_CMS_CONTEST_READY_A
	UNWINDIA_CMS_CONTEST_READY_B
	UNWINDIA_CMS_CONTEST_READY_ALL
	_max_eventid
)

var Events_name = map[int]string{
	0: "UNWINDIA_CMS_CONTEST_NEW",
	1: "UNWINDIA_CMS_CONTEST_READY_A",
	2: "UNWINDIA_CMS_CONTEST_READY_B",
	3: "UNWINDIA_CMS_CONTEST_READY_ALL",
}

var Events_value = map[string]Events{
	"UNWINDIA_CMS_CONTEST_NEW":       UNWINDIA_CMS_CONTEST_NEW,
	"UNWINDIA_CMS_CONTEST_READY_A":   UNWINDIA_CMS_CONTEST_READY_A,
	"UNWINDIA_CMS_CONTEST_READY_B":   UNWINDIA_CMS_CONTEST_READY_B,
	"UNWINDIA_CMS_CONTEST_READY_ALL": UNWINDIA_CMS_CONTEST_READY_ALL,
}

func (e Events) String() string {
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

// UnmarshalJSON unmarshals b into Events.
func (e *Events) UnmarshalJSON(b []byte) error {
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

		*e = Events(ci)
		return nil
	}

	s := string(b)
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}

	if ev, ok := Events_value[s]; ok {
		*e = ev
		return nil
	}
	return fmt.Errorf("invalid code: %q", string(b))
}

type MessageType struct {
	Created string
	Updated string
	Deleted string
}

var MessageTypes = MessageType{
	Created: "created",
	Updated: "updated",
	Deleted: "deleted",
}

type Message struct {
	Specversion    string       `json:"specversion"`
	Type           string       `json:"type"`
	ModificationID string       `json:"modification_id"`
	Data           *interface{} `json:"data,omitempty"`
}
