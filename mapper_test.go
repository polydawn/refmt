package xlate

import (
	"testing"
)

func TestMapperSetup(t *testing.T) {
	type BB struct {
		Z string
	}
	type AA struct {
		X string
		Y BB
	}

	type testRow struct {
		title  string
		setup  MapperSetup
		expect interface{} // err type or nil
	}
	for _, tr := range []testRow{{
		title:  "empty mapper setups are acceptable",
		setup:  MapperSetup{},
		expect: nil,
	}, {
		title: "short simple mapper setups are acceptable",
		setup: MapperSetup{
			{AA{}, Map_Wildcard_toString},
			{BB{}, Map_Wildcard_toString},
		},
		expect: nil,
	}, {
		title: "mapper setups with nil funcs are rejected",
		setup: MapperSetup{
			{AA{}, nil},
		},
		expect: &ErrNilMappingFunc{MapperSetupRow{AA{}, nil}},
	}, {
		title: "mapper setups with repeat types are rejected",
		setup: MapperSetup{
			{AA{}, Map_Wildcard_toStringOfType},
			{AA{}, Map_Wildcard_toString},
		},
		expect: &ErrMapperSetupNotUnique{MapperSetupRow{AA{}, Map_Wildcard_toString}},
	}} {
		err := capturePanics(func() {
			NewMapper(tr.setup)
		})
		if !stringyEquality(tr.expect, err) {
			t.Errorf("test %q FAILED:\n\texpected     %#v\n\terrored with %#v",
				tr.title, tr.expect, err)
		}
	}
}

func TestMapperBasics(t *testing.T) {
	type BB struct {
		Z string
	}
	type AA struct {
		X string
		Y BB
	}

	type testRow struct {
		title  string
		thing  interface{}
		mapper *Mapper
		expect string
	}
	for _, tr := range []testRow{{
		title: "test a single thing at its zero value",
		thing: AA{},
		mapper: NewMapper(MapperSetup{
			{BB{}, Map_Wildcard_toStringOfType},
		}),
		expect: "BB",
	}} {
		_ = tr
	}
}
