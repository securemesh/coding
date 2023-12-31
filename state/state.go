package state

import (
	"fmt"
	"maps"
	"strings"
)

type node struct {
	symbol byte
	count  int
}

type State struct {
	nodes    []*node
	bySymbol map[byte]int
}

func NewState() *State {
	st := &State{
		bySymbol: map[byte]int{},
	}

	for i := 0; i < 256; i++ {
		st.nodes = append(st.nodes, &node{
			symbol: byte(i),
		})

		st.bySymbol[byte(i)] = i
	}

	return st
}

func (st State) Clone() *State {
	st2 := &State{
		bySymbol: maps.Clone(st.bySymbol),
	}

	for _, node := range st.nodes {
		tmp := *node
		st2.nodes = append(st2.nodes, &tmp)
	}

	return st2
}

// Returns old index
func (st *State) IncrementSymbol(symbol byte) int {
	nodeIndex := st.bySymbol[symbol]
	st.nodes[nodeIndex].count++

	for iterIndex := nodeIndex; iterIndex > 0; iterIndex-- {
		prevIndex := iterIndex - 1
		iterNode := st.nodes[iterIndex]
		prevNode := st.nodes[prevIndex]

		if prevNode.count > iterNode.count {
			break
		} else if prevNode.count == iterNode.count && prevNode.symbol < iterNode.symbol {
			break
		}

		st.nodes[iterIndex] = prevNode
		st.bySymbol[prevNode.symbol] = iterIndex

		st.nodes[prevIndex] = iterNode
		st.bySymbol[iterNode.symbol] = prevIndex
	}

	return nodeIndex
}

func (st State) String() string {
	strs := []string{}

	for _, node := range st.nodes {
		if node.count == 0 {
			break
		}

		strs = append(strs, fmt.Sprintf("{%#U}=%d", node.symbol, node.count))
	}

	return strings.Join(strs, ", ")
}
