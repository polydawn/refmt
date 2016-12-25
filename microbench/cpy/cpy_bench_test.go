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

// Sanity check: strings are not noticably different:
//
//	Benchmark_CopyByValue-8                 30000000                45.5 ns/op
//	Benchmark_CopyByRef-8                   2000000000               0.59 ns/op
//	Benchmark_CopyByValue_String-8          20000000                72.3 ns/op
//	Benchmark_CopyByRef_String-8            2000000000               0.60 ns/op
//
// Not commited, but note that there is no sigificant impact from the length of the string.
// Benchmem offers some insight into why:
//
//	Benchmark_CopyByValue_String-8          20000000                73.5 ns/op            16 B/op          1 allocs/op
//
// Evidentally copy-by-value of a string requires a proportionally larger alloc to store the length;
// and furthermore despite being a single alloc, the size in bytes does visibly increase time cost.
func Benchmark_CopyByValue_String(b *testing.B) {
	type Alias interface{}
	var slot Alias
	type StructA struct {
		field string
	}
	type StructB struct {
		field string
	}
	valA := StructA{"alksjdlkjweoihgowihehgioijerg"}
	valB := StructB{}

	for i := 0; i < b.N; i++ {
		slot = valA.field
		valB.field = slot.(string)
	}
	if valB.field != valA.field {
		b.Error("final value of valB wrong")
	}
}

func Benchmark_CopyByRef_String(b *testing.B) {
	type Alias interface{}
	var slot Alias
	type StructA struct {
		field string
	}
	type StructB struct {
		field string
	}
	valA := StructA{"alksjdlkjweoihgowihehgioijerg"}
	valB := StructB{}

	for i := 0; i < b.N; i++ {
		slot = &(valA.field)
		valB.field = *(slot.(*string))
	}
	if valB.field != valA.field {
		b.Error("final value of valB wrong")
	}
}
