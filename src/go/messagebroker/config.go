package messagebroker

import (
	"fmt"
	"strconv"
)

type PulsarAuth int

const (
	AUTH_SIMPLE PulsarAuth = iota
	AUTH_OAUTH2
	_max_authid
)

var PulsarAuth_name = map[int]string{
	0: "simple",
	1: "oauth2",
}

var PulsarAuth_value = map[string]PulsarAuth{
	PulsarAuth_name[0]: AUTH_SIMPLE,
	PulsarAuth_name[1]: AUTH_OAUTH2,
}

func (p PulsarAuth) String() string {
	s, ok := PulsarAuth_name[int(p)]
	if ok {
		return s
	}
	return strconv.Itoa(int(p))
}

// UnmarshalJSON unmarshals b into PulsarAuth.
func (p *PulsarAuth) UnmarshalJSON(b []byte) error {
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
		if ci >= uint64(_max_authid) {
			return fmt.Errorf("invalid code: %q", ci)
		}

		*p = PulsarAuth(ci)
		return nil
	}

	if ev, ok := PulsarAuth_value[string(b)]; ok {
		*p = ev
		return nil
	}
	return fmt.Errorf("invalid code: %q", string(b))
}

func (p *PulsarAuth) Unmarshal(data string) error {
	return p.UnmarshalJSON([]byte(data))
}

func (p *PulsarAuth) UnmarshalText(text []byte) error {
	return p.UnmarshalJSON(text)
}
