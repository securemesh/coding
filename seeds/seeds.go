package seeds

import (
	"github.com/securemesh/coding/state"
)

var chatState = newStateFromSeed([][]byte{
	/* 01 */ []byte("',.?CIbcfjkpvxz\xea"),
	/* 02 */ []byte("dgwy"),
	/* 03 */ []byte("\nmru"),
	/* 04 */ []byte("l"),
	/* 05 */ []byte("hns"),
	/* 06 */ []byte("a"),
	/* 07 */ []byte(""),
	/* 08 */ []byte("i"),
	/* 09 */ []byte(""),
	/* 10 */ []byte(" "),
	/* 11 */ []byte(""),
	/* 12 */ []byte(""),
	/* 13 */ []byte(""),
	/* 14 */ []byte("et"),
	/* 15 */ []byte(""),
	/* 16 */ []byte("o"),
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
