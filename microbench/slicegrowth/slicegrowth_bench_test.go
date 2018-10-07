package slicegrowth

import (
	"reflect"
	"testing"
)

// okay so this append method *is* amortizing sanely.
// but we still have two allocs per set here, why?
func BenchmarkReflectAppend_Naive(b *testing.B) {
	x := []int{}
	rv := reflect.ValueOf(&x).Elem()
	rt_val := reflect.TypeOf(0)

	for i := 0; i < b.N; i++ {
		rv.Set(reflect.Append(rv, reflect.Zero(rt_val)))
	}
	b.Logf("len, cap = %d, %d", len(x), cap(x))
}

// yep that was one of 'em.  reflect.Zero causes malloc woo.
func BenchmarkReflectAppend_ConstantZero(b *testing.B) {
	x := []int{}
	rv := reflect.ValueOf(&x).Elem()
	rt_val := reflect.TypeOf(0)
	rv_valZero := reflect.Zero(rt_val)

	for i := 0; i < b.N; i++ {
		rv.Set(reflect.Append(rv, rv_valZero))
	}
	b.Logf("len, cap = %d, %d", len(x), cap(x))
}

// faster.  but still averaging an alloc per append, how and why
func BenchmarkReflectAppend_ConstZeroAndFinalSet(b *testing.B) {
	x := []int{}
	rv := reflect.ValueOf(&x).Elem()
	rt_val := reflect.TypeOf(0)
	rv_valZero := reflect.Zero(rt_val)

	rv_moving := rv
	for i := 0; i < b.N; i++ {
		rv_moving = reflect.Append(rv_moving, rv_valZero)
	}
	rv.Set(rv_moving)
	b.Logf("x  len, cap = %d, %d", len(x), cap(x))
}
