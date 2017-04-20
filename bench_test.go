package refmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/polydawn/refmt/cbor"
	refmtjson "github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj"
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

func Benchmark_ArrayFlatIntToJson_Xlate(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := NewJsonEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Marshal(&fixture_arrayFlatInt)
	}
	checkAftermath(err, buf.Bytes(), fixture_arrayFlatInt_json)
}
func Benchmark_ArrayFlatIntToCbor_Xlate(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := NewCborEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Marshal(&fixture_arrayFlatInt)
	}
	checkAftermath(err, buf.Bytes(), fixture_arrayFlatInt_cbor)
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

//
// slice of strings test
//

var fixture_arrayFlatStr = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
var fixture_arrayFlatStr_json = []byte(`["1","2","3","4","5","6","7","8","9","0"]`)
var fixture_arrayFlatStr_cbor = []byte{0x80 + 10, 0x60 + 1, 0x30 + 1, 0x60 + 1, 0x30 + 2, 0x60 + 1, 0x30 + 3, 0x60 + 1, 0x30 + 4, 0x60 + 1, 0x30 + 5, 0x60 + 1, 0x30 + 6, 0x60 + 1, 0x30 + 7, 0x60 + 1, 0x30 + 8, 0x60 + 1, 0x30 + 9, 0x60 + 1, 0x30 + 0}

func Benchmark_ArrayFlatStrToJson_Xlate(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := NewJsonEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Marshal(&fixture_arrayFlatStr)
	}
	checkAftermath(err, buf.Bytes(), fixture_arrayFlatStr_json)
}
func Benchmark_ArrayFlatStrToCbor_Xlate(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := NewCborEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Marshal(&fixture_arrayFlatStr)
	}
	checkAftermath(err, buf.Bytes(), fixture_arrayFlatStr_cbor)
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
var fixture_suiteFieldRoute = (&obj.Suite{}).
	Add(structAlpha{}, obj.Morphism{Atlas: atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "B", FieldRoute: atlas.FieldRoute{0}},
			{Name: "C", FieldRoute: atlas.FieldRoute{1}},
			{Name: "C2", FieldRoute: atlas.FieldRoute{2}},
			{Name: "X", FieldRoute: atlas.FieldRoute{3}},
			{Name: "Y", FieldRoute: atlas.FieldRoute{4}},
			{Name: "Z", FieldRoute: atlas.FieldRoute{5}},
			{Name: "W", FieldRoute: atlas.FieldRoute{6}},
		},
	}}).
	Add(structBeta{}, obj.Morphism{Atlas: atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "R", FieldRoute: atlas.FieldRoute{0}},
		},
	}}).
	Add(structGamma{}, obj.Morphism{Atlas: atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "N", FieldRoute: atlas.FieldRoute{0}},
			{Name: "M", FieldRoute: atlas.FieldRoute{1}},
		},
	}}).
	Add(structRecursive{}, obj.Morphism{Atlas: atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "R", FieldRoute: atlas.FieldRoute{0}},
			{Name: "M", FieldRoute: atlas.FieldRoute{1}},
		},
	}})
var fixture_suiteAddrFunc = (&obj.Suite{}).
	Add(structAlpha{}, obj.Morphism{Atlas: atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "B", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).B) }},
			{Name: "C", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).C) }},
			{Name: "C2", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).C2) }},
			{Name: "X", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).X) }},
			{Name: "Y", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).Y) }},
			{Name: "Z", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).Z) }},
			{Name: "W", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).W) }},
		},
	}}).
	Add(structBeta{}, obj.Morphism{Atlas: atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "R", AddrFunc: func(v interface{}) interface{} { return &(v.(*structBeta).R) }},
		},
	}}).
	Add(structGamma{}, obj.Morphism{Atlas: atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "N", AddrFunc: func(v interface{}) interface{} { return &(v.(*structGamma).N) }},
			{Name: "M", AddrFunc: func(v interface{}) interface{} { return &(v.(*structGamma).M) }},
		},
	}}).
	Add(structRecursive{}, obj.Morphism{Atlas: atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "R", AddrFunc: func(v interface{}) interface{} { return &(v.(*structRecursive).R) }},
			{Name: "M", AddrFunc: func(v interface{}) interface{} { return &(v.(*structRecursive).M) }},
		},
	}})

func Benchmark_StructToJson_XlateFieldRoute(b *testing.B) {
	var buf bytes.Buffer
	var err error
	marshaller := obj.NewMarshaler(fixture_suiteFieldRoute)
	serializer := refmtjson.NewSerializer(&buf)
	enc := TokenPump{
		marshaller,
		serializer,
	}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		marshaller.Bind(&fixture_struct)
		serializer.Reset()
		err = enc.Run()
	}
	checkAftermath(err, buf.Bytes(), fixture_struct_json)
}
func Benchmark_StructToCbor_XlateFieldRoute(b *testing.B) {
	var buf bytes.Buffer
	var err error
	marshaller := obj.NewMarshaler(fixture_suiteFieldRoute)
	encoder := cbor.NewEncoder(&buf)
	serializer := TokenPump{
		marshaller,
		encoder,
	}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		marshaller.Bind(&fixture_struct)
		encoder.Reset()
		err = serializer.Run()
	}
	checkAftermath(err, buf.Bytes(), fixture_struct_cbor)
}
func Benchmark_StructToJson_XlateAddrFunc(b *testing.B) {
	var buf bytes.Buffer
	var err error
	marshaller := obj.NewMarshaler(fixture_suiteAddrFunc)
	serializer := refmtjson.NewSerializer(&buf)
	enc := TokenPump{
		marshaller,
		serializer,
	}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		marshaller.Bind(&fixture_struct)
		serializer.Reset()
		err = enc.Run()
	}
	checkAftermath(err, buf.Bytes(), fixture_struct_json)
}
func Benchmark_StructToCbor_XlateAddrFunc(b *testing.B) {
	var buf bytes.Buffer
	var err error
	marshaller := obj.NewMarshaler(fixture_suiteAddrFunc)
	encoder := cbor.NewEncoder(&buf)
	serializer := TokenPump{
		marshaller,
		encoder,
	}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		marshaller.Bind(&fixture_struct)
		encoder.Reset()
		err = serializer.Run()
	}
	checkAftermath(err, buf.Bytes(), fixture_struct_cbor)
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
