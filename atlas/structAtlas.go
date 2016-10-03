package atlas

import (
	"reflect"

	"github.com/polydawn/go-xlate"
)

type StructAtlas []StructAtlasRow

type StructAtlasRow struct {
	Name       string
	NameBytes  []byte // []byte(name)
	FieldRoute FieldRoute
}

type FieldRoute []int // a slice because there may be nested jumps to make

func (fr FieldRoute) TraverseToValue(v reflect.Value) reflect.Value {
	for _, i := range fr {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return reflect.Value{}
			}
			v = v.Elem()
		}
		v = v.Field(i)
	}
	return v
}

func (atl StructAtlas) Visit(dispatch *xlate.Mapper, dest xlate.Destination, input interface{}) {
	rvInput := reflect.ValueOf(input)
	dest.OpenMap()
	for _, row := range atl {
		rvField := row.FieldRoute.TraverseToValue(rvInput)
		if !rvField.IsValid() {
			continue
		}
		// TODO : re-implement `row.omitEmpty && IsEmptyValue(rvField)` behavior
		//  (which probably means putting more flags into StructAtlasRow).
		dest.WriteMapKey(row.Name)
		dispatch.Map(dest, rvField.Interface())
	}
	dest.CloseMap()
}
