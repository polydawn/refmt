package atlas

import "reflect"

type MarshalTransformFunc func(liveForm reflect.Value) (serialForm reflect.Value, err error)
type UnmarshalTransformFunc func(serialForm reflect.Value) (liveForm reflect.Value, err error)

var err_rt = reflect.TypeOf((*error)(nil)).Elem()

/*
	Takes a wildcard object which must be `func (live T1) (serialable T2, error)`
	and returns a MarshalTransformFunc and the typeinfo of T2.
*/
func MakeMarshalTransformFunc(fn interface{}) (MarshalTransformFunc, reflect.Type) {
	fn_rv := reflect.ValueOf(fn)
	if fn_rv.Kind() != reflect.Func {
		panic("no")
	}
	fn_rt := fn_rv.Type()
	if fn_rt.NumIn() != 1 {
		panic("no")
	}
	if fn_rt.NumOut() != 2 {
		panic("no")
	}
	// nothing to do for `fn_rt.In(0)` -- whatever type it is... ugh, yeah, this is sketchy; we should most certainly be able to check this in the builder
	out_rt := fn_rt.Out(0)
	if !fn_rt.Out(1).AssignableTo(err_rt) {
		panic("no")
	}
	return func(liveForm reflect.Value) (serialForm reflect.Value, err error) {
		results := fn_rv.Call([]reflect.Value{liveForm})
		if results[1].IsNil() {
			return results[0], nil
		}
		return results[0], results[1].Interface().(error)
	}, out_rt
}

/*
	Takes a wildcard object which must be `func (serialable T1) (live T2, error)`
	and returns a UnmarshalTransformFunc and the typeinfo of T1.
*/
func MakeUnmarshalTransformFunc(fn interface{}) (UnmarshalTransformFunc, reflect.Type) {
	return nil, nil // TODO
}
