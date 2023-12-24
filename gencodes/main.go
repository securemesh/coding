package main

import (
	"log"
	"strings"
)

func main() {
	// Generate all possible values up to 10 bits
	byLen := map[int][]string{
		1: []string{"0", "1"},
	}

	for l := 2; l <= 10; l++ {
		values := []string{}

		for _, v := range byLen[l-1] {
			values = append(values, v+"0", v+"1")
		}

		byLen[l] = values
	}

	limits := map[int]int{
		1:  0,
		2:  0,
		3:  0,
		4:  2,
		5:  8,
		6:  8,
		7:  16,
		8:  32,
		9:  64,
		10: 256,
	}

	total := 0
	short := 0

	for l := 1; l <= 10; l++ {
		vs := byLen[l]
		limit := limits[l]
		values := []string{}

	valueLoop:
		for _, v := range vs {
			if limit == 0 {
				break
			}

			for i := 1; i < l; i++ {
				for _, v2 := range byLen[i] {
					if strings.HasPrefix(v, v2) {
						continue valueLoop
					}
				}
			}

			values = append(values, v)
			limit--
			print(v + "\n")
		}

		byLen[l] = values

		total += len(values)
		if l < 8 {
			short += len(values)
		}
	}

	for l := 1; l <= 10; l++ {
		values := byLen[l]
		log.Printf("%d: %d", l, len(values))
	}

	log.Printf("total=%d", total)
	log.Printf("short=%d", short)
}
