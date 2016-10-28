package json

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/polydawn/go-xlate/dest/tok"
)

func TestParsing(t *testing.T) {
	tt := []struct {
		title  string
		input  string
		expect tok.Tokens
	}{
		{
			"flat string",
			`"value"`,
			tok.Tokens{
				{tok.TokenKind_ValString, "value"},
			},
		},
		//		{
		//			"single row map",
		//			func(sink xlate.Destination) {
		//				sink.OpenMap()
		//				sink.WriteMapKey("key")
		//				sink.WriteString("value")
		//				sink.CloseMap()
		//			},
		//			`{"key":"value"}`,
		//		},
		//		{
		//			"duo row map",
		//			func(sink xlate.Destination) {
		//				sink.OpenMap()
		//				sink.WriteMapKey("key")
		//				sink.WriteString("value")
		//				sink.WriteMapKey("k2")
		//				sink.WriteString("v2")
		//				sink.CloseMap()
		//			},
		//			`{"key":"value","k2":"v2"}`,
		//		},
		//		{
		//			"single entry array",
		//			func(sink xlate.Destination) {
		//				sink.OpenArray()
		//				sink.WriteString("value")
		//				sink.CloseArray()
		//			},
		//			`["value"]`,
		//		},
		//		{
		//			"duo entry array",
		//			func(sink xlate.Destination) {
		//				sink.OpenArray()
		//				sink.WriteString("value")
		//				sink.WriteString("v2")
		//				sink.CloseArray()
		//			},
		//			`["value","v2"]`,
		//		},
		//		{
		//			"empty map",
		//			func(sink xlate.Destination) {
		//				sink.OpenMap()
		//				sink.CloseMap()
		//			},
		//			`{}`,
		//		},
		//		{
		//			"empty array",
		//			func(sink xlate.Destination) {
		//				sink.OpenArray()
		//				sink.CloseArray()
		//			},
		//			`[]`,
		//		},
		//		{
		//			"array nested in map as non-first and final entry",
		//			func(sink xlate.Destination) {
		//				sink.OpenMap()
		//				sink.WriteMapKey("k1")
		//				sink.WriteString("v1")
		//				sink.WriteMapKey("ke")
		//				sink.OpenArray()
		//				sink.WriteString("oh")
		//				sink.WriteString("whee")
		//				sink.WriteString("wow")
		//				sink.CloseArray()
		//				sink.CloseMap()
		//			},
		//			`{"k1":"v1","ke":["oh","whee","wow"]}`,
		//		},
		//		{
		//			"array nested in map as first and non-final entry",
		//			func(sink xlate.Destination) {
		//				sink.OpenMap()
		//				sink.WriteMapKey("ke")
		//				sink.OpenArray()
		//				sink.WriteString("oh")
		//				sink.WriteString("whee")
		//				sink.WriteString("wow")
		//				sink.CloseArray()
		//				sink.WriteMapKey("k1")
		//				sink.WriteString("v1")
		//				sink.CloseMap()
		//			},
		//			`{"ke":["oh","whee","wow"],"k1":"v1"}`,
		//		},
		//		{
		//			"arrays in arrays in arrays",
		//			func(sink xlate.Destination) {
		//				sink.OpenArray()
		//				sink.OpenArray()
		//				sink.OpenArray()
		//				sink.CloseArray()
		//				sink.CloseArray()
		//				sink.CloseArray()
		//			},
		//			`[[[]]]`,
		//		},
	}
	for _, tr := range tt {
		buf := bytes.NewBufferString(tr.input)
		tokDest := tok.NewDestination()
		dec := decodeState{
			bytes:     newReader(buf),
			tokenSink: tokDest,
		}
		dec.scan.reset()
		dec.emitValue()
		assert(t, tr.title, tr.expect, tokDest.Tokens)
	}
}

func stringyEquality(x, y interface{}) bool {
	return fmt.Sprintf("%#v", x) == fmt.Sprintf("%#v", y)
}

func assert(t *testing.T, title string, expect, actual interface{}) {
	if !stringyEquality(expect, actual) {
		t.Errorf("test %q FAILED:\n\texpected  %#v\n\tactual    %#v",
			title, expect, actual)
	}
}
