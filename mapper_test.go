package xlate

import (
	"reflect"
	"testing"

	"github.com/polydawn/go-xlate/dest/tok"
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
		title     string
		thing     interface{}
		mapper    *Mapper
		expect    tok.Tokens
		expectErr error
	}
	for _, tr := range []testRow{{
		title: "test a single thing at its zero value",
		thing: BB{},
		mapper: NewMapper(MapperSetup{
			{BB{}, Map_Wildcard_toString},
		}),
		expect: tok.Tokens{
			{tok.TokenKind_ValString, "{}"},
		},
	}, {
		title: "test a single thing that has no handler",
		thing: AA{},
		mapper: NewMapper(MapperSetup{
			{BB{}, Map_Wildcard_toString},
		}),
		expectErr: &ErrMissingMappingFunc{reflect.TypeOf(AA{})},
	}, {
		title: "test a custom func plus recursion at its zero value",
		thing: AA{},
		mapper: NewMapper(MapperSetup{
			// a fairly bizarre hack of recursion, trying to avoid triggering anything too fancy yet.
			{AA{}, func(dispatch *Mapper, dest Destination, input interface{}) {
				dest.OpenMap()
				dest.WriteMapKey("type")
				dest.WriteString("AA")
				dest.WriteMapKey("data")
				dispatch.Map(dest, input.(AA).Y)
				dest.CloseMap()
			}},
			{BB{}, Map_Wildcard_toString},
		}),
		expect: tok.Tokens{
			{tok.TokenKind_OpenMap, nil},
			{tok.TokenKind_MapKey, "type"},
			{tok.TokenKind_ValString, "AA"},
			{tok.TokenKind_MapKey, "data"},
			{tok.TokenKind_ValString, "{}"},
			{tok.TokenKind_CloseMap, nil},
		},
	}} {
		tokDest := tok.NewDestination()
		err := capturePanics(func() {
			tr.mapper.Map(
				tokDest,
				tr.thing,
			)
		})
		if !stringyEquality(tr.expectErr, err) {
			t.Errorf("test %q FAILED:\n\texpected     %#v\n\terrored with %#v",
				tr.title, tr.expectErr, err)
		}
		if !stringyEquality(tr.expect, tokDest.Tokens) {
			t.Errorf("test %q FAILED:\n\texpected: %#v\n\tactual:   %#v",
				tr.title, tr.expect, tokDest.Tokens)
		}
	}
}
