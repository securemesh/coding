package coding

import (
	"bytes"

	"github.com/icza/bitio"
	"github.com/samber/lo"
	"github.com/securemesh/coding/codes"
	"github.com/securemesh/coding/state"
)

func Encode(st *state.State, msg []byte) []byte {
	buf := &bytes.Buffer{}
	w := bitio.NewWriter(buf)

	for _, b := range msg {
		index := st.IncrementSymbol(b)
		code := codes.CodeForIndex(index)
		lo.Must0(w.WriteBits(uint64(code.Value), uint8(code.Bits)))
	}

	lo.Must0(w.Close())
	return buf.Bytes()
}
