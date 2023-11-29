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
				{"trailing commas", `{"x":"1","y":"2",,,}`, "invalid char while expecting start of key: comma"},
				{"just commas", `{,,,}`, "invalid char while expecting start of key: comma"},
				{"leading commas", `{,,,"x":"1","y":"2",,,}`, "invalid char while expecting start of key: comma"},
				{"no commas", `{"x":"1""y":"2"}`, "expected comma or map close after map value; got quote"},
				{"no commas, just spaces", `{    "x":"1"    "y":"2"   }`, "expected comma or map close after map value; got quote"},
				{"no commas, just tabs", `{	"x":"1"	"y":"2"	}`, "expected comma or map close after map value; got quote"},
			} {
				Convey(tc.name, func() {
					var slot map[string]string
					err := Unmarshal(json.DecodeOptions{}, []byte(tc.j), &slot)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, tc.err)
				})
			}
		})
		Convey("array comma handling errors", func() {
			for _, tc := range []struct {
				name string
				j    string
				err  string
			}{
				{"trailing commas", `["1","2",,,]`, "invalid char while expecting start of value: comma"},
				{"just commas", `[,,,]`, "invalid char while expecting start of value: comma"},
				{"leading commas", `[,,,"1","2",,,]`, "invalid char while expecting start of value: comma"},
				{"no commas", `["1""2"]`, "expected comma or array close after array value; got quote"},
				{"no commas, just spaces", `[    "1"    "2"   ]`, "expected comma or array close after array value; got quote"},
				{"no commas, just tabs", `[	"1"	"2"	]`, "expected comma or array close after array value; got quote"},
			} {
				Convey(tc.name, func() {
					var slot []string
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
