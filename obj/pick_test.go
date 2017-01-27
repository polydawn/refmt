package obj

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPick(t *testing.T) {
	type CC string
	type BB struct {
		Z string
	}
	type AA struct {
		X string
		Y BB
	}
	type Sly []AA

	suite := &Suite{
		map[reflect.Type]func() MarshalMachine{
			reflect.TypeOf(CC("")): func() MarshalMachine { return &MarshalMachineLiteral{} },
		},
	}
	tt := []struct {
		title string
		value interface{} // Always a ref.  Can't take addr of this without getting `*interface{}` as typeinfo, which misses point of test.
	}{{
		title: "a struct should hello",
		value: &BB{},
	}, {
		title: "a slice should hit a slice wildcard delegating to its component's machine",
		value: &[]BB{},
	}, {
		title: "a slice of pointers should hit slice wildcard delegating to its component's machine",
		value: &[]*BB{},
	}, {
		title: "a typedef'd slice should hit its own machine",
		value: &[]Sly{},
		//}, { // Perhaps later
		//	title: "odd types like funcs should error",
		//	value: func() {},
	}}
	for _, tr := range tt {
		// Call pick with the address of our value.  Remember, this is what `AddrFunc` always yields.
		m := suite.maybePickMarshalMachine(tr.value)
		// Eh?
		t.Logf("test %q:\n\tvalue %#v yielded machine %T\n", tr.title, tr.value, m)
	}

	json.Marshal(nil)
}
