package xlate

import (
	"fmt"
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

func TestMapperInterfaceDispatch(t *testing.T) {
	type Iface interface{}
	type Top struct {
		X Iface
	}
	type AA struct {
		Y string
	}
	type BB struct {
		Z string
	}

	fmt.Printf(";; %v -- %v\n", reflect.TypeOf(Iface(nil)), nil)
	fmt.Printf(";; %v -- %v\n", reflect.TypeOf(Iface(BB{})), reflect.TypeOf(Iface(BB{})).Kind())
	var fuk Iface
	fmt.Printf(";; %v -- %v\n", reflect.TypeOf(fuk), nil)
	fmt.Printf(";; %v -- %v\n", reflect.TypeOf(&fuk), reflect.TypeOf(&fuk).Kind())
	fmt.Printf(";; %v -- %v\n", reflect.TypeOf(&fuk).Elem(), reflect.TypeOf(&fuk).Elem().Kind())

	type testRow struct {
		title     string
		thing     interface{}
		mapper    *Mapper
		expect    tok.Tokens
		expectErr error
	}
	for _, tr := range []testRow{{
		title: "test a single thing at its zero value",
		thing: Top{AA{"whee"}},
		mapper: NewMapper(MapperSetup{
			{Top{}, func(dispatch *Mapper, dest Destination, input interface{}) {
				dispatch.Map(dest, Iface(input.(Top).X))
			}},
			//{Iface(nil), Map_Wildcard_toString},
			{AA{}, Map_Wildcard_toString},
		}),
		expect: tok.Tokens{
			{tok.TokenKind_ValString, "{}"},
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
			t.Errorf("test %q FAILED:\n\texpected     %#v\n\terrored with %v",
				tr.title, tr.expectErr, err)
		}
		if !stringyEquality(tr.expect, tokDest.Tokens) {
			t.Errorf("test %q FAILED:\n\texpected: %#v\n\tactual:   %#v",
				tr.title, tr.expect, tokDest.Tokens)
		}
	}
}
