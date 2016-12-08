package atlas

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/polydawn/go-xlate/testutil"
)

func TestReflectWorksLikeIThinkItDoes(t *testing.T) {
	type Fwee string
	Assert(t, "typedefs of primitives Kind() as the prim", reflect.String, reflect.TypeOf(Fwee("")).Kind())
	Assert(t, "typedefs of primitives have Type().Name()", "Fwee", reflect.TypeOf(Fwee("")).Name())
	Assert(t, "primitives also have Type().Name()", "string", reflect.TypeOf("").Name())
}

func TestAtlasAddrFunc(t *testing.T) {
	type BB struct {
		Z string
	}
	type AA struct {
		X string
		Y BB
	}

	atl := &Atlas{Fields: []Entry{
		{Name: "x", AddrFunc: func(v interface{}) interface{} { return &(v.(*AA).X) }},
		{Name: "y", AddrFunc: func(v interface{}) interface{} { return &(v.(*AA).Y) }},
	}}
	aa := AA{
		X: "qwer",
		Y: BB{},
	}
	xp := atl.Fields[0].Grab(&aa)
	Assert(t, "addrfunc yields informative ptr type", "*string", fmt.Sprintf("%T", xp))
	Assert(t, "addrfunc yields readable reference", "qwer", *xp.(*string))
	*xp.(*string) = "zxcv"
	Assert(t, "addrfunc yields writable reference", "zxcv", *xp.(*string))
}

func TestAtlasFieldRoute(t *testing.T) {
	type BB struct {
		Z string
	}
	type AA struct {
		X string
		Y BB
	}

	atl := &Atlas{Fields: []Entry{
		{Name: "x", FieldRoute: []int{0}},
		{Name: "y", FieldRoute: []int{1}},
	}}
	aa := AA{
		X: "qwer",
		Y: BB{},
	}
	xp := atl.Fields[0].Grab(&aa)
	Assert(t, "addrfunc yields informative ptr type", "*string", fmt.Sprintf("%T", xp))
	Assert(t, "addrfunc yields readable reference", "qwer", *xp.(*string))
	*xp.(*string) = "zxcv"
	Assert(t, "addrfunc yields writable reference", "zxcv", *xp.(*string))
}
