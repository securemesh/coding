package main

import (
	"bufio"
	"log"
	"os"

	"github.com/samber/lo"
	"github.com/securemesh/coding"
	"github.com/securemesh/coding/heap"
	"github.com/securemesh/coding/seeds"
)

func main() {
	samples := lo.Must(loadSamples())

	log.Printf("orig=%d", lo.SumBy(samples, func(sample []byte) int {
		return len(sample)
	}))

	def := heap.NewHeap()
	log.Printf("def=%d [%s]", totalLength(def, samples), def)

	chat := seeds.ChatHeap()
	log.Printf("chat=%d [%s]", totalLength(chat, samples), chat)

	opt := optimize(heap.NewHeap(), samples)
	log.Printf("opt=%d [%s]", totalLength(opt, samples), opt)
}

func optimize(h *heap.Heap, samples [][]byte) *heap.Heap {
	for true {
		better := optimize2(h, samples)
		if better == nil {
			return h
		}
		h = better
		log.Printf("\titer=%d [%s]", totalLength(h, samples), h)
	}

	return h
}

type sampleResult struct {
	heap *heap.Heap
	score int
}

func optimize2(baseHeap *heap.Heap, samples [][]byte) *heap.Heap {
	ch := make(chan sampleResult, 100)

	for i := 0; i < 256; i++ {
		s := byte(i)
		go func () {
			h := baseHeap.Clone()
			h.IncrementSymbol(s)
			ch <- sampleResult{
				heap: h,
				score: totalLength(h, samples),
			}
		}()
	}

	var best *heap.Heap = nil
	bestScore := totalLength(baseHeap, samples)

	for i := 0; i < 256; i++ {
		res := <-ch
		if res.score < bestScore {
			best = res.heap
			bestScore = res.score
		}
	}

	return best
}

func totalLength(heap *heap.Heap, samples [][]byte) int {
	return lo.SumBy(samples, func(sample []byte) int {
		return len(coding.Encode(heap.Clone(), sample))
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
