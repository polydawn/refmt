package obj

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/polydawn/refmt/tok"
	"github.com/polydawn/refmt/tok/fixtures"
)

var skipMe = fmt.Errorf("skipme")

type marshalResults struct {
	title string
	// Yields a value to hand to the marshaller.
	// A func returning a wildcard is used rather than just an `interface{}`, because `&target` conveys very different type information.
	valueFn func() interface{}

	expectErr error
	errString string
}
type unmarshalResults struct {
	title string
	// Yields the handle we should give to the unmarshaller to fill.
	// Like `valueFn`, the indirection here is to help
	slotFn func() interface{}

	// Yields the value we will compare the unmarshal result against.
	// A func returning a wildcard is used rather than just an `interface{}`, because `&target` conveys very different type information.
	valueFn   func() interface{}
	expectErr error
	errString string
}

type tObjStr struct {
	X string
}

type tObjStr2 struct {
	X string
	Y string
}

type tObjK struct {
	K []tObjK2
}
type tObjK2 struct {
	K2 int
}

type tDefStr string
type tDefInt int
type tDefBytes []byte

var objFixtures = []struct {
	title string

	// The serial sequence of tokens the value is isomorphic to.
	sequence fixtures.Sequence

	// The suite of mappings to use.
	atlas atlas.Atlas

	// The results to expect from various marshalling starting points.
	// This is a slice because we occasionally have several different kinds of objects
	// which we expect will converge on the same token fixture given the same atlas.
	marshalResults []marshalResults

	// The results to expect from various unmarshal situations.
	// This is a slice because unmarshal may have different outcomes (usually,
	// erroring vs not) depending on the type of value it was given to populate.
	unmarshalResults []unmarshalResults
}{
	{title: "string literal",
		sequence: fixtures.SequenceMap["flat string"],
		marshalResults: []marshalResults{
			{title: "from string literal",
				valueFn: func() interface{} { str := "value"; return str }},
			{title: "from *string",
				valueFn: func() interface{} { str := "value"; return &str }},
			{title: "from string in iface slot",
				valueFn: func() interface{} { var iface interface{}; iface = "value"; return iface }},
			{title: "from string in *iface slot",
				valueFn: func() interface{} { var iface interface{}; iface = "value"; return &iface }},
			{title: "from *string in iface slot",
				valueFn: func() interface{} { str := "value"; var iface interface{}; iface = &str; return iface }},
			{title: "from *string in *iface slot",
				valueFn: func() interface{} { str := "value"; var iface interface{}; iface = &str; return &iface }},
		},
		unmarshalResults: []unmarshalResults{
			{title: "into string",
				slotFn:    func() interface{} { var str string; return str },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf("")}},
			{title: "into *string",
				slotFn:  func() interface{} { var str string; return &str },
				valueFn: func() interface{} { str := "value"; return str }},
			{title: "into wildcard",
				slotFn:    func() interface{} { var v interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(interface{}(nil))}},
			{title: "into *wildcard",
				slotFn:  func() interface{} { var v interface{}; return &v },
				valueFn: func() interface{} { str := "value"; return str }},
			{title: "into map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(map[string]interface{}(nil))}},
			{title: "into *map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return &v },
				expectErr: ErrUnmarshalIncongruent{Token{Type: TString, Str: "value"}, reflect.ValueOf(map[string]interface{}(nil))}},
			{title: "into []iface",
				slotFn:    func() interface{} { var v []interface{}; return v },
				expectErr: skipMe},
			{title: "into *[]iface",
				slotFn:    func() interface{} { var v []interface{}; return &v },
				expectErr: skipMe},
		},
	},
	{title: "object with one string field, with atlas entry",
		sequence: fixtures.SequenceMap["single row map"],
		atlas: atlas.MustBuild(
			atlas.BuildEntry(tObjStr{}).StructMap().
				AddField("X", atlas.StructMapEntry{SerialName: "key"}).
				Complete(),
		),
		marshalResults: []marshalResults{
			{title: "from object with one field",
				valueFn: func() interface{} { return tObjStr{"value"} }},
			{title: "from map[str]iface with one entry",
				valueFn: func() interface{} { return map[string]interface{}{"key": "value"} }},
			{title: "from map[str]str with one entry",
				valueFn: func() interface{} { return map[string]string{"key": "value"} }},
			{title: "from *map[str]str",
				valueFn: func() interface{} { m := map[string]string{"key": "value"}; return &m }},
		},
		unmarshalResults: []unmarshalResults{
			{title: "into string",
				slotFn:    func() interface{} { var str string; return str },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf("")}},
			{title: "into *string",
				slotFn:    func() interface{} { var str string; return &str },
				expectErr: ErrUnmarshalIncongruent{Token{Type: TMapOpen, Length: 1}, reflect.ValueOf("")}},
			{title: "into wildcard",
				slotFn:    func() interface{} { var v interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(interface{}(nil))}},
			{title: "into *wildcard",
				slotFn:  func() interface{} { var v interface{}; return &v },
				valueFn: func() interface{} { return map[string]interface{}{"key": "value"} }},
			{title: "into map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(map[string]interface{}(nil))}},
			{title: "into made map[str]iface",
				slotFn:  func() interface{} { v := make(map[string]interface{}); return v },
				valueFn: func() interface{} { return map[string]interface{}{"key": "value"} }},
			{title: "into *map[str]iface",
				slotFn:  func() interface{} { var v map[string]interface{}; return &v },
				valueFn: func() interface{} { return map[string]interface{}{"key": "value"} }},
			{title: "into *map[str]str",
				slotFn:  func() interface{} { var v map[string]string; return &v },
				valueFn: func() interface{} { return map[string]string{"key": "value"} }},
			{title: "into []iface",
				slotFn:    func() interface{} { var v []interface{}; return v },
				expectErr: skipMe},
			{title: "into *[]iface",
				slotFn:    func() interface{} { var v []interface{}; return &v },
				expectErr: skipMe},
		},
	},
	{title: "object with two string fields, with atlas entry",
		sequence: fixtures.SequenceMap["duo row map"],
		atlas: atlas.MustBuild(
			atlas.BuildEntry(tObjStr2{}).StructMap().
				AddField("X", atlas.StructMapEntry{SerialName: "key"}).
				AddField("Y", atlas.StructMapEntry{SerialName: "k2"}).
				Complete(),
		),
		marshalResults: []marshalResults{
			{title: "from object with two fields",
				valueFn: func() interface{} { return tObjStr2{"value", "v2"} }},
		},
		unmarshalResults: []unmarshalResults{
			{title: "into string",
				slotFn:    func() interface{} { var str string; return str },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf("")}},
			{title: "into *string",
				slotFn:    func() interface{} { var str string; return &str },
				expectErr: ErrUnmarshalIncongruent{Token{Type: TMapOpen, Length: 2}, reflect.ValueOf("")}},
			{title: "into wildcard",
				slotFn:    func() interface{} { var v interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(interface{}(nil))}},
			{title: "into *wildcard",
				slotFn:  func() interface{} { var v interface{}; return &v },
				valueFn: func() interface{} { return map[string]interface{}{"key": "value", "k2": "v2"} }},
			{title: "into map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(map[string]interface{}(nil))}},
			{title: "into made map[str]iface",
				slotFn:  func() interface{} { v := make(map[string]interface{}); return v },
				valueFn: func() interface{} { return map[string]interface{}{"key": "value", "k2": "v2"} }},
			{title: "into *map[str]iface",
				slotFn:  func() interface{} { var v map[string]interface{}; return &v },
				valueFn: func() interface{} { return map[string]interface{}{"key": "value", "k2": "v2"} }},
			{title: "into *map[str]str",
				slotFn:  func() interface{} { var v map[string]string; return &v },
				valueFn: func() interface{} { return map[string]string{"key": "value", "k2": "v2"} }},
			{title: "into []iface",
				slotFn:    func() interface{} { var v []interface{}; return v },
				expectErr: skipMe},
			{title: "into *[]iface",
				slotFn:    func() interface{} { var v []interface{}; return &v },
				expectErr: skipMe},
		},
	},
	{title: "empty primitive arrays",
		sequence: fixtures.SequenceMap["empty array"],
		marshalResults: []marshalResults{
			{title: "from int array",
				valueFn: func() interface{} { return [0]int{} }},
			{title: "from int slice",
				valueFn: func() interface{} { return []int{} }},
			{title: "from iface array",
				valueFn: func() interface{} { return [0]interface{}{} }},
			{title: "from iface slice",
				valueFn: func() interface{} { return []interface{}{} }},
		},
		unmarshalResults: []unmarshalResults{
			{title: "into string",
				slotFn:    func() interface{} { var str string; return str },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf("")}},
			{title: "into *string",
				slotFn:    func() interface{} { var str string; return &str },
				expectErr: ErrUnmarshalIncongruent{Token{Type: TArrOpen, Length: 0}, reflect.ValueOf("")}},
			{title: "into wildcard",
				slotFn:    func() interface{} { var v interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(interface{}(nil))}},
			{title: "into *wildcard",
				slotFn:  func() interface{} { var v interface{}; return &v },
				valueFn: func() interface{} { return []interface{}{} }},
			{title: "into map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return v },
				expectErr: skipMe},
			{title: "into made map[str]iface",
				slotFn:    func() interface{} { v := make(map[string]interface{}); return v },
				expectErr: skipMe},
			{title: "into *map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return &v },
				expectErr: skipMe},
			{title: "into *map[str]str",
				slotFn:    func() interface{} { var v map[string]string; return &v },
				expectErr: skipMe},
			{title: "into []iface",
				slotFn: func() interface{} { var v []interface{}; return v },
				// array/slice direct: theoretically possible, as long as it's short enough.  but not supported right now.
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf([]interface{}{})}},
			{title: "into *[]iface",
				slotFn:  func() interface{} { var v []interface{}; return &v },
				valueFn: func() interface{} { return []interface{}{} }},
			{title: "into []str",
				slotFn:    func() interface{} { var v []string; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf([]string{})}},
			{title: "into *[]str",
				slotFn:  func() interface{} { var v []string; return &v },
				valueFn: func() interface{} { return []string{} }},
			{title: "into []int",
				slotFn:    func() interface{} { var v []int; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf([]int{})}},
			{title: "into *[]int",
				slotFn:  func() interface{} { var v []int; return &v },
				valueFn: func() interface{} { return []int{} }}, // fine *despite the type mismatch* because with no tokens, well, nobody's the wiser.
		},
	},
	{title: "short primitive arrays",
		sequence: fixtures.SequenceMap["duo entry array"],
		marshalResults: []marshalResults{
			{title: "from str array",
				valueFn: func() interface{} { return [2]string{"value", "v2"} }},
			{title: "from str slice",
				valueFn: func() interface{} { return []string{"value", "v2"} }},
			{title: "from iface array",
				valueFn: func() interface{} { return [2]interface{}{"value", "v2"} }},
			{title: "from iface slice",
				valueFn: func() interface{} { return []interface{}{"value", "v2"} }},
		},
		unmarshalResults: []unmarshalResults{
			{title: "into string",
				slotFn:    func() interface{} { var str string; return str },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf("")}},
			{title: "into *string",
				slotFn:    func() interface{} { var str string; return &str },
				expectErr: ErrUnmarshalIncongruent{Token{Type: TArrOpen, Length: 2}, reflect.ValueOf("")}},
			{title: "into wildcard",
				slotFn:    func() interface{} { var v interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(interface{}(nil))}},
			{title: "into *wildcard",
				slotFn:  func() interface{} { var v interface{}; return &v },
				valueFn: func() interface{} { return []interface{}{"value", "v2"} }},
			{title: "into map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return v },
				expectErr: skipMe},
			{title: "into made map[str]iface",
				slotFn:    func() interface{} { v := make(map[string]interface{}); return v },
				expectErr: skipMe},
			{title: "into *map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return &v },
				expectErr: skipMe},
			{title: "into *map[str]str",
				slotFn:    func() interface{} { var v map[string]string; return &v },
				expectErr: skipMe},
			{title: "into []iface",
				slotFn:    func() interface{} { var v []interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf([]interface{}{})}},
			{title: "into *[]iface",
				slotFn:  func() interface{} { var v []interface{}; return &v },
				valueFn: func() interface{} { return []interface{}{"value", "v2"} }},
			{title: "into []str",
				slotFn:    func() interface{} { var v []string; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf([]string{})}},
			{title: "into *[]str",
				slotFn:  func() interface{} { var v []string; return &v },
				valueFn: func() interface{} { return []string{"value", "v2"} }},
			{title: "into []int",
				slotFn:    func() interface{} { var v []int; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf([]int{})}},
			{title: "into *[]int",
				slotFn:    func() interface{} { var v []int; return &v },
				expectErr: ErrUnmarshalIncongruent{Token{Type: TString, Str: "value"}, reflect.ValueOf(0)}},
		},
	},
	{title: "array nest in map pt1",
		sequence: fixtures.SequenceMap["array nested in map as non-first and final entry"],
		marshalResults: []marshalResults{
			{title: "from map[str]iface with nested []iface",
				valueFn: func() interface{} {
					return map[string]interface{}{"k1": "v1", "ke": []interface{}{"oh", "whee", "wow"}}
				}},
			{title: "from map[str]iface with nested []str",
				valueFn: func() interface{} {
					return map[string]interface{}{"k1": "v1", "ke": []string{"oh", "whee", "wow"}}
				}},
		},
		unmarshalResults: []unmarshalResults{
			{title: "into string",
				slotFn:    func() interface{} { var str string; return str },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf("")}},
			{title: "into *string",
				slotFn:    func() interface{} { var str string; return &str },
				expectErr: ErrUnmarshalIncongruent{Token{Type: TMapOpen, Length: 2}, reflect.ValueOf("")}},
			{title: "into wildcard",
				slotFn:    func() interface{} { var v interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(interface{}(nil))}},
			{title: "into *wildcard",
				slotFn: func() interface{} { var v interface{}; return &v },
				valueFn: func() interface{} {
					return map[string]interface{}{"k1": "v1", "ke": []interface{}{"oh", "whee", "wow"}}
				}},
			{title: "into map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(map[string]interface{}(nil))}},
			{title: "into made map[str]iface",
				slotFn: func() interface{} { v := make(map[string]interface{}); return v },
				valueFn: func() interface{} {
					return map[string]interface{}{"k1": "v1", "ke": []interface{}{"oh", "whee", "wow"}}
				}},
			{title: "into *map[str]iface",
				slotFn: func() interface{} { var v map[string]interface{}; return &v },
				valueFn: func() interface{} {
					return map[string]interface{}{"k1": "v1", "ke": []interface{}{"oh", "whee", "wow"}}
				}},
			{title: "into *map[str]str",
				slotFn: func() interface{} { var v map[string]string; return &v },
				//expectErr: ErrUnmarshalIncongruent{Token{Type: TArrOpen, Length: 3}, reflect.ValueOf("")}},
				expectErr: skipMe}, // big tricky todo: currently falls in the cracks where reflect core panics.
			{title: "into []iface",
				slotFn:    func() interface{} { var v []interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf([]interface{}{})}},
			{title: "into *[]iface",
				slotFn:    func() interface{} { var v []interface{}; return &v },
				expectErr: skipMe}, // should certainly error, but not well spec'd yet
			{title: "into []str",
				slotFn:    func() interface{} { var v []string; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf([]string{})}},
			{title: "into *[]str",
				slotFn:    func() interface{} { var v []string; return &v },
				expectErr: skipMe}, // should certainly error, but not well spec'd yet
			{title: "into []int",
				slotFn:    func() interface{} { var v []int; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf([]int{})}},
			{title: "into *[]int",
				slotFn:    func() interface{} { var v []int; return &v },
				expectErr: skipMe}, // should certainly error, but not well spec'd yet
		},
	},
	{title: "nested maps and arrays with no wildcards",
		sequence: fixtures.SequenceMap["map[str][]map[str]int"],
		marshalResults: []marshalResults{
			{title: "from oh-so-much type info",
				valueFn: func() interface{} {
					return map[string][]map[string]int{"k": []map[string]int{
						map[string]int{"k2": 1},
						map[string]int{"k2": 2},
					}}
				}},
		},
		unmarshalResults: []unmarshalResults{
			{title: "into oh-so-much type info",
				slotFn: func() interface{} { return make(map[string][]map[string]int) },
				valueFn: func() interface{} {
					return map[string][]map[string]int{"k": []map[string]int{
						map[string]int{"k2": 1},
						map[string]int{"k2": 2},
					}}
				}},
		},
	},
	{title: "nested maps and arrays as atlased structs",
		sequence: fixtures.SequenceMap["map[str][]map[str]int"],
		atlas: atlas.MustBuild(
			atlas.BuildEntry(tObjK{}).StructMap().
				AddField("K", atlas.StructMapEntry{SerialName: "k"}).
				Complete(),
			atlas.BuildEntry(tObjK2{}).StructMap().
				AddField("K2", atlas.StructMapEntry{SerialName: "k2"}).
				Complete(),
		),
		marshalResults: []marshalResults{
			{title: "from tObjK{[]tObjK2{}}",
				valueFn: func() interface{} {
					return tObjK{[]tObjK2{{1}, {2}}}
				}},
		},
		unmarshalResults: []unmarshalResults{
			{title: "into tObjK{[]tObjK2{}}",
				slotFn: func() interface{} { return &tObjK{} },
				valueFn: func() interface{} {
					return tObjK{[]tObjK2{{1}, {2}}}
				}},
		},
	},
	{title: "transform funks (struct<->string)",
		sequence: fixtures.SequenceMap["flat string"],
		atlas: atlas.MustBuild(
			atlas.BuildEntry(tObjStr{}).Transform().
				TransformMarshal(atlas.MakeMarshalTransformFunc(
					func(x tObjStr) (string, error) {
						return x.X, nil
					})).
				TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(
					func(x string) (tObjStr, error) {
						return tObjStr{x}, nil
					})).
				Complete(),
		),
		marshalResults: []marshalResults{
			{title: "from tObjStr",
				valueFn: func() interface{} { return tObjStr{"value"} }},
			{title: "from *tObjStr",
				valueFn: func() interface{} { return &tObjStr{"value"} }},
		},
		unmarshalResults: []unmarshalResults{
			{title: "into *tObjStr",
				slotFn:  func() interface{} { return &tObjStr{} },
				valueFn: func() interface{} { return tObjStr{"value"} }},
			// There are no tests here for "into interface{}" because by definition
			//  those situations wouldn't provide type info that would trigger these paths.
		},
	},
	{title: "typedef strings",
		sequence: fixtures.SequenceMap["flat string"],
		// no atlas necessary: the default behavior for a kind, even if typedef'd,
		//  is to simply to the natural thing for that kind.
		marshalResults: []marshalResults{
			{title: "from tDefStr literal",
				valueFn: func() interface{} { str := tDefStr("value"); return str }},
			{title: "from *tDefStr",
				valueFn: func() interface{} { str := tDefStr("value"); return &str }},
			{title: "from tDefStr in iface slot",
				valueFn: func() interface{} { var iface interface{}; iface = tDefStr("value"); return iface }},
			{title: "from tDefStr in *iface slot",
				valueFn: func() interface{} { var iface interface{}; iface = tDefStr("value"); return &iface }},
			{title: "from *tDefStr in iface slot",
				valueFn: func() interface{} { str := tDefStr("value"); var iface interface{}; iface = &str; return iface }},
			{title: "from *tDefStr in *iface slot",
				valueFn: func() interface{} { str := tDefStr("value"); var iface interface{}; iface = &str; return &iface }},
		},
		unmarshalResults: []unmarshalResults{
			{title: "into tDefStr",
				slotFn:    func() interface{} { var str tDefStr; return str },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(tDefStr(""))}},
			{title: "into *tDefStr",
				slotFn:  func() interface{} { var str tDefStr; return &str },
				valueFn: func() interface{} { str := tDefStr("value"); return str }},
		},
	},
	{title: "typedef ints",
		sequence: fixtures.SequenceMap["integer one"],
		// no atlas necessary: the default behavior for a kind, even if typedef'd,
		//  is to simply to the natural thing for that kind.
		marshalResults: []marshalResults{
			{title: "from tDefInt literal",
				valueFn: func() interface{} { v := tDefInt(1); return v }},
			{title: "from *tDefInt",
				valueFn: func() interface{} { v := tDefInt(1); return &v }},
			{title: "from tDefInt in iface slot",
				valueFn: func() interface{} { var iface interface{}; iface = tDefInt(1); return iface }},
			{title: "from tDefInt in *iface slot",
				valueFn: func() interface{} { var iface interface{}; iface = tDefInt(1); return &iface }},
			{title: "from *tDefInt in iface slot",
				valueFn: func() interface{} { v := tDefInt(1); var iface interface{}; iface = &v; return iface }},
			{title: "from *tDefInt in *iface slot",
				valueFn: func() interface{} { v := tDefInt(1); var iface interface{}; iface = &v; return &iface }},
		},
		unmarshalResults: []unmarshalResults{
			{title: "into tDefInt",
				slotFn:    func() interface{} { var v tDefInt; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(tDefInt(0))}},
			{title: "into *tDefInt",
				slotFn:  func() interface{} { var v tDefInt; return &v },
				valueFn: func() interface{} { v := tDefInt(1); return v }},
		},
	},
	{title: "typedef bytes",
		sequence: fixtures.SequenceMap["short byte array"],
		// no atlas necessary: the default behavior for a kind, even if typedef'd,
		//  is to simply to the natural thing for that kind.
		marshalResults: []marshalResults{
			{title: "from tDefBytes literal",
				valueFn: func() interface{} { v := tDefBytes(`value`); return v }},
			{title: "from *tDefBytes",
				valueFn: func() interface{} { v := tDefBytes(`value`); return &v }},
			{title: "from tDefBytes in iface slot",
				valueFn: func() interface{} { var iface interface{}; iface = tDefBytes(`value`); return iface }},
			{title: "from tDefBytes in *iface slot",
				valueFn: func() interface{} { var iface interface{}; iface = tDefBytes(`value`); return &iface }},
			{title: "from *tDefBytes in iface slot",
				valueFn: func() interface{} { v := tDefBytes(`value`); var iface interface{}; iface = &v; return iface }},
			{title: "from *tDefBytes in *iface slot",
				valueFn: func() interface{} { v := tDefBytes(`value`); var iface interface{}; iface = &v; return &iface }},
		},
		unmarshalResults: []unmarshalResults{
			{title: "into tDefBytes",
				slotFn:    func() interface{} { var v tDefBytes; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(tDefBytes{})}},
			{title: "into *tDefBytes",
				slotFn:  func() interface{} { var v tDefBytes; return &v },
				valueFn: func() interface{} { v := tDefBytes(`value`); return v }},
		},
	},
}

func TestMarshaller(t *testing.T) {
	// Package all the values from one step into a struct, just so that
	// we can assert on them all at once and make one green checkmark render per step.
	// Stringify the token first so extraneous fields in the union are hidden.
	type step struct {
		tok string
		err error
	}

	Convey("Marshaller suite:", t, func() {
		for _, tr := range objFixtures {
			Convey(fmt.Sprintf("%q fixture sequence:", tr.title), func() {
				for _, trr := range tr.marshalResults {
					maybe := Convey
					if trr.expectErr == skipMe {
						maybe = SkipConvey
					}
					// Conjure value.  Also format title for test, using its type info.
					value := trr.valueFn()
					valueKind := reflect.ValueOf(value).Kind()
					maybe(fmt.Sprintf("working %s (%s|%T):", trr.title, valueKind, value), func() {
						// Set up marshaller.
						marshaller := NewMarshaler(tr.atlas)
						marshaller.Bind(value)

						Convey("Steps...", func() {
							// Run steps until the marshaller says done or error.
							// For each step, assert the token matches fixtures;
							// when error and expected one, skip token check on that step
							// and finalize with the assertion.
							// If marshaller doesn't stop when we expected it to based
							// on fixture length, let it keep running three more steps
							// so we get that much more debug info.
							var done bool
							var err error
							var tok Token
							expectSteps := len(tr.sequence.Tokens) - 1
							for nStep := 0; nStep < expectSteps+3; nStep++ {
								done, err = marshaller.Step(&tok)
								if err != nil && trr.expectErr != nil {
									Convey("Result (error expected)", func() {
										So(err.Error(), ShouldResemble, trr.expectErr.Error())
									})
									return
								}
								if nStep <= expectSteps {
									So(
										step{tok.String(), err},
										ShouldResemble,
										step{tr.sequence.Tokens[nStep].String(), nil},
									)
								} else {
									So(
										step{tok.String(), err},
										ShouldResemble,
										step{Token{}.String(), fmt.Errorf("overshoot")},
									)
								}
								if done {
									Convey("Result (halted correctly)", func() {
										So(nStep, ShouldEqual, expectSteps)
									})
									return
								}
							}
						})
					})
				}
			})
		}
	})
}

func TestUnmarshaller(t *testing.T) {
	// Package all the values from one step into a struct, just so that
	// we can assert on them all at once and make one green checkmark render per step.
	// Stringify the token first so extraneous fields in the union are hidden.
	type step struct {
		tok  string
		err  error
		done bool
	}

	Convey("Unmarshaller suite:", t, func() {
		for _, tr := range objFixtures {
			Convey(fmt.Sprintf("%q fixture sequence:", tr.title), func() {
				for _, trr := range tr.unmarshalResults {
					maybe := Convey
					if trr.expectErr == skipMe {
						maybe = SkipConvey
					}
					// Conjure slot.  Also format title for test, using its type info.
					slot := trr.slotFn()
					slotKind := reflect.ValueOf(slot).Kind()
					maybe(fmt.Sprintf("targetting %s (%s|%T):", trr.title, slotKind, slot), func() {

						// Set up unmarshaller.
						unmarshaller := NewUnmarshaler(tr.atlas)
						err := unmarshaller.Bind(slot)
						if err != nil && trr.expectErr != nil {
							Convey("Result (error expected)", func() {
								So(err.Error(), ShouldResemble, trr.expectErr.Error())
							})
							return
						}

						Convey("Steps...", func() {
							// Run steps.
							// This is less complicated than the marshaller test
							// because we know exactly when we'll run out of them.
							var done bool
							var err error
							expectSteps := len(tr.sequence.Tokens) - 1
							for nStep, tok := range tr.sequence.Tokens {
								done, err = unmarshaller.Step(&tok)
								if err != nil && trr.expectErr != nil {
									Convey("Result (error expected)", func() {
										So(err.Error(), ShouldResemble, trr.expectErr.Error())
									})
									return
								}
								if nStep == expectSteps {
									So(
										step{tok.String(), err, done},
										ShouldResemble,
										step{tr.sequence.Tokens[nStep].String(), nil, true},
									)
								} else {
									So(
										step{tok.String(), err, done},
										ShouldResemble,
										step{tr.sequence.Tokens[nStep].String(), nil, false},
									)
								}
							}

							Convey("Result", func() {
								// Get value back out.  Some reflection required to get around pointers.
								rv := reflect.ValueOf(slot)
								if rv.Kind() == reflect.Ptr {
									rv = rv.Elem()
								}
								So(rv.Interface(), ShouldResemble, trr.valueFn())
							})
						})
					})
				}
			})
		}
	})
}
