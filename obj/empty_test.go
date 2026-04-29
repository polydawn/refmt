package obj

import (
	"reflect"
	"testing"
)

type T1 struct{}

// T2 and T3 carry fields specifically to exercise isEmptyValue against
// structs that look non-trivial but have only zero-valued nilable members.
type T2 struct {
	//lint:ignore U1000 used via reflection in TestIsEmptyValue
	array []byte
}
type T3 struct {
	//lint:ignore U1000 used via reflection in TestIsEmptyValue
	array []*T3
}

func TestIsEmptyValue(t *testing.T) {
	isEmpty := func(v interface{}) {
		if !isEmptyValue(reflect.ValueOf(v)) {
			t.Fatalf("expected value of type %T to be empty", v)
		}
	}
	isEmpty(T1{})
	isEmpty(T2{})
	isEmpty(T3{})
}
