package seeds

import (
	"github.com/securemesh/coding/state"
)

var chatState = newStateFromSeed([][]byte{
	/* 01 */ []byte("',.0:?CIbgjkpvxz\xea"),
	/* 02 */ []byte("\nfw"),
	/* 03 */ []byte("cdmuy"),
	/* 04 */ []byte("l"),
	/* 05 */ []byte("r"),
	/* 06 */ []byte("t"),
	/* 07 */ []byte("ahos"),
	/* 08 */ []byte("in"),
	/* 09 */ []byte(""),
	/* 10 */ []byte(""),
	/* 11 */ []byte(" "),
	/* 12 */ []byte(""),
	/* 13 */ []byte("e"),
})

func ChatState() *state.State {
	return chatState.Clone()
}

func newStateFromSeed(seed [][]byte) *state.State {
	st := state.NewState()

	for i := range seed {
		for _, s := range seed[i:] {
			for _, b := range s {
				st.IncrementSymbol(b)
			}
		}
	}

	return st
}
