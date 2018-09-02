package bench

import (
	"bytes"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/cbor"
	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj/atlas"
)

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

var fixture_structAlpha = structAlpha{
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
var fixture_structAlpha_json = []byte(`{"B":{"R":{"R":{"R":{"R":null,"M":""},"M":"asdf"},"M":"quir"}},"C":{"N":"n","M":13},"C2":{"N":"n2","M":14},"X":1,"Y":2,"Z":"3","W":"4"}`)
var fixture_structAlpha_cbor = []byte{0xa7, 0x61, 0x42, 0xa1, 0x61, 0x52, 0xa2, 0x61, 0x52, 0xa2, 0x61, 0x52, 0xa2, 0x61, 0x52, 0xff, 0x61, 0x4d, 0x60, 0x61, 0x4d, 0x64, 0x61, 0x73, 0x64, 0x66, 0x61, 0x4d, 0x64, 0x71, 0x75, 0x69, 0x72, 0x61, 0x43, 0xa2, 0x61, 0x4e, 0x61, 0x6e, 0x61, 0x4d, 0x0d, 0x62, 0x43, 0x32, 0xa2, 0x61, 0x4e, 0x62, 0x6e, 0x32, 0x61, 0x4d, 0x0e, 0x61, 0x58, 0x01, 0x61, 0x59, 0x02, 0x61, 0x5a, 0x61, 0x33, 0x61, 0x57, 0x61, 0x34}
var fixture_structAlpha_atlas = atlas.MustBuild(
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

func Benchmark_StructAlpha_MarshalToCborRefmt(b *testing.B) {
	var buf bytes.Buffer
	exerciseMarshaller(b,
		refmt.NewMarshaller(cbor.EncodeOptions{}, &buf), &buf,
		fixture_arrayFlatInt, fixture_arrayFlatInt_cbor,
	)
}

func Benchmark_StructAlpha_MarshalToJsonRefmt(b *testing.B) {
	var buf bytes.Buffer
	exerciseMarshaller(b,
		refmt.NewMarshaller(json.EncodeOptions{}, &buf), &buf,
		fixture_arrayFlatInt, fixture_arrayFlatInt_json,
	)
}

func Benchmark_StructAlpha_MarshalToJsonStdlib(b *testing.B) {
	exerciseStdlibJsonMarshaller(b,
		fixture_arrayFlatInt, fixture_arrayFlatInt_json,
	)
}
