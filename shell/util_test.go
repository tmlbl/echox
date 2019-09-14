package shell

import (
	"testing"
)

func TestEndsWith(t *testing.T) {
	for _, tst := range []struct {
		b      []byte
		sub    []byte
		expect bool
	}{
		{
			b:      []byte("hello world"),
			sub:    []byte("world"),
			expect: true,
		},
		{
			b:      []byte("hello world"),
			sub:    []byte("welp"),
			expect: false,
		},
		{
			b:      []byte(""),
			sub:    []byte("welp"),
			expect: false,
		},
	} {
		if endswith(tst.b, tst.sub) != tst.expect {
			t.Errorf("Expected %s ends with %s to be %v",
				tst.b, tst.sub, tst.expect)
		}
	}
}
