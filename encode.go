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

	for i := 0; i < len(msg); {
		l, index := st.IncrementSymbol(msg[i:])
		i += l
		code := codes.CodeForIndex(uint32(index))
		lo.Must0(w.WriteBits(uint64(code.Value), uint8(code.Bits)))
	}

	lo.Must0(w.Close())
	return buf.Bytes()
}
