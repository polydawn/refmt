package atlas

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStructMapAutogen(t *testing.T) {

	Convey("StructMap Autogen:", t, func() {
		type BB struct {
			Z string
		}
		type AA struct {
			X string
			Y BB
		}
		Convey("for a type which references other types, but is flat", func() {
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

		type CC struct {
			A AA
			BB
		}
		Convey("for a type which has some embedded structs", func() {
			entry := AutogenerateStructMapEntry(reflect.TypeOf(CC{}))

			So(len(entry.StructMap.Fields), ShouldEqual, 2)
			So(entry.StructMap.Fields[0].SerialName, ShouldEqual, "A")
			So(entry.StructMap.Fields[0].ReflectRoute, ShouldResemble, ReflectRoute{0})
			So(entry.StructMap.Fields[0].Type, ShouldEqual, reflect.TypeOf(AA{}))
			So(entry.StructMap.Fields[0].OmitEmpty, ShouldEqual, false)
			So(entry.StructMap.Fields[1].SerialName, ShouldEqual, "Z") // dives straight through embed!
			So(entry.StructMap.Fields[1].ReflectRoute, ShouldResemble, ReflectRoute{1, 0})
			So(entry.StructMap.Fields[1].Type, ShouldEqual, reflect.TypeOf(""))
			So(entry.StructMap.Fields[1].OmitEmpty, ShouldEqual, false)
		})
	})
}
