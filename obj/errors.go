package obj

import (
	"fmt"
	"reflect"
)

var (
	_ error = ErrNoHandler{}
	_ error = ErrUnreachable{}
)

/*
	Error raised as a panic when marshalling or unmarshalling an object, and
	no handler can be found for a referenced type.
*/
type ErrNoHandler struct {
	Valptr interface{}
}

func (e ErrNoHandler) Error() string {
	val_rv := reflect.ValueOf(e.Valptr).Elem()
	return fmt.Sprintf("no machine available in suite for %s of type %T",
		val_rv.Kind(),
		val_rv.Interface())
}

/*
	Error raised from paths the library should be unable to reach.  File bugs.
*/
type ErrUnreachable struct {
	Msg string
}

func (e ErrUnreachable) Fmt(format string, a ...interface{}) ErrUnreachable {
	return ErrUnreachable{fmt.Sprintf(format, a...)}
}

func (e ErrUnreachable) Error() string {
	return "refmt bug: " + e.Msg
}
