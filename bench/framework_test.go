package bench

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"testing"

	wish "github.com/warpfork/go-wish"
	"github.com/warpfork/go-wish/cmp"

	"github.com/polydawn/refmt"
)

func exerciseMarshaller(
	b *testing.B,
	subj refmt.Marshaller,
	buf *bytes.Buffer,
	val interface{},
	expect []byte,
) {
	var err error
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = subj.Marshal(val)
	}
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(buf.Bytes(), expect) {
		panic(fmt.Errorf("result \"% x\"\nmust equal \"% x\"", buf.Bytes(), expect))
	}
}

func exerciseStdlibJsonMarshaller(
	b *testing.B,
	val interface{},
	expect []byte,
) {
	var err error
	var buf bytes.Buffer
	subj := stdjson.NewEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = subj.Encode(val)
	}
	if err != nil {
		panic(err)
	}
	buf.Truncate(buf.Len() - 1) // Stdlib suffixes a linebreak.
	if !bytes.Equal(buf.Bytes(), expect) {
		panic(fmt.Errorf("result \"% x\"\nmust equal \"% x\"", buf.Bytes(), expect))
	}
}

func exerciseUnmarshaller(
	b *testing.B,
	subj refmt.Unmarshaller,
	buf *bytes.Buffer,
	src []byte,
	blankFn func() interface{},
	expect interface{},
) {
	var err error
	var targ interface{}
	for i := 0; i < b.N; i++ {
		targ = blankFn()
		buf.Reset()
		buf.Write(src)
		err = subj.Unmarshal(targ)
	}
	if err != nil {
		panic(err)
	}
	if detail, pass := wish.ShouldEqual(targ, expect); !pass {
		panic(fmt.Errorf("difference:\n%s", detail))
	}
}

func exerciseStdlibJsonUnmarshaller(
	b *testing.B,
	src []byte,
	blankFn func() interface{},
	expect interface{},
) {
	var err error
	var targ interface{}
	for i := 0; i < b.N; i++ {
		targ = blankFn()
		subj := stdjson.NewDecoder(bytes.NewBuffer(src))
		err = subj.Decode(targ)
	}
	if err != nil {
		panic(err)
	}
	if diff := cmp.Diff(fixFloatsToInts(targ), fixFloatsToInts(expect), cmpOpt_flattenFloats); diff != "" {
		panic(fmt.Errorf("difference:\n%s", diff))
	}
}

// This function normalizes floats to ints, and we use it so the same fixtures
// work for CBOR and refmt-JSON and stdlib-JSON -- the latter of which only
// produces floats when unmarshalling into an empty interface.
//
// The whole function is fairly absurd, but so is refusing to admit ints exist.
func fixFloatsToInts(in interface{}) interface{} {
	switch in2 := in.(type) {
	case *map[string]interface{}:
		return fixFloatsToInts(*in2)
	case map[string]interface{}:
		out := make(map[string]interface{}, len(in2))
		for k, v := range in2 {
			out[k] = fixFloatsToInts(v)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(in2))
		for i, v := range in2 {
			out[i] = fixFloatsToInts(v)
		}
		return out
	case float64:
		return int(in2)
	default:
		return in
	}
}

// This feature...... doesn't seem to... work.  I'm sure I'm holding it wrong,
// but I don't know how.
// I'll remove this in a moment, I just wanted a commit saying "I tried it".
// (Strip the 'fixFloatsToInts' uses, and flip targ and expect either way you
// like; I cannot seem to trigger the 'hello' panic. :( )
var cmpOpt_flattenFloats = cmp.Transformer(
	"flattenFloats",
	func(in float64) int { panic("hello?"); return int(in) },
)
