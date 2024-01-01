package main

import (
	"bufio"
	"bytes"
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

	dict := buildDictionary(samples, 1024)

	def := state.NewState()
	log.Printf("def=%d [%s]", totalLength(def, samples), def)

	chat := seeds.ChatState()
	log.Printf("chat=%d [%s]", totalLength(chat, samples), chat)

	opt := optimize(state.NewState(), samples)
	log.Printf("opt=%d [%s]", totalLength(opt, samples), opt)
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
			for j := 2; j < min(5, len(sub)); j++ {
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

func optimize(st *state.State, samples [][]byte) *state.State {
	st.AddSymbol([]byte("it "))

	for true {
		better := optimize2(st, samples)
		if better == nil {
			return st
		}
		st = better
		log.Printf("\titer=%d [%s]", totalLength(st, samples), st)
	}

	return st
}

type sampleResult struct {
	symbol []byte
	state  *state.State
	score  int
}

func optimize2(baseState *state.State, samples [][]byte) *state.State {
	ch := make(chan sampleResult, 100)
	symbols := baseState.Symbols()

	for _, symbol := range symbols {
		res := sampleResult{
			symbol: symbol,
		}

		go func() {
			st := baseState.Clone()
			st.IncrementSymbol(res.symbol)
			res.state = st
			res.score = totalLength(st, samples)
			ch <- res
		}()
	}

	results := []sampleResult{}

	for _ = range symbols {
		results = append(results, <-ch)
	}

	slices.SortFunc(results, func(a, b sampleResult) int { return bytes.Compare(a.symbol, b.symbol) })
	best := slices.MaxFunc(results, func(a, b sampleResult) int { return b.score - a.score })

	if best.score == totalLength(baseState, samples) {
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
	fh, err := os.Open("sms.txt")
	if err != nil {
		return nil, err
	}

	defer fh.Close()

	s := bufio.NewScanner(fh)
	ret := [][]byte{}

	for s.Scan() {
		ret = append(ret, s.Bytes())
	}

	return ret, nil
}
