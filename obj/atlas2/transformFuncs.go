package atlas

import "reflect"

type MarshalTransformFunc func(live reflect.Value) (serialable reflect.Value, err error)
type UnmarshalTransformFunc func(serialable reflect.Value) (live reflect.Value, err error)

/*
	Takes a wildcard object which must be `func (live T1) (serialable T2, error)`
	and returns a MarshalTransformFunc and the typeinfo of T2.
*/
func MakeMarshalTransformFunc(fn interface{}) (MarshalTransformFunc, reflect.Type) {
	return nil, nil // TODO
}

/*
	Takes a wildcard object which must be `func (serialable T1) (live T2, error)`
	and returns a UnmarshalTransformFunc and the typeinfo of T1.
*/
func MakeUnmarshalTransformFunc(fn interface{}) (UnmarshalTransformFunc, reflect.Type) {
	return nil, nil // TODO
}
