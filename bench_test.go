package refmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/polydawn/refmt/obj/atlas"
)

func checkAftermath(err error, result []byte, expect []byte) {
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(result, expect) {
		// fmt note: "space-x" is nice to read as hex; "%q" will try harder to print ascii, but often looks fairly bonkers anyway on e.g. cbor.
		panic(fmt.Errorf("result \"% x\"\nmust equal \"% x\"", result, expect))
	}
}

//
// slice of ints test
//

var fixture_arrayFlatInt = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
var fixture_arrayFlatInt_json = []byte(`[1,2,3,4,5,6,7,8,9,0]`)
var fixture_arrayFlatInt_cbor = []byte{0x80 + 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

func Benchmark_ArrayFlatIntToJson_Refmt(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := NewJsonEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Marshal(&fixture_arrayFlatInt)
	}
	checkAftermath(err, buf.Bytes(), fixture_arrayFlatInt_json)
}
func Benchmark_ArrayFlatIntToJson_Stdlib(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := json.NewEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Encode(&fixture_arrayFlatInt)
	}
	buf.Truncate(buf.Len() - 1) // Stdlib suffixes a linebreak.
	checkAftermath(err, buf.Bytes(), fixture_arrayFlatInt_json)
}
func Benchmark_ArrayFlatIntToCbor_Refmt(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := NewCborEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Marshal(&fixture_arrayFlatInt)
	}
	checkAftermath(err, buf.Bytes(), fixture_arrayFlatInt_cbor)
}

//
// slice of strings test
//

var fixture_arrayFlatStr = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
var fixture_arrayFlatStr_json = []byte(`["1","2","3","4","5","6","7","8","9","0"]`)
var fixture_arrayFlatStr_cbor = []byte{0x80 + 10, 0x60 + 1, 0x30 + 1, 0x60 + 1, 0x30 + 2, 0x60 + 1, 0x30 + 3, 0x60 + 1, 0x30 + 4, 0x60 + 1, 0x30 + 5, 0x60 + 1, 0x30 + 6, 0x60 + 1, 0x30 + 7, 0x60 + 1, 0x30 + 8, 0x60 + 1, 0x30 + 9, 0x60 + 1, 0x30 + 0}

func Benchmark_ArrayFlatStrToJson_Refmt(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := NewJsonEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Marshal(&fixture_arrayFlatStr)
	}
	checkAftermath(err, buf.Bytes(), fixture_arrayFlatStr_json)
}
func Benchmark_ArrayFlatStrToJson_Stdlib(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := json.NewEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Encode(&fixture_arrayFlatStr)
	}
	buf.Truncate(buf.Len() - 1) // Stdlib suffixes a linebreak.
	checkAftermath(err, buf.Bytes(), fixture_arrayFlatStr_json)
}
func Benchmark_ArrayFlatStrToCbor_Refmt(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := NewCborEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Marshal(&fixture_arrayFlatStr)
	}
	checkAftermath(err, buf.Bytes(), fixture_arrayFlatStr_cbor)
}

//
// object traversal tests
//

type structAlpha struct {
	B  *structBeta
	C  structGamma
	C2 structGamma
	X  int
	Y  int
	Z  string
	W  string
}
type structBeta struct {
	R *structRecursive
}
type structGamma struct {
	N string
	M int
}
type structRecursive struct {
	R *structRecursive
	M string
}

var fixture_struct = structAlpha{
	&structBeta{
		&structRecursive{
			&structRecursive{
				&structRecursive{},
				"asdf",
			},
			"quir",
		},
	},
	structGamma{"n", 13},
	structGamma{"n2", 14},
	1, 2, "3", "4",
}

// note: 18 string keys, 7 string values; total 25 strings.
var fixture_struct_json = []byte(`{"B":{"R":{"R":{"R":{"R":null,"M":""},"M":"asdf"},"M":"quir"}},"C":{"N":"n","M":13},"C2":{"N":"n2","M":14},"X":1,"Y":2,"Z":"3","W":"4"}`)
var fixture_struct_cbor = []byte{0xa7, 0x61, 0x42, 0xa1, 0x61, 0x52, 0xa2, 0x61, 0x52, 0xa2, 0x61, 0x52, 0xa2, 0x61, 0x52, 0xff, 0x61, 0x4d, 0x60, 0x61, 0x4d, 0x64, 0x61, 0x73, 0x64, 0x66, 0x61, 0x4d, 0x64, 0x71, 0x75, 0x69, 0x72, 0x61, 0x43, 0xa2, 0x61, 0x4e, 0x61, 0x6e, 0x61, 0x4d, 0x0d, 0x62, 0x43, 0x32, 0xa2, 0x61, 0x4e, 0x62, 0x6e, 0x32, 0x61, 0x4d, 0x0e, 0x61, 0x58, 0x01, 0x61, 0x59, 0x02, 0x61, 0x5a, 0x61, 0x33, 0x61, 0x57, 0x61, 0x34}
var fixture_struct_atlas = atlas.MustBuild(
	atlas.BuildEntry(structAlpha{}).StructMap().
		AddField("B", atlas.StructMapEntry{SerialName: "B"}).
		AddField("C", atlas.StructMapEntry{SerialName: "C"}).
		AddField("C2", atlas.StructMapEntry{SerialName: "C2"}).
		AddField("X", atlas.StructMapEntry{SerialName: "X"}).
		AddField("Y", atlas.StructMapEntry{SerialName: "Y"}).
		AddField("Z", atlas.StructMapEntry{SerialName: "Z"}).
		AddField("W", atlas.StructMapEntry{SerialName: "W"}).
		Complete(),
	atlas.BuildEntry(structBeta{}).StructMap().
		AddField("R", atlas.StructMapEntry{SerialName: "R"}).
		Complete(),
	atlas.BuildEntry(structGamma{}).StructMap().
		AddField("N", atlas.StructMapEntry{SerialName: "N"}).
		AddField("M", atlas.StructMapEntry{SerialName: "M"}).
		Complete(),
	atlas.BuildEntry(structRecursive{}).StructMap().
		AddField("R", atlas.StructMapEntry{SerialName: "R"}).
		AddField("M", atlas.StructMapEntry{SerialName: "M"}).
		Complete(),
)

func Benchmark_StructToJson_Refmt(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := NewAtlasedJsonEncoder(&buf, fixture_struct_atlas)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Marshal(&fixture_struct)
	}
	checkAftermath(err, buf.Bytes(), fixture_struct_json)
}
func Benchmark_StructToJson_Stdlib(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := json.NewEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Encode(&fixture_struct)
	}
	buf.Truncate(buf.Len() - 1) // Stdlib suffixes a linebreak.
	checkAftermath(err, buf.Bytes(), fixture_struct_json)
}
func Benchmark_StructToCbor_Refmt(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := NewAtlasedCborEncoder(&buf, fixture_struct_atlas)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Marshal(&fixture_struct)
	}
	checkAftermath(err, buf.Bytes(), fixture_struct_cbor)
}
