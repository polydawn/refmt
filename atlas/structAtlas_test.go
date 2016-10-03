package atlas

import (
	"reflect"
	"testing"
)

func TestGenerateStructAtlasAndTraversal(t *testing.T) {
	type BB struct {
		Z string
	}
	type AA struct {
		X string
		Y BB
	}

	atl := GenerateStructAtlas(reflect.TypeOf(AA{}))
	assert(t, "AA atlas length check", 2, len(atl))
	assert(t, "AA atlas field1 name check", "X", atl[0].Name)
	assert(t, "AA atlas field2 name check", "Y", atl[1].Name)
	var aa AA
	rvf1 := atl[0].FieldRoute.TraverseToValue(reflect.ValueOf(aa))
	assert(t, "atlas traverse 1 val type check", "string", rvf1.Type().Name())
	rvf1 = atl[0].FieldRoute.TraverseToValue(reflect.ValueOf(&aa))
	assert(t, "atlas traverse 1 val type check from ref", "string", rvf1.Type().Name())
	rvf1.SetString("magical message") // will panic if starting without '&' because unaddressable
	assert(t, "can use atlas traverse 1 to set value", "magical message", aa.X)

	type CC struct {
		A AA
		BB
	}
	atl = GenerateStructAtlas(reflect.TypeOf(CC{}))
	assert(t, "CC atlas length check", len(atl), 2)
	assert(t, "CC atlas field1 name check", "A", atl[0].Name)
	assert(t, "CC atlas field2 name check", "Z", atl[1].Name) // dives through embed!
	assert(t, "fieldRoute to direct fields is len 1", 1, len(atl[0].FieldRoute))
	assert(t, "fieldRoute to embed fields is len 2", 2, len(atl[1].FieldRoute))
	var cc CC
	atl[1].FieldRoute.TraverseToValue(reflect.ValueOf(&cc)).SetString("mystical string")
	assert(t, "can use atlas traverse to set embed value", "mystical string", cc.BB.Z)
}
