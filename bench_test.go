package xlate

import (
	"encoding/json"
	"testing"
)

var fixture_arrayFlatInt = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

func Benchmark_ArrayFlatIntToJson_Xlate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JsonEncode(&fixture_arrayFlatInt)
	}
}
func Benchmark_ArrayFlatIntToJson_Stdlib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(fixture_arrayFlatInt)
	}
}

var fixture_arrayFlatStr = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

func Benchmark_ArrayFlatStrToJson_Xlate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JsonEncode(&fixture_arrayFlatStr)
	}
}
func Benchmark_ArrayFlatStrToJson_Stdlib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(fixture_arrayFlatStr)
	}
}
