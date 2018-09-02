package bench

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"testing"

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
