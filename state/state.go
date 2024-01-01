package state

import (
	"bytes"
	"fmt"
	"strings"
)

type node struct {
	symbol []byte
	count  int
	index  int
}

type trieNode struct {
	node     *node
	children [16]*trieNode
}

type State struct {
	nodes []*node
	root  *trieNode
}

func NewState() *State {
	st := &State{
		root: &trieNode{},
	}

	for i := 0; i < 256; i++ {
		node := &node{
			symbol: []byte{byte(i)},
			index:  i,
		}

		st.nodes = append(st.nodes, node)
		st.insertTrieNode(node)
	}

	return st
}

func (st State) Clone() *State {
	st2 := &State{
		root: &trieNode{},
	}

	for _, node := range st.nodes {
		node2 := *node
		st2.nodes = append(st2.nodes, &node2)
		st2.insertTrieNode(&node2)
	}

	return st2
}

func (st *State) AddSymbol(symbol []byte) {
	node := &node{
		symbol: symbol,
		index:  len(st.nodes),
	}

	st.nodes = append(st.nodes, node)
	st.insertTrieNode(node)
}

// Returns (symbol_length, old_index)
func (st *State) IncrementSymbol(symbols []byte) (int, int) {
	node := st.nodeFromSymbols(symbols)
	node.count++
	origIndex := node.index

	for iterIndex := origIndex; iterIndex > 0; iterIndex-- {
		prevIndex := iterIndex - 1
		iterNode := st.nodes[iterIndex]
		prevNode := st.nodes[prevIndex]

		if prevNode.count > iterNode.count {
			break
		} else if prevNode.count == iterNode.count && bytes.Compare(prevNode.symbol, iterNode.symbol) < 0 {
			break
		}

		st.nodes[iterNode.index], st.nodes[prevNode.index] = prevNode, iterNode
		prevNode.index, iterNode.index = iterIndex, prevIndex
	}

	return len(node.symbol), origIndex
}

func (st State) String() string {
	strs := []string{}

	for _, node := range st.nodes {
		if node.count == 0 {
			break
		}

		strs = append(strs, fmt.Sprintf("%#U=%d", node.symbol, node.count))
	}

	return strings.Join(strs, ", ")
}

func (st State) Symbols() [][]byte {
	ret := [][]byte{}

	for _, node := range st.nodes {
		ret = append(ret, node.symbol)
	}

	return ret
}

func (st State) insertTrieNode(node *node) {
	iter := st.root

	for _, b := range node.symbol {
		nibbleTrieNode := iter.getOrInsertChild((b & 0xf0) >> 4)
		iter = nibbleTrieNode.getOrInsertChild(b & 0x0f)
	}

	iter.node = node
}

func (st State) nodeFromSymbols(symbols []byte) *node {
	var lastFound *node
	iter := st.root

	for _, b := range symbols {
		nibbleTrieNode := iter.children[(b&0xf0)>>4]
		if nibbleTrieNode == nil {
			break
		}

		iter = nibbleTrieNode.children[b&0x0f]
		if iter == nil {
			break
		}

		if iter.node != nil {
			lastFound = iter.node
		}
	}

	return lastFound
}

func (tn *trieNode) getOrInsertChild(val byte) *trieNode {
	child := tn.children[val]
	if child != nil {
		return child
	}

	child = &trieNode{}
	tn.children[val] = child
	return child
}
