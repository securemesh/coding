package heap

import (
	"maps"
)

type node struct {
	symbol byte
	count  int
}

type Heap struct {
	nodes    [256]node
	bySymbol map[byte]int
}

func NewHeap() *Heap {
	h := &Heap{
		bySymbol: map[byte]int{},
	}

	for i := 0; i < 256; i++ {
		h.nodes[i].symbol = byte(i)
		h.bySymbol[byte(i)] = i
	}

	return h
}

func (h Heap) Clone() *Heap {
	return &Heap{
		nodes:    h.nodes,
		bySymbol: maps.Clone(h.bySymbol),
	}
}

func (h *Heap) IncrementSymbol(symbol byte) int {
	nodeIndex := h.bySymbol[symbol]
	h.nodes[nodeIndex].count++

	iterIndex := nodeIndex
	for iterIndex != 0 {
		parentIndex := h.parentIndex(iterIndex)

		if h.nodes[iterIndex].count <= h.nodes[parentIndex].count {
			break
		}

		h.nodes[iterIndex], h.nodes[parentIndex] = h.nodes[parentIndex], h.nodes[iterIndex]
		h.bySymbol[h.nodes[iterIndex].symbol] = iterIndex
		h.bySymbol[h.nodes[parentIndex].symbol] = parentIndex
		iterIndex = parentIndex
	}

	return nodeIndex
}

func (h Heap) parentIndex(nodeIndex int) int {
	return (nodeIndex - 1) / 2
}
