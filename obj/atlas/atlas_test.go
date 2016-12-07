package atlas

import (
	"fmt"
	"testing"

	. "github.com/polydawn/go-xlate/testutil"
)

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
	Assert(t, "addrfunc yields readable ptr type", "*string", fmt.Sprintf("%T", xp))
	Assert(t, "addrfunc yields readable reference", "qwer", *xp.(*string))
	*xp.(*string) = "zxcv"
	Assert(t, "addrfunc yields writable reference", "zxcv", *xp.(*string))
}
