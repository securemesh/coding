package main

import (
	"bufio"
	"log"
	"os"

	"github.com/samber/lo"
	"github.com/securemesh/coding"
	"github.com/securemesh/coding/state"
	"github.com/securemesh/coding/seeds"
)

func main() {
	samples := lo.Must(loadSamples())

	log.Printf("orig=%d", lo.SumBy(samples, func(sample []byte) int {
		return len(sample)
	}))

	def := state.NewState()
	log.Printf("def=%d [%s]", totalLength(def, samples), def)

	chat := seeds.ChatState()
	log.Printf("chat=%d [%s]", totalLength(chat, samples), chat)

	opt := optimize(state.NewState(), samples)
	log.Printf("opt=%d [%s]", totalLength(opt, samples), opt)
}

func optimize(st *state.State, samples [][]byte) *state.State {
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
	state *state.State
	score int
}

func optimize2(baseState *state.State, samples [][]byte) *state.State {
	ch := make(chan sampleResult, 100)

	for i := 0; i < 256; i++ {
		s := byte(i)
		go func () {
			st := baseState.Clone()
			st.IncrementSymbol(s)
			ch <- sampleResult{
				state: st,
				score: totalLength(st, samples),
			}
		}()
	}

	var best *state.State = nil
	bestScore := totalLength(baseState, samples)

	for i := 0; i < 256; i++ {
		res := <-ch
		if res.score < bestScore {
			best = res.state
			bestScore = res.score
		}
	}

	return best
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
