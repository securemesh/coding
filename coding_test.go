package coding_test

import (
	"bufio"
	"os"
	"testing"

	"github.com/samber/lo"
	"github.com/securemesh/coding"
)

func TestSimple(t *testing.T) {
	msg := []byte("this is a test. this is only a test.")
	encoded := coding.Encode(coding.ChatHeap(), msg)
	t.Logf("orig=%d encoded=%d", len(msg), len(encoded))
}

func TestSMS(t *testing.T) {
	fh := lo.Must(os.Open("sms.txt"))
	defer fh.Close()

	s := bufio.NewScanner(fh)

	orig := 0
	encoded := 0

	for s.Scan() {
		msg := s.Bytes()
		e := coding.Encode(coding.ChatHeap(), msg)
		orig += len(msg)
		encoded += len(e)
	}

	t.Logf("orig=%d encoded=%d", orig, encoded)
}
