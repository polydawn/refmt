package refmt

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj/atlas"
)

func TestUnmarshal(t *testing.T) {
	Convey("json", t, func() {
		Convey("string", func() {
			var slot string
			bs := []byte(`"str"`)
			err := Unmarshal(json.DecodeOptions{}, bs, &slot)
			So(err, ShouldBeNil)
			So(slot, ShouldEqual, "str")
		})
		Convey("map", func() {
			var slot map[string]string
			bs := []byte(`{"x":"1"}`)
			err := Unmarshal(json.DecodeOptions{}, bs, &slot)
			So(err, ShouldBeNil)
			So(slot, ShouldResemble, map[string]string{"x": "1"})
		})
		Convey("map comma handling", func() {
			var slot map[string]string
			bs := []byte(`{"x":"1","y":"2"}`)
			err := Unmarshal(json.DecodeOptions{}, bs, &slot)
			So(err, ShouldBeNil)
			So(slot, ShouldResemble, map[string]string{"x": "1", "y": "2"})
		})
		Convey("map comma handling errors", func() {
			for _, tc := range []struct {
				name string
				j    string
				err  string
			}{
				// fails with: "expected colon after map key; got 0x2c"
				{"trailing commas", `{"x":"1","y":"2",,,}`, "expected key after comma, got comma"},
				// fails with: "expected colon after map key; got 0x2c"
				{"just commas", `{,,,}`, "expected key after comma, got comma"},
				// fails with: "expected colon after map key; got 0x2c"
				{"leading commas", `{,,,"x":"1","y":"2",,,}`, "expected key after map open, got comma"},
				// doesn't error
				{"no commas", `{"x":"1""y":"2"}`, "expected comma after value, got quote"},
				// doesn't error
				{"no commas, just spaces", `{    "x":"1"    "y":"2"   }`, "expected comma after value, got quote"},
				// doesn't error
				{"no commas, just tabs", `{	"x":"1"	"y":"2"	}`, "expected comma after value, got quote"},
			} {
				Convey(tc.name, func() {
					var slot map[string]string
					err := Unmarshal(json.DecodeOptions{}, []byte(tc.j), &slot)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, tc.err)
				})
			}
		})
	})
}

func TestUnmarshalAtlased(t *testing.T) {
	Convey("json", t, func() {
		Convey("obj", func() {
			type testObj struct {
				X string
				Y string
			}
			var slot testObj
			atl := atlas.MustBuild(
				atlas.BuildEntry(testObj{}).
					StructMap().Autogenerate().
					Complete(),
			)
			bs := []byte(`{"x":"1","y":"2"}`)
			err := UnmarshalAtlased(json.DecodeOptions{}, bs, &slot, atl)
			So(err, ShouldBeNil)
		})
	})
}
