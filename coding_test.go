package coding_test

import (
	"encoding/csv"
	"io"
	"os"
	"testing"

	"github.com/samber/lo"
	"github.com/securemesh/coding"
	"github.com/securemesh/coding/seeds"
)

func TestSimple(t *testing.T) {
	msg := []byte("this is a test. this is only a test.")
	encoded := coding.Encode(seeds.ChatState(), msg)
	t.Logf("orig=%d encoded=%d", len(msg), len(encoded))
}

func TestSMS(t *testing.T) {
	fh := lo.Must(os.Open("text.csv"))
	defer fh.Close()

	r := csv.NewReader(fh)

	orig := 0
	encoded := 0

	for true {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		msg := []byte(row[0])
		e := coding.Encode(seeds.ChatState(), msg)
		orig += len(msg)
		encoded += len(e)
	}

	t.Logf("orig=%d encoded=%d", orig, encoded)
}
