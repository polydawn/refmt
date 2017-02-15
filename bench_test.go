package xlate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	xlatejson "github.com/polydawn/go-xlate/json"
	"github.com/polydawn/go-xlate/obj"
	"github.com/polydawn/go-xlate/obj/atlas"
)

//
// slice of ints test
//

var fixture_arrayFlatInt = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
var fixture_arrayFlatInt_expect = []byte(`[1,2,3,4,5,6,7,8,9,0]`)

func Benchmark_ArrayFlatIntToJson_Xlate(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := NewJsonEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Marshal(&fixture_arrayFlatInt)
	}
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(buf.Bytes(), fixture_arrayFlatInt_expect) {
		panic(fmt.Errorf("result \"%s\"\nmust equal \"%s\"", buf.Bytes(), fixture_arrayFlatInt_expect))
	}
}
func Benchmark_ArrayFlatIntToJson_Stdlib(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := json.NewEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Encode(&fixture_arrayFlatInt)
	}
	if err != nil {
		panic(err)
	}
	buf.Truncate(buf.Len() - 1) // Stdlib suffixes a linebreak.
	if !bytes.Equal(buf.Bytes(), fixture_arrayFlatInt_expect) {
		panic(fmt.Errorf("result \"%s\"\nmust equal \"%s\"", buf.Bytes(), fixture_arrayFlatInt_expect))
	}
}

//
// slice of strings test
//

var fixture_arrayFlatStr = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
var fixture_arrayFlatStr_expect = []byte(`["1","2","3","4","5","6","7","8","9","0"]`)

func Benchmark_ArrayFlatStrToJson_Xlate(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := NewJsonEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Marshal(&fixture_arrayFlatStr)
	}
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(buf.Bytes(), fixture_arrayFlatStr_expect) {
		panic(fmt.Errorf("result \"%s\"\nmust equal \"%s\"", buf.Bytes(), fixture_arrayFlatStr_expect))
	}
}
func Benchmark_ArrayFlatStrToJson_Stdlib(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := json.NewEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Encode(&fixture_arrayFlatStr)
	}
	if err != nil {
		panic(err)
	}
	buf.Truncate(buf.Len() - 1) // Stdlib suffixes a linebreak.
	if !bytes.Equal(buf.Bytes(), fixture_arrayFlatStr_expect) {
		panic(fmt.Errorf("result \"%s\"\nmust equal \"%s\"", buf.Bytes(), fixture_arrayFlatStr_expect))
	}
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
var fixture_struct_expect = []byte(`{"B":{"R":{"R":{"R":{"R":null,"M":""},"M":"asdf"},"M":"quir"}},"C":{"N":"n","M":13},"C2":{"N":"n2","M":14},"X":1,"Y":2,"Z":"3","W":"4"}`)
var fixture_suiteFieldRoute = (&obj.Suite{}).
	Add(structAlpha{}, obj.MarshalMachineStructAtlasFactory(atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "B", FieldRoute: atlas.FieldRoute{0}},
			{Name: "C", FieldRoute: atlas.FieldRoute{1}},
			{Name: "C2", FieldRoute: atlas.FieldRoute{2}},
			{Name: "X", FieldRoute: atlas.FieldRoute{3}},
			{Name: "Y", FieldRoute: atlas.FieldRoute{4}},
			{Name: "Z", FieldRoute: atlas.FieldRoute{5}},
			{Name: "W", FieldRoute: atlas.FieldRoute{6}},
		},
	})).
	Add(structBeta{}, obj.MarshalMachineStructAtlasFactory(atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "R", FieldRoute: atlas.FieldRoute{0}},
		},
	})).
	Add(structGamma{}, obj.MarshalMachineStructAtlasFactory(atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "N", FieldRoute: atlas.FieldRoute{0}},
			{Name: "M", FieldRoute: atlas.FieldRoute{1}},
		},
	})).
	Add(structRecursive{}, obj.MarshalMachineStructAtlasFactory(atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "R", FieldRoute: atlas.FieldRoute{0}},
			{Name: "M", FieldRoute: atlas.FieldRoute{1}},
		},
	}))
var fixture_suiteAddrFunc = (&obj.Suite{}).
	Add(structAlpha{}, obj.MarshalMachineStructAtlasFactory(atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "B", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).B) }},
			{Name: "C", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).C) }},
			{Name: "C2", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).C2) }},
			{Name: "X", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).X) }},
			{Name: "Y", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).Y) }},
			{Name: "Z", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).Z) }},
			{Name: "W", AddrFunc: func(v interface{}) interface{} { return &(v.(*structAlpha).W) }},
		},
	})).
	Add(structBeta{}, obj.MarshalMachineStructAtlasFactory(atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "R", AddrFunc: func(v interface{}) interface{} { return &(v.(*structBeta).R) }},
		},
	})).
	Add(structGamma{}, obj.MarshalMachineStructAtlasFactory(atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "N", AddrFunc: func(v interface{}) interface{} { return &(v.(*structGamma).N) }},
			{Name: "M", AddrFunc: func(v interface{}) interface{} { return &(v.(*structGamma).M) }},
		},
	})).
	Add(structRecursive{}, obj.MarshalMachineStructAtlasFactory(atlas.Atlas{
		Fields: []atlas.Entry{
			{Name: "R", AddrFunc: func(v interface{}) interface{} { return &(v.(*structRecursive).R) }},
			{Name: "M", AddrFunc: func(v interface{}) interface{} { return &(v.(*structRecursive).M) }},
		},
	}))

func Benchmark_StructToJson_XlateFieldRoute(b *testing.B) {
	var buf bytes.Buffer
	var err error
	marshaller := obj.NewMarshaler(fixture_suiteFieldRoute)
	serializer := xlatejson.NewSerializer(&buf)
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
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(buf.Bytes(), fixture_struct_expect) {
		panic(fmt.Errorf("result \"%s\"\nmust equal \"%s\"", buf.Bytes(), fixture_struct_expect))
	}
}
func Benchmark_StructToJson_XlateAddrFunc(b *testing.B) {
	var buf bytes.Buffer
	var err error
	marshaller := obj.NewMarshaler(fixture_suiteAddrFunc)
	serializer := xlatejson.NewSerializer(&buf)
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
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(buf.Bytes(), fixture_struct_expect) {
		panic(fmt.Errorf("result \"%s\"\nmust equal \"%s\"", buf.Bytes(), fixture_struct_expect))
	}
}
func Benchmark_StructToJson_Stdlib(b *testing.B) {
	var buf bytes.Buffer
	var err error
	enc := json.NewEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err = enc.Encode(&fixture_struct)
	}
	if err != nil {
		panic(err)
	}
	buf.Truncate(buf.Len() - 1) // Stdlib suffixes a linebreak.
	if !bytes.Equal(buf.Bytes(), fixture_struct_expect) {
		panic(fmt.Errorf("result \"%s\"\nmust equal \"%s\"", buf.Bytes(), fixture_struct_expect))
	}
}
