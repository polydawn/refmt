package atlas

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStructMapAutogen(t *testing.T) {
	type BB struct {
		Z string
	}
	type AA struct {
		X string
		Y BB
	}

	Convey("StructMap Autogen:", t, func() {
		Convey("atlas entry should build for fixture type AA", func() {
			entry := AutogenerateStructMapEntry(reflect.TypeOf(AA{}))

			So(len(entry.StructMap.Fields), ShouldEqual, 2)
			So(entry.StructMap.Fields[0].SerialName, ShouldEqual, "X")
			So(entry.StructMap.Fields[0].ReflectRoute, ShouldResemble, ReflectRoute{0})
			So(entry.StructMap.Fields[0].Type, ShouldEqual, reflect.TypeOf(""))
			So(entry.StructMap.Fields[0].OmitEmpty, ShouldEqual, false)
			So(entry.StructMap.Fields[1].SerialName, ShouldEqual, "Y")
			So(entry.StructMap.Fields[1].ReflectRoute, ShouldResemble, ReflectRoute{1})
			So(entry.StructMap.Fields[1].Type, ShouldEqual, reflect.TypeOf(BB{}))
			So(entry.StructMap.Fields[1].OmitEmpty, ShouldEqual, false)
		})
	})
}
