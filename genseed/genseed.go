package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"os"
	"slices"

	"github.com/samber/lo"
	"github.com/securemesh/coding"
	"github.com/securemesh/coding/seeds"
	"github.com/securemesh/coding/state"
)

func main() {
	samples := lo.Must(loadSamples())

	log.Printf("orig=%d", lo.SumBy(samples, func(sample []byte) int {
		return len(sample)
	}))

	def := state.NewState()
	log.Printf("def=%d {%s}", totalLength(def, samples), def)

	chat := seeds.ChatState()
	log.Printf("chat=%d {%s}", totalLength(chat, samples), chat)

	words := buildDictionary(samples, 1024)
	dict := optimizeDict(state.NewState(), samples, words)
	log.Printf("dict=%d {%#U}", totalLength(dict, samples), dict.Symbols()[256:])

	opt := optimize(dict, samples)
	log.Printf("opt=%d {%s}", totalLength(opt, samples), opt)
}

type pair struct {
	symbol []byte
	count  int
}

func buildDictionary(samples [][]byte, num int) [][]byte {
	counts := map[uint64]*pair{}

	for _, sample := range samples {
		for i := 0; i < len(sample); i++ {
			sub := sample[i:]
			for j := 2; j < min(9, len(sub)); j++ {
				sub2 := sub[:j]
				k := toUint64(sub2)

				p := counts[k]
				if p == nil {
					counts[k] = &pair{
						symbol: sub2,
						count:  1,
					}
				} else {
					p.count++
				}
			}
		}
	}

	pairs := []*pair{}

	for _, p := range counts {
		pairs = append(pairs, p)
	}

	slices.SortFunc(pairs, func(a, b *pair) int { return bytes.Compare(a.symbol, b.symbol) })
	slices.SortStableFunc(pairs, func(a, b *pair) int { return b.score() - a.score() })

	ret := [][]byte{}

	for i := 0; i < num && i < len(pairs); i++ {
		ret = append(ret, pairs[i].symbol)
	}

	return ret
}

func toUint64(bs []byte) uint64 {
	var ret uint64

	for _, b := range bs {
		ret = (ret << 8) | uint64(b)
	}

	return ret
}

func (p pair) score() int {
	return p.count * ((len(p.symbol) * 8) - 11)
}

type sampleResult struct {
	symbol []byte
	state  *state.State
	score  int
}

func optimizeDict(st *state.State, samples [][]byte, dict [][]byte) *state.State {
	log.Printf("optDict:")

	for len(dict) > 0 {
		better := optimizeDict2(st, samples, dict)
		if better == nil {
			return st
		}
		st = better.state
		log.Printf("\titer=%d {%#U}", totalLength(st, samples), better.symbol)

		dict2 := [][]byte{}
		for _, symbol := range dict {
			if !bytes.Equal(symbol, better.symbol) {
				dict2 = append(dict2, symbol)
			}
		}
		dict = dict2
	}

	return st
}

func optimizeDict2(baseState *state.State, samples [][]byte, dict [][]byte) *sampleResult {
	ch := make(chan *sampleResult, 100)

	for _, symbol := range dict {
		res := &sampleResult{
			symbol: symbol,
		}

		go func() {
			res.state = baseState.Clone()
			res.state.AddSymbol(res.symbol)
			res.score = totalLength(res.state, samples)
			ch <- res
		}()
	}

	results := []*sampleResult{}

	for _ = range dict {
		results = append(results, <-ch)
	}

	slices.SortFunc(results, func(a, b *sampleResult) int { return bytes.Compare(a.symbol, b.symbol) })
	best := slices.MaxFunc(results, func(a, b *sampleResult) int { return b.score - a.score })

	if best.score >= totalLength(baseState, samples) {
		return nil
	}

	return best
}

func optimize(st *state.State, samples [][]byte) *state.State {
	log.Printf("opt:")

	for true {
		better := optimize2(st, samples)
		if better == nil {
			return st
		}
		st = better
		log.Printf("\titer=%d {%s}", totalLength(st, samples), st)
	}

	return st
}

func optimize2(baseState *state.State, samples [][]byte) *state.State {
	ch := make(chan *sampleResult, 100)
	symbols := baseState.Symbols()

	for _, symbol := range symbols {
		res := &sampleResult{
			symbol: symbol,
		}

		go func() {
			res.state = baseState.Clone()
			res.state.IncrementSymbol(res.symbol)
			res.score = totalLength(res.state, samples)
			ch <- res
		}()
	}

	results := []*sampleResult{}

	for _ = range symbols {
		results = append(results, <-ch)
	}

	slices.SortFunc(results, func(a, b *sampleResult) int { return bytes.Compare(a.symbol, b.symbol) })
	best := slices.MaxFunc(results, func(a, b *sampleResult) int { return b.score - a.score })

	if best.score >= totalLength(baseState, samples) {
		return nil
	}

	return best.state
}

func totalLength(st *state.State, samples [][]byte) int {
	return lo.SumBy(samples, func(sample []byte) int {
		return len(coding.Encode(st.Clone(), sample))
	})
}

func loadSamples() ([][]byte, error) {
	fh, err := os.Open("text.csv")
	if err != nil {
		return nil, err
	}

	defer fh.Close()

	r := csv.NewReader(fh)
	ret := [][]byte{}

	for true {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		ret = append(ret, []byte(row[0]))
	}

	return ret, nil
}
