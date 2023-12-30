package heap

import (
	"fmt"
	"maps"
	"slices"
	"strings"
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

		if h.nodes[iterIndex].count < h.nodes[parentIndex].count || (h.nodes[iterIndex].count == h.nodes[parentIndex].count && h.nodes[iterIndex].symbol > h.nodes[parentIndex].symbol) {
			break
		}

		h.nodes[iterIndex], h.nodes[parentIndex] = h.nodes[parentIndex], h.nodes[iterIndex]
		h.bySymbol[h.nodes[iterIndex].symbol] = iterIndex
		h.bySymbol[h.nodes[parentIndex].symbol] = parentIndex
		iterIndex = parentIndex
	}

	return nodeIndex
}

func (h Heap) String() string {
	nodes := []node{}

	for _, node := range h.nodes {
		if node.count == 0 {
			continue
		}

		nodes = append(nodes, node)
	}

	slices.SortStableFunc(nodes, func(a, b node) int { return int(a.symbol) - int(b.symbol) })
	slices.SortStableFunc(nodes, func(a, b node) int { return a.count - b.count })

	strs := []string{}

	for _, node := range nodes {
		strs = append(strs, fmt.Sprintf("{%#U}=%d", node.symbol, node.count))
	}

	return strings.Join(strs, ", ")
}

func (h Heap) parentIndex(nodeIndex int) int {
	return (nodeIndex - 1) / 2
}
