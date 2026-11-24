package obj

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/polydawn/refmt/tok"
)

func TestUnmarshalIntegerRange(t *testing.T) {
	for _, tc := range []struct {
		token   Token
		initial any
		want    any // nil expects ErrUnmarshalTypeCantFit and the initial value.
	}{
		{Token{Type: TInt, Int: -128}, int8(7), int8(-128)},
		{Token{Type: TInt, Int: -129}, int8(7), nil},
		{Token{Type: TUint, Uint: 127}, int8(7), int8(127)},
		{Token{Type: TUint, Uint: 128}, int8(7), nil},
		{Token{Type: TInt, Int: -1}, uint8(7), nil},
		{Token{Type: TInt, Int: 255}, uint8(7), uint8(255)},
		{Token{Type: TInt, Int: 256}, uint8(7), nil},
		{Token{Type: TUint, Uint: 256}, uint8(7), nil},
		{Token{Type: TUint, Uint: math.MaxInt64}, int64(7), int64(math.MaxInt64)},
		{Token{Type: TUint, Uint: 1 << 63}, int64(7), nil},
		{Token{Type: TUint, Uint: 1 << 63}, uint64(7), uint64(1 << 63)},
	} {
		slot := reflect.New(reflect.TypeOf(tc.initial))
		slot.Elem().Set(reflect.ValueOf(tc.initial))
		u := NewUnmarshaller(atlas.MustBuild())
		if err := u.Bind(slot.Interface()); err != nil {
			t.Fatal(err)
		}
		_, err := u.Step(&tc.token)
		want := tc.want
		if want == nil {
			var cantFit ErrUnmarshalTypeCantFit
			if !errors.As(err, &cantFit) {
				t.Errorf("%v into %T: got %v; want ErrUnmarshalTypeCantFit", tc.token, tc.initial, err)
			}
			want = tc.initial
		} else if err != nil {
			t.Errorf("%v into %T: %v", tc.token, tc.initial, err)
		}
		if got := slot.Elem().Interface(); got != want {
			t.Errorf("%v into %T: got %v; want %v", tc.token, tc.initial, got, want)
		}
	}
}
