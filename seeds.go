package coding

import (
	"github.com/securemesh/coding/heap"
)

var chatHeap = newHeapFromSeed([][]byte{
	[]byte(`]\_}`),
	[]byte(`[ê%=Z`),
	[]byte(`#ÄQ<>`),
	[]byte(`&X@+*`),
	[]byte(`$~"V;`),
	[]byte(`/78q9`),
	[]byte("zRE54F(U-6\n"),
	[]byte(`NLx:C01D2BJ)K3GP`),
	[]byte(`STWH!OYAjM`),
	[]byte(`?,'`),
	[]byte(`bIv`),
	[]byte(`mygwc.pfk`),
	[]byte(`isrhlud`),
	[]byte(`eotan`),
	[]byte(` `),
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
