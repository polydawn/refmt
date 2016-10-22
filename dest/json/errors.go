package json

import (
	"reflect"
)

type ErrUnsupportedValue struct {
	Value reflect.Value
	Str   string
}

func (e *ErrUnsupportedValue) Error() string {
	return "unsupported value: " + e.Str
}
