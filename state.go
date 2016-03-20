package pool

import "fmt"

type State uint32

// States used for the event loop
const (
	RUNNING State = 1 << iota
	FINISHED
	CHAN_HAS_SUCCESS
	CHAN_HAS_FAILURE
)

// String implementation for for the state
func (s State) String() string {
	var ret string
	is_state := func(c, g State) bool { return c == (c & g) }
	states := map[State]string{
		0:                "NIL",
		RUNNING:          "RUNNING",
		FINISHED:         "FINISHED",
		CHAN_HAS_SUCCESS: "CHAN_HAS_SUCCESS",
		CHAN_HAS_FAILURE: "CHAN_HAS_FAILURE",
	}

	for i, n := range states {
		if is_state(i, s) {
			ret += "|" + n
		}
	}

	return fmt.Sprintf("state: %s (0x%04x)", ret[1:], (uint32)(s))
}
