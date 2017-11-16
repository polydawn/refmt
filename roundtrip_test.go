package refmt_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/cbor"
	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj/atlas"
)

func TestRoundTrip(t *testing.T) {
	t.Run("nil nil", func(t *testing.T) {
		testRoundTripAllEncodings(t, nil, atlas.MustBuild())
	})
	t.Run("empty []interface{}", func(t *testing.T) {
		testRoundTripAllEncodings(t, []interface{}{}, atlas.MustBuild())
	})
	t.Run("nil []interface{}", func(t *testing.T) {
		testRoundTripAllEncodings(t, []interface{}(nil), atlas.MustBuild())
	})
	t.Run("empty map[string]interface{}", func(t *testing.T) {
		testRoundTripAllEncodings(t, map[string]interface{}(nil), atlas.MustBuild())
	})
	t.Run("nil map[string]interface{}", func(t *testing.T) {
		testRoundTripAllEncodings(t, map[string]interface{}(nil), atlas.MustBuild())
	})
}

func testRoundTripAllEncodings(
	t *testing.T,
	value interface{},
	atl atlas.Atlas,
) {
	t.Run("cbor", func(t *testing.T) {
		roundTrip(t, value, cbor.EncodeOptions{}, cbor.DecodeOptions{}, atl)
	})
	t.Run("json", func(t *testing.T) {
		roundTrip(t, value, json.EncodeOptions{}, json.DecodeOptions{}, atl)
	})
}

func roundTrip(
	t *testing.T,
	value interface{},
	encodeOptions refmt.EncodeOptions,
	decodeOptions refmt.DecodeOptions,
	atl atlas.Atlas,
) {
	// Encode.
	var buf bytes.Buffer
	encoder := refmt.NewMarshallerAtlased(encodeOptions, &buf, atl)
	if err := encoder.Marshal(value); err != nil {
		t.Fatalf("failed encoding: %s", err)
	}

	// Decode back to obj.
	decoder := refmt.NewUnmarshallerAtlased(decodeOptions, bytes.NewBuffer(buf.Bytes()), atl)
	var slot interface{}
	if err := decoder.Unmarshal(&slot); err != nil {
		t.Fatalf("failed decoding: %s", err)
	}
	t.Logf("%#T -- %#v", slot, slot)

	// Re-encode.  Expect to get same encoded form.
	var buf2 bytes.Buffer
	encoder2 := refmt.NewMarshallerAtlased(encodeOptions, &buf2, atl)
	if err := encoder2.Marshal(slot); err != nil {
		t.Fatalf("failed re-encoding: %s", err)
	}

	// Stringify.  (Plain "%q" escapes unprintables quite nicely.)
	str1 := fmt.Sprintf("%q", buf.String())
	str2 := fmt.Sprintf("%q", buf2.String())
	if str1 != str2 {
		t.Errorf("%q != %q", str1, str2)
	}
	t.Logf("%#v == %q", value, str1)
}
