/*
	Guides the design of `tok.Token`:
	test whether it's important to yield pointers to the values of interest,
	versus simply putting the values in an `interface{}` slot.

	(It is.)

	Though syntactically irritating to use of pointers to primitives for all tokens,
	this bench demonstrates that doing so avoids a source of allocs,
	and thus has significant performance implications.
*/
package bench

import (
	"testing"
)

// Std:  Benchmark_CopyByValue-8         30000000                43.4 ns/op
// noGC: Benchmark_CopyByValue-8         30000000                34.0 ns/op
// mem:  Benchmark_CopyByValue-8         30000000                44.4 ns/op             8 B/op          1 allocs/op
func Benchmark_CopyByValue(b *testing.B) {
	type Alias interface{}
	var slot Alias
	type StructA struct {
		field int
	}
	type StructB struct {
		field int
	}
	valA := StructA{4}
	valB := StructB{}

	for i := 0; i < b.N; i++ {
		slot = valA.field
		valB.field = slot.(int)
	}
	if valB.field != 4 {
		b.Error("final value of valB wrong")
	}
}

// Std:  Benchmark_CopyByRef-8           2000000000               0.59 ns/op
// noGC: Benchmark_CopyByRef-8           2000000000               0.59 ns/op
// mem:  Benchmark_CopyByRef-8           2000000000               0.59 ns/op            0 B/op          0 allocs/op
func Benchmark_CopyByRef(b *testing.B) {
	type Alias interface{}
	var slot Alias
	type StructA struct {
		field int
	}
	type StructB struct {
		field int
	}
	valA := StructA{4}
	valB := StructB{}

	for i := 0; i < b.N; i++ {
		slot = &(valA.field)
		valB.field = *(slot.(*int))
	}
	if valB.field != 4 {
		b.Error("final value of valB wrong")
	}
}
