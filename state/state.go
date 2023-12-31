package state

import (
	"fmt"
	"strings"
)

type node struct {
	symbol byte
	count  int
	index  int
}

type State struct {
	nodes    []*node
	bySymbol map[byte]*node
}

func NewState() *State {
	st := &State{
		bySymbol: map[byte]*node{},
	}

	for i := 0; i < 256; i++ {
		node := &node{
			symbol: byte(i),
			index:  i,
		}

		st.nodes = append(st.nodes, node)
		st.bySymbol[node.symbol] = node
	}

	return st
}

func (st State) Clone() *State {
	st2 := &State{
		bySymbol: map[byte]*node{},
	}

	for _, node := range st.nodes {
		tmp := *node
		st2.nodes = append(st2.nodes, &tmp)
		st2.bySymbol[tmp.symbol] = &tmp
	}

	return st2
}

// Returns old index
func (st *State) IncrementSymbol(symbol byte) int {
	node := st.nodeFromSymbol(symbol)
	node.count++
	origIndex := node.index

	for iterIndex := origIndex; iterIndex > 0; iterIndex-- {
		prevIndex := iterIndex - 1
		iterNode := st.nodes[iterIndex]
		prevNode := st.nodes[prevIndex]

		if prevNode.count > iterNode.count {
			break
		} else if prevNode.count == iterNode.count && prevNode.symbol < iterNode.symbol {
			break
		}

		st.nodes[iterNode.index], st.nodes[prevNode.index] = prevNode, iterNode
		prevNode.index, iterNode.index = iterIndex, prevIndex
	}

	return origIndex
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

func (st State) nodeFromSymbol(symbol byte) *node {
	return st.bySymbol[symbol]
}
