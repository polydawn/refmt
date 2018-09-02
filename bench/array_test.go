package bench

import (
	"bytes"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/cbor"
	"github.com/polydawn/refmt/json"
)

var fixture_arrayFlatInt = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
var fixture_arrayFlatInt_json = []byte(`[1,2,3,4,5,6,7,8,9,0]`)
var fixture_arrayFlatInt_cbor = []byte{0x80 + 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

func Benchmark_ArrayFlatInt_MarshalToCborRefmt(b *testing.B) {
	var buf bytes.Buffer
	exerciseMarshaller(b,
		refmt.NewMarshaller(cbor.EncodeOptions{}, &buf), &buf,
		fixture_arrayFlatInt, fixture_arrayFlatInt_cbor,
	)
}

func Benchmark_ArrayFlatInt_MarshalToJsonRefmt(b *testing.B) {
	var buf bytes.Buffer
	exerciseMarshaller(b,
		refmt.NewMarshaller(json.EncodeOptions{}, &buf), &buf,
		fixture_arrayFlatInt, fixture_arrayFlatInt_json,
	)
}

func Benchmark_ArrayFlatInt_MarshalToJsonStdlib(b *testing.B) {
	exerciseStdlibJsonMarshaller(b,
		fixture_arrayFlatInt, fixture_arrayFlatInt_json,
	)
}

var fixture_arrayFlatStr = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
var fixture_arrayFlatStr_json = []byte(`["1","2","3","4","5","6","7","8","9","0"]`)
var fixture_arrayFlatStr_cbor = []byte{0x80 + 10, 0x60 + 1, 0x30 + 1, 0x60 + 1, 0x30 + 2, 0x60 + 1, 0x30 + 3, 0x60 + 1, 0x30 + 4, 0x60 + 1, 0x30 + 5, 0x60 + 1, 0x30 + 6, 0x60 + 1, 0x30 + 7, 0x60 + 1, 0x30 + 8, 0x60 + 1, 0x30 + 9, 0x60 + 1, 0x30 + 0}

func Benchmark_ArrayFlatStr_MarshalToCborRefmt(b *testing.B) {
	var buf bytes.Buffer
	exerciseMarshaller(b,
		refmt.NewMarshaller(cbor.EncodeOptions{}, &buf), &buf,
		fixture_arrayFlatStr, fixture_arrayFlatStr_cbor,
	)
}

func Benchmark_ArrayFlatStr_MarshalToJsonRefmt(b *testing.B) {
	var buf bytes.Buffer
	exerciseMarshaller(b,
		refmt.NewMarshaller(json.EncodeOptions{}, &buf), &buf,
		fixture_arrayFlatStr, fixture_arrayFlatStr_json,
	)
}

func Benchmark_ArrayFlatStr_MarshalToJsonStdlib(b *testing.B) {
	exerciseStdlibJsonMarshaller(b,
		fixture_arrayFlatStr, fixture_arrayFlatStr_json,
	)
}
