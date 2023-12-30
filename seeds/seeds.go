package seeds

import (
	"github.com/securemesh/coding/heap"
)

var chatHeap = newHeapFromSeed([][]byte{
	/* 01 */ []byte("\x07'(,-8?ACDFHJLMNPRSTUWYbcfgjkpxzê"),
	/* 02 */ []byte("\n.dvw"),
	/* 03 */ []byte("Ihlmor"),
	/* 04 */ []byte("nu"),
	/* 05 */ []byte("ey"),
	/* 06 */ []byte("i"),
	/* 07 */ []byte("s"),
	/* 08 */ []byte(""),
	/* 09 */ []byte(""),
	/* 10 */ []byte(""),
	/* 11 */ []byte("at"),
	/* 12 */ []byte(""),
	/* 13 */ []byte(""),
	/* 14 */ []byte(""),
	/* 15 */ []byte(" "),
})

func ChatHeap() *heap.Heap {
	return chatHeap.Clone()
}

func newHeapFromSeed(seed [][]byte) *heap.Heap {
	h := heap.NewHeap()

	for i := range seed {
		for _, s := range seed[i:] {
			for _, b := range s {
				h.IncrementSymbol(b)
			}
		}
	}

	return h
}
