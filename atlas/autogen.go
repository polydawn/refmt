package atlas

import (
	"reflect"

	"github.com/polydawn/go-xlate"
)

/*
	Produces a MappingFunc which will transform inputs of this type to tokens
	much as std json.Marshal would (following struct tags, etc).
*/
func GenerateMappingFunc(typeExample interface{}) xlate.MappingFunc {
	rt := reflect.TypeOf(typeExample)
	switch rt.Kind() {
	case reflect.Bool:
		return func(_ *xlate.Mapper, dest xlate.Destination, input interface{}) { dest.WriteNull() } // TODO valuemapper thunk
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(_ *xlate.Mapper, dest xlate.Destination, input interface{}) { dest.WriteNull() } // TODO valuemapper thunk
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return func(_ *xlate.Mapper, dest xlate.Destination, input interface{}) { dest.WriteNull() } // TODO valuemapper thunk
	case reflect.Float32:
		return func(_ *xlate.Mapper, dest xlate.Destination, input interface{}) { dest.WriteNull() } // TODO valuemapper thunk
	case reflect.Float64:
		return func(_ *xlate.Mapper, dest xlate.Destination, input interface{}) { dest.WriteNull() } // TODO valuemapper thunk
	case reflect.String:
		return func(_ *xlate.Mapper, dest xlate.Destination, input interface{}) { dest.WriteString(input.(string)) }
	case reflect.Interface:
		panic("TODO")
	case reflect.Struct:
		return GenerateStructAtlas(rt).Visit
	case reflect.Map:
		panic("TODO")
	case reflect.Slice:
		panic("TODO")
	case reflect.Array:
		panic("TODO")
	case reflect.Ptr:
		panic("TODO")
	default:
		panic("unsupported type")
	}
}
