package obj

import "reflect"

// ErrInvalidUnmarshalTarget describes an invalid argument passed to UnmarshalDriver.Bind.
// (Unmarshalling must target a non-nil pointer so that it can address the value.)
type ErrInvalidUnmarshalTarget struct {
	Type reflect.Type
}

func (e ErrInvalidUnmarshalTarget) Error() string {
	if e.Type == nil {
		return "invalid unmarshal target (nil)"
	}
	if e.Type.Kind() != reflect.Ptr {
		return "invalid unmarshal target (non-pointer " + e.Type.String() + ")"
	}
	return "invalid unmarshal target: (nil " + e.Type.String() + ")"
}
