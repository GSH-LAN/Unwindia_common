package messagebroker

import (
	"fmt"
	"strconv"
)

type MessageTypes int

const (
	MessageTypeCreated MessageTypes = iota
	MessageTypeUpdated
	MessageTypeDeleted
	_max_messagetype
)

var MessageTypes_name = map[int]string{
	0: "create",
	1: "update",
	2: "delete",
}

var MessageTypes_value = map[string]MessageTypes{
	MessageTypes_name[0]: MessageTypeCreated,
	MessageTypes_name[1]: MessageTypeUpdated,
	MessageTypes_name[2]: MessageTypeDeleted,
}

func (m MessageTypes) String() string {
	s, ok := MessageTypes_name[int(m)]
	if ok {
		return s
	}
	return strconv.Itoa(int(m))
}

// UnmarshalJSON unmarshals b into MessageTypes.
func (m *MessageTypes) UnmarshalJSON(b []byte) error {
	// From json.Unmarshaler: By convention, to approximate the behavior of
	// Unmarshal itself, Unmarshalers implement UnmarshalJSON([]byte("null")) as
	// a no-op.
	if string(b) == "null" {
		return nil
	}
	if m == nil {
		return fmt.Errorf("nil receiver passed to UnmarshalJSON")
	}

	if ci, err := strconv.ParseUint(string(b), 10, 32); err == nil {
		if ci >= uint64(_max_messagetype) {
			return fmt.Errorf("invalid code: %q", ci)
		}

		*m = MessageTypes(ci)
		return nil
	}

	if mv, ok := MessageTypes_value[string(b)]; ok {
		*m = mv
		return nil
	}
	return fmt.Errorf("invalid code: %q", string(b))
}

type Message struct {
	Type    MessageTypes `json:"type"`
	SubType string       `json:"subtype"`
	Data    interface{}  `json:"data,omitempty"`
}
