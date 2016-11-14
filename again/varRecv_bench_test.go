package again

import (
	"encoding/json"
	"testing"
)

// Force bench.N to a fixed number.
// This makes it easier to take a peek at a pprof output covering
//  different tests and get a fair(ish) understanding of relative costs.
func forceN(b *testing.B) {
	b.N = 1000000
}

func Benchmark_VarRecvTinyMap(b *testing.B) {
	forceN(b)
	var v interface{}
	x := []Token{
		"k1",
	}
	for i := 0; i < b.N; i++ {
		sink := NewVarReceiver(&v)
		sink.Step(&Token_MapOpen)
		sink.Step(&x[0])
		sink.Step(&x[0])
		sink.Step(&Token_MapClose)
	}
}

func Benchmark_JsonUnmarshalTinyMap(b *testing.B) {
	forceN(b)
	var v interface{}
	byt := []byte(`{"k1":"k1"}`)
	for i := 0; i < b.N; i++ {
		json.Unmarshal(byt, &v)
	}
}
