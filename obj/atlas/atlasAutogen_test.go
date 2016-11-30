package atlas

import (
	"reflect"
	"testing"

	. "github.com/polydawn/go-xlate/testutil"
)

func TestGenerateAtlasAndTraversal(t *testing.T) {
	type BB struct {
		Z string
	}
	type AA struct {
		X string
		Y BB
	}

	atl := GenerateAtlas(reflect.TypeOf(AA{}))
	Assert(t, "AA atlas length check", 2, len(atl.Fields))
	Assert(t, "AA atlas field1 name check", "X", atl.Fields[0].Name)
	Assert(t, "AA atlas field2 name check", "Y", atl.Fields[1].Name)
	var aa AA
	rvf1 := atl.Fields[0].fieldRoute.TraverseToValue(reflect.ValueOf(aa))
	Assert(t, "atlas traverse 1 val type check", "string", rvf1.Type().Name())
	rvf1 = atl.Fields[0].fieldRoute.TraverseToValue(reflect.ValueOf(&aa))
	Assert(t, "atlas traverse 1 val type check from ref", "string", rvf1.Type().Name())
	rvf1.SetString("magical message") // will panic if starting without '&' because unaddressable
	Assert(t, "can use atlas traverse 1 to set value", "magical message", aa.X)

	type CC struct {
		A AA
		BB
	}
	atl = GenerateAtlas(reflect.TypeOf(CC{}))
	Assert(t, "CC atlas length check", len(atl.Fields), 2)
	Assert(t, "CC atlas field1 name check", "A", atl.Fields[0].Name)
	Assert(t, "CC atlas field2 name check", "Z", atl.Fields[1].Name) // dives through embed!
	Assert(t, "fieldRoute to direct fields is len 1", 1, len(atl.Fields[0].fieldRoute))
	Assert(t, "fieldRoute to embed fields is len 2", 2, len(atl.Fields[1].fieldRoute))
	var cc CC
	atl.Fields[1].fieldRoute.TraverseToValue(reflect.ValueOf(&cc)).SetString("mystical string")
	Assert(t, "can use atlas traverse to set embed value", "mystical string", cc.BB.Z)
}
