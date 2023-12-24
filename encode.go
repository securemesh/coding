package coding

import (
	"bytes"

	"github.com/icza/bitio"
	"github.com/samber/lo"
	"github.com/securemesh/coding/heap"
)

func Encode(h *heap.Heap, msg []byte) []byte {
	buf := &bytes.Buffer{}
	w := bitio.NewWriter(buf)

	for _, b := range msg {
		index := h.IncrementSymbol(b)
		code := codes[index]
		lo.Must0(w.WriteBits(uint64(code.value), uint8(code.bits)))
	}

	lo.Must0(w.Close())
	return buf.Bytes()
}
