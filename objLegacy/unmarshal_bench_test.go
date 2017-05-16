package objLegacy

import (
	"encoding/json"
	"testing"

	. "github.com/polydawn/refmt/tok"
)

// Force bench.N to a fixed number.
// This makes it easier to take a peek at a pprof output covering
//  different tests and get a fair(ish) understanding of relative costs.
func forceN(b *testing.B) {
	b.N = 1000000
}

func Benchmark_UnmarshalTinyMap(b *testing.B) {
	forceN(b)
	var v interface{}
	x := []Token{
		{Type: TString, Str: "k1"},
	}
	for i := 0; i < b.N; i++ {
		sink := NewUnmarshaler(&v)
		sink.Step(&Token{Type: TMapOpen})
		sink.Step(&x[0])
		sink.Step(&x[0])
		sink.Step(&Token{Type: TMapClose})
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

func Benchmark_UnmarshalLongArray(b *testing.B) {
	forceN(b)
	var v interface{}
	x := []Token{
		{Type: TArrOpen},
		{Type: TInt, Int: 1}, {Type: TInt, Int: 2}, {Type: TInt, Int: 3}, {Type: TInt, Int: 4},
		{Type: TInt, Int: 5}, {Type: TInt, Int: 6}, {Type: TInt, Int: 7}, {Type: TInt, Int: 8},
		{Type: TInt, Int: 9}, {Type: TInt, Int: 10}, {Type: TInt, Int: 11}, {Type: TInt, Int: 12},
		{Type: TInt, Int: 13}, {Type: TInt, Int: 14}, {Type: TInt, Int: 15}, {Type: TInt, Int: 16},
		{Type: TArrClose},
	}
	for i := 0; i < b.N; i++ {
		sink := NewUnmarshaler(&v)
		for j := 0; j < len(x); j++ {
			sink.Step(&x[j])
		}
	}
}

func Benchmark_JsonUnmarshalLongArray(b *testing.B) {
	forceN(b)
	var v interface{}
	byt := []byte(`[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16]`)
	for i := 0; i < b.N; i++ {
		json.Unmarshal(byt, &v)
	}
}
